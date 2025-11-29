package main

import (
	"agenttry/mdl"
	"agenttry/runner"
	"fmt"
	"log"
	"time"

	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// NOTE: long running functions are implemented as functiontools by setting functiontool.Config IsLongRunning to true.
// See: https://pkg.go.dev/google.golang.org/adk@v0.1.0/tool/functiontool#Config
// See also: https://github.com/google/adk-go/blob/main/tool/functiontool/long_running_function_test.go

const LargeOrderThreshold = 5

// PlaceShippingOrderArgs holds the arguments for the placeShippingOrder function.
type PlaceShippingOrderArgs struct {
	NumContainers int    `json:"num_containers"`
	Dest          string `json:"dest"`
}

// ShippingOrderStatus holds the result of the placeShippingOrder function.
type ShippingOrderStatus struct {
	Status        string `json:"status"`
	OrderID       string `json:"order_id,omitempty"`
	NumContainers int    `json:"num_containers"`
	Dest          string `json:"destination"`
	Message       string `json:"message"`
	JobID         string `json:"job_id,omitempty"` // For long-running tasks
	Error         string `json:"error,omitempty"`  // New error field
}

func placeShippingOrder(ctx tool.Context, args PlaceShippingOrderArgs) ShippingOrderStatus {
	if args.NumContainers < LargeOrderThreshold {
		return autoApproveSmallOrder(args.NumContainers, args.Dest)
	}

	// Large order, requires human approval.
	// We return a pending status and a JobID. The agent will need to handle this.
	jobID := fmt.Sprintf("SHIPPING-JOB-%d", time.Now().UnixNano())
	return ShippingOrderStatus{
		Status:        "pending",
		NumContainers: args.NumContainers,
		Dest:          args.Dest,
		Message:       fmt.Sprintf("Large order of %d containers to %s requires manual approval.", args.NumContainers, args.Dest),
		JobID:         jobID,
	}
}

func autoApproveSmallOrder(numContainer int, dest string) ShippingOrderStatus {
	return ShippingOrderStatus{
		Status:        "APPROVED",
		OrderID:       fmt.Sprintf("ORD-%d-AUTO", numContainer),
		NumContainers: numContainer,
		Dest:          dest,
		Message:       fmt.Sprintf("Order auto-approved: %d containers to %s", numContainer, dest),
	}
}

func main() {
	shippingTool, err := functiontool.New(functiontool.Config{
		Name:          "placeShippingOrder",
		Description:   "Places a shipping order. Small orders are auto-approved. Large orders require manual approval.",
		IsLongRunning: true,
	}, placeShippingOrder)
	if err != nil {
		log.Fatal(err)
	}

	shippingAgent, err := llmagent.New(llmagent.Config{
		Name:  "ShippingAgent",
		Model: mdl.FromEnv(),
		Instruction: `You are a shipping order assistant.
- Use the 'placeShippingOrder' tool to place shipping orders.
- Report the final order status to the user.
- If the status is pending, inform the user that the order requires manual approval and provide them with the JobID. The user will then need to use a separate mechanism to approve or deny the order.
`,
		Tools: []tool.Tool{shippingTool},
	})
	if err != nil {
		log.Fatal(err)
	}

	runner.Run(shippingAgent)
}
