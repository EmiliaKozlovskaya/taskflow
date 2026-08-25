package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

// создаем валидатор в глобальной области видимости
var requestValidator = validator.New()

type validatable interface {
	Validate() error
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf(
			"decode json %v: %w",
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	var (
		err error
	)

	//type asserstion - проверяем удовлетворяет ли dest интерфейсу validatable(имеет ли кастомный метод .Validate())
	v, ok := dest.(validatable)
	if ok {
		//будет значить, что для входящей дто определены кастомные правила валидации
		//а если они определены, значит нам надо ими воспользоваться
		err = v.Validate()
	} else {
		//значит переданная структура валидируется при помощи validator (для базовых типов)
		//для валидации дто хорошо подходит библиотека validator
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf(
			"request validation: %v: %w",
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}
