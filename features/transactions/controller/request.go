package controller

type TransactionRequest struct {
	CartIds   []string `json:"cart_ids"`
	UsingCoin bool     `json:"using_coin"`
}
