package luhn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"ya41-56/internal/shared/luhn"
)

func TestIsValidLuhn(t *testing.T) {
	valid := []string{
		"79927398713",
		"12345678903",
		"4111111111111111",
	}

	invalid := []string{
		"79927398710",
		"abcdef",
		"",
		"123456",
	}

	for _, v := range valid {
		assert.True(t, luhn.IsValidLuhn(v), "ожидали валидный номер: %s", v)
	}

	for _, v := range invalid {
		assert.False(t, luhn.IsValidLuhn(v), "ожидали невалидный номер: %s", v)
	}
}
