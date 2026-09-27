package main

import (
	"fmt"

	"github.com/kimjooyoon/gooo-jev-lab/receipt"
)

func main() {
	route := receipt.RouteByConfidence(0.92, 0.85)
	fmt.Printf("support-triage observation route: %s\n", route)
}