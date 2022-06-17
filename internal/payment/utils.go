package payment

import (
	"errors"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/gorilla/mux"
)

// Он принимает http-запрос, получает идентификатор из запроса, преобразует его в int64 и возвращает.
func GetQueryId(r *http.Request) (int64, error) {
	id := mux.Vars(r)["id"]
	convertedID, err := ConverteIDtoI64(id)
	if err != nil {
		return 0, err
	}

	return convertedID, nil
}

// Он преобразует строку в int64
func ConverteIDtoI64(id string) (int64, error) {
	convertedID, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}

	return int64(convertedID), nil
}

// Функция называется checkTerminalStatusRow и принимает один аргумент, строку, которая имеет тип
// int64. Функция возвращает ошибку
func checkTerminalStatusRow(row int64) error {
	if row == 0 {
		return errors.New("terminal status")
	}

	return nil
}

// Возвращает true, если данная строка является действительным адресом электронной почты, и false в
// противном случае.
func isEmail(address string) bool {
	_, err := mail.ParseAddress(address)

	return err == nil
}
