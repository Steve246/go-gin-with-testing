package model

import "fmt"

type Product struct {
	Nama string
	Alamat string 
}
func (c *Product) GetCustomerInfo() string {
	return fmt.Sprintf("Nama : %s, Alamat: %s", c.Nama, c.Alamat)
}
