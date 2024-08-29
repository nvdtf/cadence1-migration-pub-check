package main

import (
	"context"
	"os"
	"strings"

	"github.com/onflow/cadence"
	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/flow-go-sdk/access/http"
)

func ListStagedContracts() map[string][]string {

	ctx := context.Background()
	// flowClient, err := http.NewClient(http.TestnetHost)
	flowClient, err := http.NewClient(http.MainnetHost)
	if err != nil {
		panic(err)
	}

	script := []byte(`
		// testnet
		// import MigrationContractStaging from 0x2ceae959ed1a7e7a
		import MigrationContractStaging from 0x56100d46aa9b0212

		access(all) fun main(): {Address: [String]} {
			return MigrationContractStaging.getAllStagedContracts()
		}
    `)

	value, err := flowClient.ExecuteScriptAtLatestBlock(ctx, script, nil)
	if err != nil {
		panic(err)
	}

	returnRes := make(map[string][]string)

	v := value.(cadence.Dictionary)
	pairs := make([]struct {
		Key   string
		Value string
	}, len(v.Pairs))

	for i, pair := range v.Pairs {
		address := pair.Key.String()
		contracts := pair.Value.(cadence.Array).Values

		for _, c := range contracts {
			cName := strings.Replace(c.String(), "\"", "", -1)
			returnRes[address] = append(returnRes[address], cName)
		}

		pairs[i] = struct {
			Key   string
			Value string
		}{
			Key:   pair.Key.String(),
			Value: pair.Value.String(),
		}
	}

	return returnRes

}

func GetStagedContract(address string, name string) (string, bool) {

	// fungible token
	if (address == "9a0766d93b6608b7" || address == "f233dcee88fe0abe") && name == "FungibleToken" {
		return GetStagedFungibleToken(), true
	}

	// non-fungible token
	if (address == "631e88ae7f1d7c20" || address == "1d7e57aa55817448") && name == "NonFungibleToken" {
		return GetStagedNonFungibleToken(), true
	}

	// view resolver
	if (address == "631e88ae7f1d7c20" || address == "1d7e57aa55817448") && name == "ViewResolver" {
		return GetStagedViewResolver(), true
	}

	// burner
	if (address == "9a0766d93b6608b7" || address == "f233dcee88fe0abe") && name == "Burner" {
		return GetStagedBurner(), true
	}

	ctx := context.Background()
	// flowClient, err := http.NewClient(http.TestnetHost)
	flowClient, err := http.NewClient(http.MainnetHost)
	if err != nil {
		panic(err)
	}

	script := []byte(`
		// testnet
		// import MigrationContractStaging from 0x2ceae959ed1a7e7a
		import MigrationContractStaging from 0x56100d46aa9b0212

		access(all) fun main(contractAddress: Address, contractName: String): String? {
			return MigrationContractStaging.getStagedContractCode(address: contractAddress, name: contractName)
			//return "\n"
		}
    `)

	addressCdc, err := common.HexToAddress(address)
	if err != nil {
		panic(err)
	}

	nameCdc, err := cadence.NewString(name)
	if err != nil {
		panic(err)
	}

	value, err := flowClient.ExecuteScriptAtLatestBlock(ctx, script, []cadence.Value{
		cadence.NewAddress(addressCdc),
		nameCdc,
	})
	if err != nil {
		panic(err)
	}

	// res := value.(cadence.String).String()

	optValue, ok := value.(cadence.Optional)
	if !ok {
		panic("not optional")
	}
	if optValue.Value == nil {
		return "", false
	}

	res := optValue.Value.(cadence.String).String()

	res = strings.Trim(res, "\"")
	res = strings.ReplaceAll(res, "\\\\n", "new_line_in_string")
	res = strings.ReplaceAll(res, "\\n", "\n")
	res = strings.ReplaceAll(res, "new_line_in_string", "\\\\n")
	res = strings.ReplaceAll(res, "\\r", "\r")
	res = strings.ReplaceAll(res, "\\\"", "\"")
	res = strings.ReplaceAll(res, "\\\\\"", "\\\"")
	res = strings.ReplaceAll(res, "\\t", "\t")
	// res = strings.ReplaceAll(res, "\\", "\"")
	// res = strings.ReplaceAll(res, "\"\"\"", "\\\"")
	// res = strings.ReplaceAll(res, "\"\"n\"\"", "\\n")

	return res, true
}

func GetStagedFungibleToken() string {
	res, err := os.ReadFile("contracts/ft.cdc")
	if err != nil {
		panic(err)
	}
	return string(res)
}

func GetStagedNonFungibleToken() string {
	res, err := os.ReadFile("contracts/nft.cdc")
	if err != nil {
		panic(err)
	}
	return string(res)
}

func GetStagedViewResolver() string {
	res, err := os.ReadFile("contracts/view-resolver.cdc")
	if err != nil {
		panic(err)
	}
	return string(res)
}

func GetStagedBurner() string {
	res, err := os.ReadFile("contracts/burner.cdc")
	if err != nil {
		panic(err)
	}
	return string(res)
}
