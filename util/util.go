package util

import "math/rand"

func GeneratePrettyId(n int) string {
	symbols := make([]rune, 0, 36)
	for i := 'a'; i <= 'z'; i++ {
		symbols = append(symbols, i)
	}
	for i := '0'; i <= '9'; i++ {
		symbols = append(symbols, i)
	}

	symbolsLen := len(symbols)

	parts := make([]rune, 0, n)
	for i := 0; i < n; i++ {
		parts = append(parts, symbols[rand.Intn(symbolsLen)])
	}

	return string(parts)
}
