package receipt

import (
	"strings"
	"testing"
)

func fixtureReceipt() Receipt {
	return Receipt{
		Schema:        "jev-lab/v1",
		WorkloadID:    "spiffe://example.test/ns/lab/sa/jev",
		Model:         "jev-latest",
		RequestDigest: digest('a'),
		StateDigest:   digest('b'),
		Questions: []Question{
			{ID: "route", Kind: KindChoice, Choices: []string{"accept", "review"}},
			{ID: "urgent", Kind: KindNoul},
		},
		Signals: []Signal{
			{QuestionID: "route", Kind: KindChoice, Probabilities: map[string]float64{"accept": 0.8, "review": 0.2}, Confidence: 0.9, ConfidenceMethod: ConfidenceMethodCalibrated},
			{QuestionID: "urgent", Kind: KindNoul, Value: 0.7, Confidence: 0.8, ConfidenceMethod: ConfidenceMethodMaxProbability},
		},
		Evidence: Evidence{
			DeclarationDigest:        digest('c'),
			IRDigest:                 digest('d'),
			GenerationDigest:         digest('e'),
			ReverseObservationDigest: digest('f'),
		},
	}
}

func digest(ch byte) string { return "sha256:" + strings.Repeat(string(ch), 64) }

func TestReceiptValidationAndEvidenceDigest(t *testing.T) {
	receipt := fixtureReceipt()
	if err := receipt.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	first, err := receipt.EvidenceDigest()
	if err != nil {
		t.Fatalf("EvidenceDigest() error = %v", err)
	}
	second, err := receipt.EvidenceDigest()
	if err != nil || first != second {
		t.Fatalf("EvidenceDigest() is not stable: %q != %q, err=%v", first, second, err)
	}
	if got := receipt.CanonicalQuestionIDs(); len(got) != 2 || got[0] != "route" || got[1] != "urgent" {
		t.Fatalf("CanonicalQuestionIDs() = %v", got)
	}
}

func TestReceiptRejectsTampering(t *testing.T) {
	receipt := fixtureReceipt()
	receipt.Evidence.IRDigest = "sha256:tampered"
	if err := receipt.Validate(); err == nil {
		t.Fatal("Validate() accepted a malformed digest")
	}
	receipt = fixtureReceipt()
	receipt.Signals[0].Probabilities["accept"] = 0.95
	if err := receipt.Validate(); err == nil {
		t.Fatal("Validate() accepted a distribution that does not sum to one")
	}
	receipt = fixtureReceipt()
	receipt.Signals[0].ConfidenceMethod = "provider_guess"
	if err := receipt.Validate(); err == nil {
		t.Fatal("Validate() accepted an unknown confidence method")
	}
}

func TestRouteByConfidenceNeverAuthorizes(t *testing.T) {
	if got := RouteByConfidence(0.95, 0.9); got != RouteAccept {
		t.Fatalf("RouteByConfidence() = %q, want %q", got, RouteAccept)
	}
	if got := RouteByConfidence(0.4, 0.9); got != RouteReview {
		t.Fatalf("RouteByConfidence() = %q, want %q", got, RouteReview)
	}
	if got := RouteByConfidence(0.95, 2); got != RouteReview {
		t.Fatalf("RouteByConfidence() invalid threshold = %q, want %q", got, RouteReview)
	}
}

func TestRouteBySignalRequiresExplicitAllowedMethod(t *testing.T) {
	signal := Signal{Confidence: 0.95, ConfidenceMethod: ConfidenceMethodCalibrated}
	if got := RouteBySignal(signal, 0.9, ConfidenceMethodCalibrated); got != RouteAccept {
		t.Fatalf("RouteBySignal() = %q, want %q", got, RouteAccept)
	}
	if got := RouteBySignal(signal, 0.9, ConfidenceMethodMaxProbability); got != RouteReview {
		t.Fatalf("RouteBySignal() with disallowed method = %q, want %q", got, RouteReview)
	}
	signal.ConfidenceMethod = ConfidenceMethodUnspecified
	if got := RouteBySignal(signal, 0.9, ConfidenceMethodUnspecified); got != RouteReview {
		t.Fatalf("RouteBySignal() with unspecified method = %q, want %q", got, RouteReview)
	}
}
