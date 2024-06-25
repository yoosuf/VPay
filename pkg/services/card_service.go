package services

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/yoosuf/VPay/pkg/models"
)

func GenerateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func StoreCardDetails(db *sql.DB, userID int, cardNumber string, expiryDate string) (*models.Card, error) {
	token, err := GenerateToken()
	if err != nil {
		return nil, err
	}

	card := &models.Card{
		UserID:      userID,
		Token:       token,
		ExpiryDate:  expiryDate,
		Last4Digits: cardNumber[len(cardNumber)-4:],
	}

	// Save tokenized card to database
	_, err = db.Exec("INSERT INTO cards (user_id, token, expiry_date, last4_digits) VALUES (?, ?, ?, ?)",
		card.UserID, card.Token, card.ExpiryDate, card.Last4Digits)
	if err != nil {
		return nil, err
	}

	return card, nil
}

func AuthorizeCard(card models.Card, amount float64) (bool, error) {
	// Simulate card authorization process
	if card.ExpiryDate < time.Now().Format("2006-01") {
		return false, errors.New("card expired")
	}

	// Assume all other checks are passed
	return true, nil
}
