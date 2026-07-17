// доменная сущность
package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/Emilia20112005/golang-todoapp/internal/core/errors"
)

type User struct {
	ID      int
	Version int //не только для оптимистической блокировки(уровня репозитория), но и для отслеживания изменений

	FullName    string
	PhoneNumber *string //указатель потому что может быть null в БД
}

// конструктор для бд
func NewUser(
	id int,
	version int,
	fullName string,
	phoneNumber *string,
) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

// конструктор для дто
func NewUserUninitialized(
	fullName string,
	phoneNumber *string,
) User {
	return NewUser(
		UninitializedID,
		UninitializedVersion,
		fullName,
		phoneNumber,
	)
}

// тут будем проверять соответствует ли пришедшая доменная сущность бизнес требованиям
func (u *User) Validate() error {
	fullNameLength := len([]rune(u.FullName)) //количество символов в строке (len(u.FullName)-кол-во байт)
	if fullNameLength < 3 || fullNameLength > 100 {
		return fmt.Errorf(
			"invalid `FullName` len: %d: %w",
			fullNameLength,
			core_errors.ErrInvalidArgument,
		)
	}

	if u.PhoneNumber != nil {
		phoneNumberLen := len([]rune(*u.PhoneNumber))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf(
				"invalid `PhoneNumber` len: %d: %w",
				phoneNumberLen,
				core_errors.ErrInvalidArgument,
			)
		}

		re := regexp.MustCompile(`^\+[0-9]+$`) //регулярное выражение, взяли из файлов миграции
		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf(
				"invalid `PhoneNumber` format: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}
	return nil
}

// структура для описания поведения при изменении пользователя (чтобы знать как что менять(как новый патч на игру который её меняет))
type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

// нужно уметь валидировать этот патч
func (p *UserPatch) Validate() error {
	if p.FullName.Set && p.FullName.Value == nil {
		//то есть кто-то пытается выставить FullName в бд null (так нельзя!! тк поле not null)
		return fmt.Errorf(
			"`FullName` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	//в остальных всех случаях патч валидный, и остальная валидация уже будет относиться к пользователю
	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil { //убеждаемся что патч валидный
		return fmt.Errorf("validate user patch: %w", err)
	}
	//нужно для начала провалидировать полученного пользователя
	tmp := *u
	//применяем патч к временной переменной и валидируем временную переменную

	if patch.FullName.Set {
		tmp.FullName = *patch.FullName.Value
	}
	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}
	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}
	//обратно присваиваем значение уже провалидированные изменения на входящего юзера
	*u = tmp

	return nil
}
