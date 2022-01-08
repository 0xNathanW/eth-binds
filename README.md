# Address -> Go Bindings

 Generate Go bindings for contract deployed at the given address using the Etherscan API.
 
 Essentially the same as the geth abigen tool but reduces arduousness of generating bindings of each contract.
 
 Will update further to increase ease of use even further + add flags for configuration and whatnot.
 
 ## Usage
 
 Clone, set API_KEY in gobind.go as your Etherscan API key, build, then pass addresses in as args.  You can pass multiple addresses in at a time.
