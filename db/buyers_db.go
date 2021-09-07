package db

import (
	"context"
	"fmt"
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

	q := fmt.Sprintf(`
		query {
			buyer_id as var(func: eq(id, %q))
	  	}`, buyer.Id)

	mutation := fmt.Sprintf(`
	uid(buyer_id) <id> %q .
    uid(buyer_id) <name> %q .
	uid(buyer_id) <age> "%d" .
	`, buyer.Id, buyer.Name, buyer.Age)

	mu := &api.Mutation{
		SetNquads: []byte(mutation),
	}

	req := &api.Request{
		Query:     q,
		Mutations: []*api.Mutation{mu},
		CommitNow: true,
	}

	if _, err := dg.NewTxn().Do(ctx, req); err != nil {
		log.Fatal(err)
	}

}
