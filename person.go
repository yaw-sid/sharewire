package main

type Person struct {
	Address *Address `json:"address"`
	Name    string   `json:"name"`
	Age     uint8    `json:"age"`
}

type Address struct {
	Street   string `json:"street"`
	Postcode string `json:"postcode"`
}
