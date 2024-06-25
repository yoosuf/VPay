package models

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Card struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Token       string `json:"token"`
	ExpiryDate  string `json:"expiry_date"`
	Last4Digits string `json:"last4_digits"`
}

type Transaction struct {
	ID              int     `json:"id"`
	UserID          int     `json:"user_id"`
	CardID          int     `json:"card_id"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	Status          string  `json:"status"`
	TransactionDate string  `json:"transaction_date"`
}
