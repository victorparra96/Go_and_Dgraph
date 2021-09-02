package utils

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func ReadCsv(ruta string) ([][]string, error) {
	csvFile, err := os.Open(ruta)
	if err != nil {
		log.Printf("Error al abrir el archivo: %v", err)
	}

	defer csvFile.Close()

	r := csv.NewReader(csvFile)

	var rawData [][]string

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Error leyendo la linea: %v", err)
		}

		fmt.Println(record, len(record), "Record")
		trim_slice := strings.Split(record[0], "'")
		if len(trim_slice) > 3 {
			name_product := trim_slice[1 : len(trim_slice)-1]
			union_strings := strings.Join(name_product, "")
			unquote_name_product, err := strconv.Unquote(union_strings)
			if err != nil {
				log.Printf("Error leyendo la linea: %v", err)
			}
			trim_slice[1] = unquote_name_product
			trim_slice[2] = trim_slice[len(trim_slice)-1]
			trim_slice = trim_slice[:3]
			record = trim_slice
		} else {
			record = trim_slice
		}

		rawData = append(rawData, record)
		fmt.Println(record, len(record))
	}

	return rawData, nil
}

func ReadCSVFromUrl(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	r := csv.NewReader(resp.Body)
	r.LazyQuotes = true

	var rawData [][]string

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Error leyendo la linea: %v", err)
		}

		fmt.Println(record)
		trim_slice := strings.Split(record[0], "'")
		if len(trim_slice) > 3 {
			name_product := trim_slice[1 : len(trim_slice)-1]
			union_strings := strings.Join(name_product, "")
			unquote_name_product, err := strconv.Unquote(union_strings)
			if err != nil {
				log.Printf("Error leyendo la linea: %v", err)
			}
			trim_slice[1] = unquote_name_product
			trim_slice[2] = trim_slice[len(trim_slice)-1]
			trim_slice = trim_slice[:3]
			record = trim_slice
		} else {
			record = trim_slice
		}

		rawData = append(rawData, record)
	}

	return rawData, nil
}

func ReadTextStandardFromUrl(url string) ([][]string, [][]string, error) {

	// Create the file
	out, err := os.Create("data/transactions_01.txt")
	if err != nil {
		log.Printf("Error creando el archivo: %v", err)
	}
	defer out.Close()

	// Leer el archivo desde el body
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error leyendo la linea: %v", err)
	}

	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Printf("Error al copiar la data: %v", err)
	}

	file, err := os.Open("data/transactions_01.txt")
	if err != nil {
		log.Printf("Error abriendo el archivo: %v", err)
	}

	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	fileScanner.Buffer(buf, 1024*1024)
	ips_validator := regexp.MustCompile("(?:[0-9]{1,3}\\.){3}[0-9]{1,3}")
	device_validator := regexp.MustCompile("(windows|linux|mac|ios|android)")
	buyers_validator := regexp.MustCompile("(\\(.*?\\))")
	var info_data [][]string
	var products_data [][]string
	for fileScanner.Scan() {
		line := fileScanner.Text()
		new_line := strings.Split(line, "#")
		for _, line := range new_line {
			if len(line) > 2 {
				replace_space := strings.ReplaceAll(line, " ", "")
				id_transaction := replace_space[:12]
				id_buyer := replace_space[12:21]
				ip := ips_validator.FindAllString(replace_space, -1)
				device := device_validator.FindAllString(replace_space, -1)
				buyers := buyers_validator.FindAllString(replace_space, -1)
				// products
				replace_products := strings.ReplaceAll(buyers[0], "(", "")
				replace_products_2 := strings.ReplaceAll(replace_products, ")", "")
				products := strings.Split(replace_products_2, ",")
				temp := make([]string, 0)
				temp = append(temp, id_transaction, id_buyer, ip[0], device[0])
				info_data = append(info_data, temp)
				products_data = append(products_data, products)
			}

		}

	}
	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	return info_data, products_data, nil

}
