package util

import (
	"net/http"
	"net/mail"
	"strconv"

	"github.com/gorilla/mux"
)

func GetQueryID(r *http.Request) (int64, error) {
	id := mux.Vars(r)["id"]
	convertedID, err := Atoi64(id)
	if err != nil {
		return 0, err
	}

	return convertedID, nil
}

func Atoi64(id string) (int64, error) {
	i, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}

	return int64(i), nil
}

func IsEmail(address string) bool {
	_, err := mail.ParseAddress(address)

	return err == nil
}
