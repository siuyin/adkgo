package main

import (
	"agenttry/mdl"
	"agenttry/runner"
	"fmt"
	"log"
	"strings"

	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type exchangeRateArgs struct {
	BaseCurrency   string `json:"base_currency"`   // The ISO 4217 currency code for the currency you are converting from (eg. "USD").
	TargetCurrency string `json:"target_currency"` // The ISO 4217 currency code of the currency you are converting to (eg. "EUR").
}
type exchangeRateResult struct {
	Rate string `json:"rate"` // The exchange rate. (eg. 0.93)
}
type baseRates map[string]string

func exchangeRate(ctx tool.Context, inp exchangeRateArgs) exchangeRateResult {
	log.Println("terpau called:", inp)
	rates := map[string]baseRates{"usd": baseRates{"eur": "0.93", "jpy": "157.50", "inr": "83.58"}}
	base := strings.ToLower(inp.BaseCurrency)
	target := strings.ToLower(inp.TargetCurrency)
	rate, ok := rates[base][target]
	if !ok {
		return exchangeRateResult{Rate: fmt.Sprintf("Sorry, we can't convert %s to %s", inp.BaseCurrency, inp.TargetCurrency)}
	}
	log.Println("terpau returning:", rate)
	return exchangeRateResult{Rate: rate}
}

func main() {
	exchangeRateTool, err := functiontool.New(functiontool.Config{
		Name: "ExchangeRate",
		Description: `Looks up and returns the exchange rate between two currencies.
		Check for error messages in the response and format the response in a professional way to the user.`,
	}, exchangeRate)
	if err != nil {
		log.Fatal(err)
	}

	exchangeRateAgent, err := llmagent.New(llmagent.Config{
		Name:  "ExchangeRateAgent",
		Model: mdl.FromEnv(),
		Instruction: `You are specialist exchange rate agent that only answer queries relating to exchange rates.
		You MUST call the ExchangeRate tool. Check the output status and then use the rate from the output.`,
		Tools: []tool.Tool{exchangeRateTool},
	})
	if err != nil {
		log.Fatal(err)
	}

	runner.Run(exchangeRateAgent)

}
