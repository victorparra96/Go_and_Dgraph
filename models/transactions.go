package models

type Transaction struct {
	Id_transaction string   `json:"id_transaction"`
	Id_buyer       string   `json:"id_buyer"`
	Ip             string   `json:"ip"`
	Device         string   `json:"device"`
	Product        []string `json:"product"`
}

func Add_product(id, buyer, ip, device string, values [][]string) {
	//var prod string
	/* for _, line := range values {
		transaction := Transaction{
			Id_transaction: id,
			Id_buyer:       buyer,
			Ip:             ip,
			Device:         device,
			Product: []string{
				values[][],
			},
		}
	} */

	//fmt.Println(transaction)
}
