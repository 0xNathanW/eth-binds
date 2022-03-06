package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

const API_KEY = "YOUR_API_KEY"

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		ABI          string `json:"ABI"`
		ContractName string `json:"ContractName"`
	}
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: gobind [contract address] [package name]")
		return
	}

	if err := getBinding(os.Args[1], os.Args[2]); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Done")
}

func getBinding(address string, pkg string) error {
	abi, name, err := makeRequest(address)
	if err != nil {
		return err
	}

	abis := []string{abi}
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
	if err := ioutil.WriteFile(name+".go", []byte(code), 0600); err != nil {
		return err
	}
	return nil
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
