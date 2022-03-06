package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

var API_KEY string

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		ABI          string `json:"ABI"`
		ContractName string `json:"ContractName"`
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("error loading .env file: %v", err))
	}
	API_KEY = os.Getenv("ETHERSCAN_API_KEY")
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: gobind [contract address] [pkg name]")
		os.Exit(1)
	}
	if err := getBinding(os.Args[:]); err != nil {
		log.Fatal(err)
	}
}

func getBinding(args []string) error {
	abi, name, err := makeRequest(os.Args[0])
	if err != nil {
		return err
	}

	abis := []string{abi}
	pkg := os.Args[1]
	types := []string{name}
	bins := []string{string([]byte{})}
	var sigs []map[string]string
	libs := make(map[string]string)
	aliases := make(map[string]string)

	fmt.Println("Generating bindings for contract ", pkg)

	code, err := bind.Bind(types, abis, bins, sigs, pkg, 0, libs, aliases)
	if err != nil {
		return err
	}
	// Write to file.
	return ioutil.WriteFile(name+".go", []byte(code), 0600)
}

func makeRequest(address string) (string, string, error) {
	url := fmt.Sprintf(
		"https://api.etherscan.io/api?module=contract"+
			"&action=getsourcecode"+
			"&address=%s"+
			"&apikey=%s",
		address, API_KEY,
	)
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return "", "", err
	}
	return response.Result[0].ABI, response.Result[0].ContractName, nil
}
