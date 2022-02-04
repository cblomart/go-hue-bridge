package api

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
)

const (
	apiResponse = "[{\"success\":{\"username\":\"<username>\"}}]"
)

func randUser() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func Api(w http.ResponseWriter, r *http.Request) {
	// get a user
	user, _ := randUser()
	// prepare response
	resp := strings.ReplaceAll(apiResponse, "<username>", user)
	// send response
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(resp))
}
