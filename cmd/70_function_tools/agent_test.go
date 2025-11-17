package main

import (
	"testing"
)

func TestExchRate(t *testing.T) {
	dat := []struct {
		b string
		t string
		r string
	}{
		{"USD", "EUR", "0.93"},
		{"USD", "JPY", "157.50"},
		{"USD", "INR", "83.58"},
		{"SGD", "INR", "Sorry, we can't convert SGD to INR"},
		{"USD", "SGD", "Sorry, we can't convert USD to SGD"},
	}
	for i, d := range dat {
		outp := exchangeRate(nil, exchangeRateArgs{BaseCurrency: d.b, TargetCurrency: d.t})
		if outp.Rate != d.r {
			t.Errorf("case %d: expected: %s, got: %s", i, d.r, outp.Rate)
		}

	}
}
