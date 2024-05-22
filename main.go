package main

func main() {

	allContracts := ListStagedContracts()

	for address, contracts := range allContracts {
		for _, contract := range contracts {
			// fmt.Printf("Address: %s, Contract: %s\n", address, contract)
			code := GetStagedContract(address, contract)
			Check(code)
			// fmt.Println(code)
			// break
		}
		// break
	}

}
