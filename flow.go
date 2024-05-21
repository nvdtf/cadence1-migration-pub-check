package main

import (
	"context"

	"github.com/onflow/cadence"
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
			returnRes[address] = append(returnRes[address], c.String())
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
