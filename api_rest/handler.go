package api_rest

import (
	"chi_api_rest_products/db"
	"chi_api_rest_products/models"
	"chi_api_rest_products/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func CreateBuyer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := http.Get("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/buyers")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	dec := json.NewDecoder(strings.NewReader(string(responseData)))

	t, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%T: %v\n", t, t)

	for dec.More() {
		var buyer models.Buyer
		// decode an array value (Message)
		err := dec.Decode(&buyer)
		if err != nil {
			log.Fatal(err)
		}

		db.DbNewBuyer(&buyer)

	}

	t, err = dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%T: %v\n", t, t)

}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	url := "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/products"
	csvLines, err := utils.ReadCSVFromUrl(url)
	if err != nil {
		http.Error(w, "Error with the url", http.StatusBadRequest)
		return
	}

	for _, line := range csvLines {

		price, err := strconv.ParseUint(line[2], 10, 32)
		if err != nil {
			http.Error(w, "Invalid Value", http.StatusBadRequest)
			return
		}

		products := models.Product{
			Id_product:   line[0],
			Name_product: line[1],
			Price:        int(price),
		}

		db.DbNewProduct(&products)

	}
	fmt.Println("Carga de productos realizada")

}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	url := "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/transactions"

	infoData, products, err := utils.ReadTextStandardFromUrl(url)
	if err != nil {
		http.Error(w, "Error with the url", http.StatusBadRequest)
		return
	}

	fmt.Println("Iniciando carga de transacciones...")
	cont := 0
	for _, info := range infoData {
		transaction := models.Transaction{
			Id_transaction: info[0],
			Id_buyer:       info[1],
			Ip:             info[2],
			Device:         info[3],
			Product:        products[cont],
		}
		cont++
		db.DbNewTransaction(&transaction)
	}
	fmt.Println("Carga de transacciones realizada")
}

func ChargeData(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	CreateBuyer(w, r)
	CreateProduct(w, r)
	CreateTransaction(w, r)

}

func GetBuyBuyers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Write([]byte(db.DbGetBuyBuyers()))
}

func GetBuyerInformation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	buyerID := chi.URLParam(r, "buyerId")

	w.Write([]byte(db.DbGetBuyerInformation(buyerID)))
}
