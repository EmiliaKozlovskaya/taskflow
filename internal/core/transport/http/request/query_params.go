package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/Emilia20112005/golang-todoapp/internal/core/errors"
)

//query parameters будут двух видов (для get_users limit & offset) и (для get_statistics дата from & to)
//не будем мудрить напишем две отдельные функции

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid integer: %v, %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}
	return &val, nil
}
