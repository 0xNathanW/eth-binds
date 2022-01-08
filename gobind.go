package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

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
	addresses := os.Args[1:]
	for _, address := range addresses {
		if err := getBinding(address); err != nil {
			fmt.Printf("Error for address %s, %v\n", address, err)
		}
	}
	fmt.Println("Done")
}

func getBinding(address string) error {
	abi, name, err := makeRequest(address)
	if err != nil {
		return err
	}

	abis := []string{abi}
	pkg := strings.ToLower(name)
	types := []string{name}
	bins := []string{string([]byte{})}
	var sigs []map[string]string
	libs := make(map[string]string)
	aliases := make(map[string]string)

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
