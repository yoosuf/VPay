package services

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"log"
)

var otpStore = make(map[int]string)

func GenerateOTP() (string, error) {
	bytes := make([]byte, 3) // 6 digits OTP
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func SendOTP(db *sql.DB, userID int) (string, error) {
	otp, err := GenerateOTP()
	if err != nil {
		log.Printf("Failed to generate OTP: %v", err)
		return "", err
	}

	// Store OTP in database
	_, err = db.Exec("INSERT INTO otps (user_id, otp) VALUES (?, ?) ON DUPLICATE KEY UPDATE otp = ?", userID, otp, otp)
	if err != nil {
		log.Printf("Failed to store OTP in database: %v", err)
		return "", err
	}

	// Simulate sending OTP to user (e.g., via email or SMS)
	log.Printf("Generated OTP for user %d: %s", userID, otp)
	return otp, nil
}

func ValidateOTP(db *sql.DB, userID int, otp string) (bool, error) {
	var storedOTP string
	err := db.QueryRow("SELECT otp FROM otps WHERE user_id = ?", userID).Scan(&storedOTP)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("no OTP generated for this user")
		}
		return false, err
	}

	if storedOTP != otp {
		return false, errors.New("invalid OTP")
	}

	_, err = db.Exec("DELETE FROM otps WHERE user_id = ?", userID)
	if err != nil {
		return false, err
	}

	return true, nil
}
