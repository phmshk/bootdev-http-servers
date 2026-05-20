package main

import (
	"slices"
	"strings"
)

func getCleanedBody(msg string) string {
	bannedWords := []string{"kerfuffle", "sharbert", "fornax"}

	splitted := strings.Split(msg, " ")

	for i, word := range splitted {
		lowered := strings.ToLower(word)
		if slices.Contains(bannedWords, lowered) {
			splitted[i] = "****"
		}
	}
	return strings.Join(splitted, " ")
}
