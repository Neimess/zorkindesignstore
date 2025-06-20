package validator

import (
    "sync"

    "github.com/go-playground/validator/v10"
)

var (
    once     sync.Once
    instance *validator.Validate
)

func GetValidator() *validator.Validate {
    once.Do(func() {
        v := validator.New()
        instance = v
    })
    return instance
}
