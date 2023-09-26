package main

type item struct {
	ID               string `json:"id"`
	ShortDescription string `json:"short_description"`
	Price            string `json:"price"`
}

// receipt represents data about a receipt
type receipt struct {
	ID           string `json:"id"`
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchase_date"`
	PurchaseTime string `json:"purchase_time"`
	Items        []item `json:"items"`
	Total        string `json:"total"`
}

func main() {

}
