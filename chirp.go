package main

import (
	"errors"

	"github.com/google/uuid"
)

const (
	minChirpLength = 1
	maxChirpLenth  = 200
)

type chirpParamaters struct {
	Body    string    `json:"body"`
	User_id uuid.UUID `json:"user_id"`
}

func ValidateChirp(chirp chirpParamaters) error {
	err := validateChirpLength(chirp.Body)
	if err != nil {
		return err
	}
	return nil
}

func validateChirpLength(body string) error {
	length := len(body)
	if length < minChirpLength {
		return errors.New("chirp is too short")
	}
	if length > maxChirpLenth {
		return errors.New("chirp is too long")
	}
	return nil
}
