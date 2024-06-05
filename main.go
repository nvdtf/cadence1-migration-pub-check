package main

import (
	"fmt"
	"time"
)

func main() {

	// code, _ := GetStagedContract("0x74ad08095d92192a", "FRC20Indexer")
	// Check(code)

	startTime := time.Now()
	allContracts := ListStagedContracts()

	for address, contracts := range allContracts {
		for _, contract := range contracts {

			// base 64 encoded staged contracts
			if (address == "0x34f3140b7f54c743" || address == "0xb45e7992680a0f7f" || address == "0x2d0d952e760d1770") && (contract == "CricketMoments" || contract == "FazeUtilityCoin") {
				continue
			}

			// fmt.Printf("Address: %s, Contract: %s\n", address, contract)
			code, found := GetStagedContract(address, contract)
			if !found {
				panic(fmt.Sprintf("staged contract code not found: %s.%s", address, contract))
			}

			res := Check(code)
			if len(res) > 0 {
				fmt.Println("---------------------------------")
				fmt.Printf("Address: %s, Contract: %s\n", address, contract)
			}
			for _, v := range res {
				fmt.Printf("❗ Resource %s is exposing %s\n", v.resourceName, v.functionName)
			}
			// break
		}
		// break
	}

	fmt.Printf("\n✅ Done! %d contracts checked in %s\n", len(allContracts), time.Since(startTime))

}
