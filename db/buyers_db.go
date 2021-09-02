package db

import (
	"context"
	"encoding/json"
	"log"

	"chi_api_rest_products/dgraph_client"
	"chi_api_rest_products/models"

	"github.com/dgraph-io/dgo/v2/protos/api"
)

func DbNewBuyer(buyer *models.Buyer) {

	dg, cancel := dgraph_client.GetDgraphClient()
	defer cancel()

	op := &api.Operation{}
	op.Schema = `
			id: string @index(hash).
			name: string .
			age: int .

			type Buyer {
				id
				name
				age
			}
		`

	ctx := context.Background()
	if err := dg.Alter(ctx, op); err != nil {
		log.Fatal(err)
	}

	mu := &api.Mutation{
		CommitNow: true,
	}

	pb, err := json.Marshal(buyer)
	if err != nil {
		log.Fatal(err)
	}

	mu.SetJson = pb
	_, err = dg.NewTxn().Mutate(ctx, mu)
	if err != nil {
		log.Fatal(err)
	}

}
