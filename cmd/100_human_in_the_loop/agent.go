package main

import (
	"agenttry/mdl"
	"agenttry/runner"
	"fmt"
	"log"

	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/tool"
)

// NOTE: long running functions are implemented as functiontools by setting functiontool.Config IsLongRunning to true.
// See: https://pkg.go.dev/google.golang.org/adk@v0.1.0/tool/functiontool#Config
// See also: https://github.com/google/adk-go/blob/main/tool/functiontool/long_running_function_test.go

const LargeOrderThreshold = 5

type orderStatus struct {
	Status        string
	OrderID       string
	NumContainers int
	Dest          string
	Message       string
}

func placeShippingOrder(numContainer int, dest string, toolContext tool.Context) orderStatus {
	if numContainer < LargeOrderThreshold {
		return autoApproveSmallOrder(numContainer, dest)
	}
	return orderStatus{}
}

func autoApproveSmallOrder(numContainer int, dest string) orderStatus {
	return orderStatus{
		Status:        "Approved",
		OrderID:       fmt.Sprintf("ORD-%d-AUTO", numContainer),
		NumContainers: numContainer,
		Dest:          dest,
		Message:       fmt.Sprintf("Order auto-approved: %d containers to %s", numContainer, dest),
	}
}

func main() {

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
	})
	if err != nil {
		log.Fatal(err)
	}

	runner.Run(currencyAgent)

}
