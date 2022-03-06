# Address -> Go Bindings

 A modified version of geth's abigen for the lazy.  Given addresses, will fetch ABI from etherscan, then generate the bindings. 
 
 Binding struct and the output file name are set as the contract name, package name is set as the lowercase of the contract name.
 
 Will update further to increase ease of use even further + add flags for configuration and whatnot.
 
 ## Usage
 
 Clone, set API_KEY in gobind.go as your Etherscan API key, build, then pass addresses in as args.  You can pass multiple addresses in at a time.
