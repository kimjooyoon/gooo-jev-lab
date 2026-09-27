package main

import (
	"fmt"

	"github.com/kimjooyoon/gooo-jev-lab/receipt"
)

func main() {
	signal := receipt.Signal{
		Confidence:       0.92,
		ConfidenceMethod: receipt.ConfidenceMethodCalibrated,
	}
	route := receipt.RouteBySignal(signal, 0.85, receipt.ConfidenceMethodCalibrated)
	fmt.Printf("support-triage observation route: %s\n", route)
}
