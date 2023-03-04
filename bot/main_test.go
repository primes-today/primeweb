package main

import (
	"math/big"
	"testing"
)

func TestFormat(t *testing.T) {
	cases := map[*big.Int]string{
		big.NewInt(1):          "1",
		big.NewInt(12):         "12",
		big.NewInt(123):        "123",
		big.NewInt(1234):       "1,234",
		big.NewInt(12345):      "12,345",
		big.NewInt(123456):     "123,456",
		big.NewInt(1234567):    "1,234,567",
		big.NewInt(12345678):   "12,345,678",
		big.NewInt(123456789):  "123,456,789",
		big.NewInt(1234567890): "1,234,567,890",
	}

	for c, e := range cases {
		if a := format(c); a != e {
			t.Errorf("expected %s, got %s", e, a)
		}
	}
}
