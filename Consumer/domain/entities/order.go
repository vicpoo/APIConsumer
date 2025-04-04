// order.go
package models

type Order struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Quantity int    `json:"quantity"`
	BedName  string `json:"bed_name"`
}

func NewOrder(name, lastName, phone, email, bedName string, quantity int) *Order {
	return &Order{
		Name:     name,
		LastName: lastName,
		Phone:    phone,
		Email:    email,
		Quantity: quantity,
		BedName:  bedName,
	}
}
