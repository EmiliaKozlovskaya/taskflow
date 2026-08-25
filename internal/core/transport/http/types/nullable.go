package core_http_types

import (
	"encoding/json"

	"github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/domain"
)

// создаем еще одну структуру Nullable которая встраивает в себя Nullable из пакета домена,
// потому что с помощью этой структуры хотим декодировать входящий запрос, а это чисто забота транспортного уровня
type Nullable[T any] struct {
	domain.Nullable[T]
}

// переопределяем метод UnmarshalJSON чтобы при json.NewDecoder().Decode()(под капотом) вызывался этот метод
func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
	//если был вызвал метод UnmarshalJSON был вызван, то json уже точно есть -> n.Set = true
	n.Set = true

	if string(b) == "null" {
		n.Value = nil

		return nil
	}

	var value T
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	n.Value = &value

	return nil
}

func (n *Nullable[T]) ToDomain() domain.Nullable[T] {
	return domain.Nullable[T]{
		Value: n.Value,
		Set:   n.Set,
	}
}

/*
------------------
JSON: {}
Nullable:
	- Value: *nil
	- Set: false

------------------
JSON: {
"phone_number": "+79998887766"
}
Nullable:
	- Value: *"+79998887766"
	- Set: true

------------------
JSON: {
"phone_number": null
}
Nullable:
	- Value: *nil
	- Set: true
------------------
*/
