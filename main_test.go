package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkMain(b *testing.B) {
	//for i := 0; i < b.N; i++ {
	main()
	//}
}
func TestGetLowestPrise(t *testing.T) {
	assert := assert.New(t)
	type GetPrice struct {
		a, b, c string
		d, e, f int
	}
	var tests = []struct {
		input    GetPrice
		expected int
	}{
		{GetPrice{`Kiev`, `Odessa`, `[3,4,5]`, 2, 3, 4}, 2},
		{GetPrice{`Kiev`, `Odessa`, `[3,4,5]`, 0, 0, 0}, 0},
		{GetPrice{`Kiev`, `Odessa`, `[3,4,5]`, 200, 0, 400}, 200},
	}

	for _, test := range tests {
		assert.Equal(getLowestPrise(test.input.a, test.input.b, test.input.c, test.input.d, test.input.e, test.input.f), test.expected)
	}
}
