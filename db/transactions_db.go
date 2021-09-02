package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"chi_api_rest_products/dgraph_client"
	"chi_api_rest_products/models"

	"github.com/dgraph-io/dgo/v2/protos/api"
)

func DbNewTransaction(transaction *models.Transaction) {

	dg, cancel := dgraph_client.GetDgraphClient()
	defer cancel()

	op := &api.Operation{}
	op.Schema = `
			id_transaction: string @index(hash).
			id_buyer: [uid] @reverse .
			ip: string @index(hash) .
			device: string .
			product: [uid] .

			type Transaction {
				id_transaction
				id_buyer
				ip
				device
				product
			}
		`

	ctx := context.Background()
	if err := dg.Alter(ctx, op); err != nil {
		log.Fatal(err)
	}

	// Llenar los productos
	result := ""
	for _, line := range transaction.Product {
		result += fmt.Sprintf(",%q", line)
	}

	products := strings.Replace(result, ",", "", 1)

	pb, err := json.Marshal(transaction.Id_buyer)
	if err != nil {
		log.Fatal(err)
	}

	q := fmt.Sprintf(`
		query {
			v as var(func: eq(id, %q))
			p as var(func: eq(id_product, %s))
	  	}`, pb[7:15], products)

	mutation := fmt.Sprintf(`
	_:Transaction <id_buyer>  uid(v) .
	_:Transaction <product>  uid(p) .
	_:Transaction <id_transaction>  %q .
	_:Transaction <device> %q .
	_:Transaction <ip> %q .
	`, transaction.Id_transaction, transaction.Device, transaction.Ip)

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

func DbGetBuyBuyers() []byte {

	dg, cancel := dgraph_client.GetDgraphClient()
	defer cancel()

	resp, err := dg.NewTxn().Query(context.Background(), `{
		buyer_list(func: has(id)) @filter(gt(count(~id_buyer), 0)){
			id
			name
			age
			total_buys: count(~id_buyer)
		}
	  }`)

	if err != nil {
		log.Fatal(err)
	}

	return resp.Json
}

func DbGetBuyerInformation(buyerID string) []byte {

	dg, cancel := dgraph_client.GetDgraphClient()
	defer cancel()

	resp, err := dg.NewTxn().Query(context.Background(), fmt.Sprintf(`{
		var(func: has(id), first: 5) @filter(NOT eq(id, %s)) {
			id
			name
			age
			~id_buyer{
				id_transaction
				product {
					id_product
					n_product as name_product
		  		}
		  	}
	  	}
	  	var(func: eq(id, %s)){
			id
			name
			age
			~id_buyer{
				ips as ip
			}
	  	}
	  	buy_history(func: eq(id, %s)) {
			id
			name
			age
			transaction_list: ~id_buyer {
		  		id_transaction
				ip
		  		device
		  		product{
					id_product
					name_product
					price
		  		}
			}
	  	}
		find_ip(func: eq(ip, val(ips))){
			buyer_list: @groupby(ip){
				total_buyers: count(uid)
			}
		}
	  	find_product_recomendation(func: eq(name_product, val(n_product))){
			id_product
			name_product
	 	}
	  }`, buyerID, buyerID, buyerID))

	if err != nil {
		log.Fatal(err)
	}

	replace_products := strings.ReplaceAll(string(resp.Json), "@", "")

	return []byte(replace_products)

}
