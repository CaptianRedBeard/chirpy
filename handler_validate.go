package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type paramaters struct {
		Body string `json:"body"`
	}

	type returnError struct {
		Error string `json:"error"`
	}

	params := paramaters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respBody := returnError{
			Error: err.Error(),
		}

		dat, _ := json.Marshal(respBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(dat)
		return
	}

	if len(params.Body) > 140 {
		respBody := returnError{
			Error: "Chirp is too long",
		}

		dat, _ := json.Marshal(respBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	verbotenWords := []string{"kerfuffle", "sharbert", "fornax"}

	var sentence strings.Builder
	for _, word := range strings.Fields(params.Body) {
		var cleanedWord strings.Builder
		for _, char := range word {
			if unicode.IsLetter(char) {
				cleanedWord.WriteRune(char)
			}
		}

		vFlag := false
		for _, vWord := range verbotenWords {
			if strings.ToLower(cleanedWord.String()) == vWord {
				vFlag = true
				sentence.WriteString("**** ")
				break
			}
		}
		if !vFlag {
			sentence.WriteString(word + " ")
		}
	}

	finalSentence := strings.TrimSpace(sentence.String())

	if finalSentence != params.Body {
		params.Body = finalSentence
	}

	type returnVals struct {
		Cleaned_body string `json:"cleaned_body"`
	}

	respBody := returnVals{
		Cleaned_body: params.Body,
	}

	dat, _ := json.Marshal(respBody)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}
