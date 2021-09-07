package models

type Transaction struct {
	Id_transaction string   `json:"id_transaction"`
	Id_buyer       string   `json:"id_buyer"`
	Ip             string   `json:"ip"`
	Device         string   `json:"device"`
	Product        []string `json:"product"`
}
