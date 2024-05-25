package main

func main() {

	code, _ := GetStagedContract("0xcbdb5a7b89c3c844", "PriceOracle")
	Check(code)

	// allContracts := ListStagedContracts()

	// for address, contracts := range allContracts {
	// 	for _, contract := range contracts {

	// 		// base 64 encoded staged contracts
	// 		if (address == "0x34f3140b7f54c743" || address == "0xb45e7992680a0f7f" || address == "0x2d0d952e760d1770") && (contract == "CricketMoments" || contract == "FazeUtilityCoin") {
	// 			continue
	// 		}

	// 		fmt.Printf("Address: %s, Contract: %s\n", address, contract)
	// 		code, found := GetStagedContract(address, contract)
	// 		if !found {
	// 			panic(fmt.Sprintf("staged contract code not found: %s.%s", address, contract))
	// 		}
	// 		Check(code)
	// 		// break
	// 	}
	// 	// break
	// }

}
