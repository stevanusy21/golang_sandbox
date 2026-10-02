package request

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func DecodeJSON[T any](r *http.Request) (T, error) {
	var payload T

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return payload, errors.New("Format payload JSON tidak valid")
	}

	if value, ok := any(payload).(interface{ Validate() error }); ok {
		if err := value.Validate(); err != nil {
			return payload, err
		}
	}

	return payload, nil
}

func GetIntParam(r *http.Request, paramName string, required bool) (int, error) {
	valStr := r.PathValue(paramName)
	if valStr == "" {
		valStr = r.URL.Query().Get(paramName)
	}

	if valStr == "" {
		if required {
			return 0, errors.New("Parameter " + paramName + " tidak ditemukan")
		}
		return 0, nil
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, errors.New("Format parameter " + paramName + " harus berupa angka")
	}

	return val, nil
}

func GetStringParam(r *http.Request, paramName string, required bool) (string, error) {
	valStr := r.PathValue(paramName)
	if valStr == "" {
		valStr = r.URL.Query().Get(paramName)
	}

	if valStr == "" {
		if required {
			return "", errors.New("Parameter " + paramName + " tidak ditemukan")
		}
		return "", nil
	}

	return valStr, nil
}
