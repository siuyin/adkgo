package main

import (
	"agenttry/mdl"
	"agenttry/runner"
	"fmt"
	"log"
	"strings"

	"github.com/mnogu/go-calculator"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type calcArgs struct {
	Expr string `json:"expression"` // Expression of math to be calculated. eg1. "2 + 3 / 4", eg2. "1.1 /2.2 + 3.3 - 4.4"
}
type calcOut struct {
	Res    string `json:"result"` // String representation of calculation result. eg1. "42", eg2. "1.23"
	Status string `json:"status"` // Status of calculation. eg1. "OK" or eg2. "Error: division by 0."
}

func calc(ctx tool.Context, inp calcArgs) calcOut {
	val, err := calculator.Calculate(inp.Expr)
	if err != nil {
		return calcOut{Res: "", Status: fmt.Sprintf("Error: %v", err)}
	}
	return calcOut{Res: fmt.Sprintf("%v", val), Status: "OK"}
}

type exchangeRateArgs struct {
	BaseCurrency   string `json:"base_currency"`   // The ISO 4217 currency code for the currency you are converting from (eg. "USD").
	TargetCurrency string `json:"target_currency"` // The ISO 4217 currency code of the currency you are converting to (eg. "EUR").
}
type exchangeRateResult struct {
	Rate string `json:"rate"` // The exchange rate. (eg. 0.93)
}
type baseRates map[string]string

func exchangeRate(ctx tool.Context, inp exchangeRateArgs) exchangeRateResult {
	rates := map[string]baseRates{"usd": baseRates{"eur": "0.93", "jpy": "157.50", "inr": "83.58"}}
	base := strings.ToLower(inp.BaseCurrency)
	target := strings.ToLower(inp.TargetCurrency)
	rate, ok := rates[base][target]
	if !ok {
		return exchangeRateResult{Rate: fmt.Sprintf("Sorry, we can't convert %s to %s", inp.BaseCurrency, inp.TargetCurrency)}
	}
	return exchangeRateResult{Rate: rate}
}

type paymentMethodArgs struct {
	Method string `json:"method"` // The name of the payment method. It should be descriptive. Eg. "platinum credit card" or "bank transfer".
}
type paymentMethodFee struct {
	Fee string `json:"fee"` // Fee rate associated with payment method. To get a percentage rate, multiply fee rate by 100.
}

func feeForPaymentMethod(ctx tool.Context, inp paymentMethodArgs) paymentMethodFee {
	fees := map[string]float32{
		"platinum credit card": 0.02,
		"gold debit card":      0.035,
		"bank transfer":        0.01,
	}
	m := strings.ToLower(inp.Method)
	fee, ok := fees[m]
	if !ok {
		return paymentMethodFee{fmt.Sprintf("Sorry, we do not accept %q as a payment method.", inp.Method)}
	}
	return paymentMethodFee{fmt.Sprintf("%v", fee)}
}

func main() {
	calcTool, err := functiontool.New(functiontool.Config{
		Name: "Calculator",
		Description: `Takes a math expression as a string and outputs and object with 'result' and 'status' fields.
		Check the status field for error. Strictly use result only if status shows OK. All other cases indicate an error.
		`,
	}, calc)
	if err != nil {
		log.Fatal(err)
	}

	exchangeRateTool, err := functiontool.New(functiontool.Config{
		Name: "ExchangeRate",
		Description: `Looks up and returns the exchange rate between two currencies.
		Check for error messages in the response and format the response in a professional way to the user.`,
	}, exchangeRate)
	if err != nil {
		log.Fatal(err)
	}

	feeForPaymentMethodTool, err := functiontool.New(functiontool.Config{
		Name: "FeeForPaymentMethod",
		Description: `Looks up the fee charged for using a payment method.
	Check for error messages in the response and format the overall response in a professional way to the user.`,
	}, feeForPaymentMethod)
	if err != nil {
		log.Fatal(err)
	}

	currencyAgent, err := llmagent.New(llmagent.Config{
		Name:  "CurrencyAgent",
		Model: mdl.FromEnv(),
		Instruction: `You are a smart currency conversion assistant.
	For currency conversion requests:
	1. Use 'FeeForPaymentMethod' to find transaction fees.
	2. Use 'ExchangeRate' to get currency conversion rates.
	3. Check the response for each tool call for errors.
	4. You MUST use the Calculator tool to calculate the final amount after fees based on the output of 'FeeForPaymentMethod' and 'ExchangeRate' methods and provide a clear breakdown.
	5. First, state the final converted amount.
	   Then, explain how you got the amount by showing the intermediate amounts. Your explanation must include: the fee percentage and its value in the original currency, the amount remaining after the fee, and the exchange rate used for the final conversion.
	
	If any tool returns an error, explain the issue clearly to the user.
	`,
		Tools: []tool.Tool{feeForPaymentMethodTool, exchangeRateTool, calcTool},
	})
	if err != nil {
		log.Fatal(err)
	}

	runner.Run(currencyAgent)

}
