package main

import (
	"context"
	"strings"

	"github.com/onflow/cadence"
	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/flow-go-sdk/access/http"
)

func ListStagedContracts() map[string][]string {

	ctx := context.Background()
	flowClient, err := http.NewClient(http.TestnetHost)
	if err != nil {
		panic(err)
	}

	script := []byte(`
		import MigrationContractStaging from 0x2ceae959ed1a7e7a

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

func GetStagedContract(address string, name string) string {

	ctx := context.Background()
	flowClient, err := http.NewClient(http.TestnetHost)
	if err != nil {
		panic(err)
	}

	script := []byte(`
		import MigrationContractStaging from 0x2ceae959ed1a7e7a

		access(all) fun main(contractAddress: Address, contractName: String): String {
			return MigrationContractStaging.getStagedContractCode(address: contractAddress, name: contractName)!
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

	res := strings.Trim(value.String(), "\"")
	res = strings.ReplaceAll(res, "\\\\n", "doobidydoo")
	res = strings.ReplaceAll(res, "\\n", "\n")
	res = strings.ReplaceAll(res, "doobidydoo", "\\\\n")
	res = strings.ReplaceAll(res, "\\\"", "\"")
	res = strings.ReplaceAll(res, "\\t", "\t")
	res = strings.ReplaceAll(res, "\\", "\"")
	res = strings.ReplaceAll(res, "\"\"\"", "\\\"")
	res = strings.ReplaceAll(res, "\"\"n\"\"", "\\n")

	return res
}
