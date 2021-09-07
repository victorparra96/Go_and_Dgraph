package db

import (
	"context"
	"fmt"
	"log"

	"chi_api_rest_products/dgraph_client"
	"chi_api_rest_products/models"

	"github.com/dgraph-io/dgo/v2/protos/api"
)

func DbNewProduct(product *models.Product) {

	dg, cancel := dgraph_client.GetDgraphClient()
	defer cancel()

	op := &api.Operation{}
	op.Schema = `
			id_product: string @index(hash).
			name_product: string @index(hash).
			price: int .

			type Product {
				id_product
				name_product
				price
			}
		`

	ctx := context.Background()
	if err := dg.Alter(ctx, op); err != nil {
		log.Fatal(err)
	}

	q := fmt.Sprintf(`
		query {
			product_id as var(func: eq(id_product, %q))
	  	}`, product.Id_product)

	mutation := fmt.Sprintf(`
	uid(product_id) <id_product> %q .
    uid(product_id) <name_product> %q .
	uid(product_id) <price> "%d" .
	`, product.Id_product, product.Name_product, product.Price)

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
