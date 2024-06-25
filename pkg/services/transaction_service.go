package services

import (
	"database/sql"
	"errors"
	"time"

	"github.com/yoosuf/VPay/pkg/models"
)

func ProcessTransaction(db *sql.DB, user models.User, card models.Card, amount float64, currency string) (*models.Transaction, error) {
	// Step 1: Authorize the card
	authorized, err := AuthorizeCard(card, amount)
	if err != nil || !authorized {
		return nil, errors.New("authorization failed")
	}

	// Step 2: Simulate transaction processing
	transaction := &models.Transaction{
		UserID:          user.ID,
		CardID:          card.ID,
		Amount:          amount,
		Currency:        currency,
		Status:          "Approved",
		TransactionDate: time.Now().Format(time.RFC3339),
	}

	// Save transaction to database
	_, err = db.Exec("INSERT INTO transactions (user_id, card_id, amount, currency, status, transaction_date) VALUES (?, ?, ?, ?, ?, ?)",
		transaction.UserID, transaction.CardID, transaction.Amount, transaction.Currency, transaction.Status, transaction.TransactionDate)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
