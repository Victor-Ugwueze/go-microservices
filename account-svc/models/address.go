package models


type AddressDetail struct {
	State string      `json:"state"`
	City string        `json:"city"`
	PostalCode int16   `json:"postalCode"`
}
 
type Address struct {}