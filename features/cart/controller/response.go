package controller

import products "greenenvironment/features/products/controller"

type CartResponse struct {
	User  User        `json:"user"`
	Items []CartItems `json:"items"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Address  string `json:"address"`
}

type CartItems struct {
	ID       string                   `json:"id"`
	Quantity int                      `json:"quantity"`
	Product  products.ProductResponse `json:"product"`
}
