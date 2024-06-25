package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/yoosuf/VPay/pkg/models"
	"github.com/yoosuf/VPay/pkg/services"
)

type PaymentRequest struct {
	UserID     int     `json:"user_id"`
	CardNumber string  `json:"card_number"`
	ExpiryDate string  `json:"expiry_date"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	OTP        string  `json:"otp"`
}

func ProcessPaymentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var paymentRequest PaymentRequest
		if err := json.NewDecoder(r.Body).Decode(&paymentRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		valid, err := services.ValidateOTP(db, paymentRequest.UserID, paymentRequest.OTP)
		if err != nil || !valid {
			http.Error(w, "invalid OTP", http.StatusUnauthorized)
			return
		}

		user := models.User{ID: paymentRequest.UserID, Name: "Alice", Email: "alice@example.com"} // Fetch user from DB
		card, err := services.StoreCardDetails(db, user.ID, paymentRequest.CardNumber, paymentRequest.ExpiryDate)
		if err != nil {
			http.Error(w, "error storing card details", http.StatusInternalServerError)
			return
		}

		transaction, err := services.ProcessTransaction(db, user, *card, paymentRequest.Amount, paymentRequest.Currency)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, _ := json.Marshal(transaction)
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func RequestOTPHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}

		log.Printf("Received OTP request for user ID: %d", userID)

		otp, err := services.SendOTP(db, userID)
		if err != nil {
			log.Printf("Error generating OTP: %v", err)
			http.Error(w, "error generating OTP", http.StatusInternalServerError)
			return
		}

		// Simulate sending OTP (e.g., via email/SMS)
		response := map[string]string{"otp": otp}
		responseJSON, _ := json.Marshal(response)
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}
}
