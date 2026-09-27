package receipt

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
)

// Kind is the bounded decision primitive returned by a JEV-style provider.
type Kind string

const (
	KindChoice Kind = "choice"
	KindScore  Kind = "score"
	KindNoul   Kind = "noul"
)

// ConfidenceMethod identifies how a provider derived Signal.Confidence.
// An unspecified method is intentionally not eligible for automatic routing.
type ConfidenceMethod string

const (
	ConfidenceMethodUnspecified      ConfidenceMethod = "unspecified"
	ConfidenceMethodCalibrated       ConfidenceMethod = "calibrated"
	ConfidenceMethodMaxProbability   ConfidenceMethod = "max_probability"
	ConfidenceMethodTopTwoMargin     ConfidenceMethod = "top_two_margin"
	ConfidenceMethodOneMinusEntropy  ConfidenceMethod = "one_minus_entropy"
)

// Question defines the expected shape of one typed decision.
type Question struct {
	ID      string   `json:"id"`
	Kind    Kind     `json:"kind"`
	Choices []string `json:"choices,omitempty"`
}

// Signal is an observation, never an authorization.
type Signal struct {
	QuestionID    string             `json:"question_id"`
	Kind          Kind               `json:"kind"`
	Value         float64            `json:"value,omitempty"`
	Probabilities map[string]float64 `json:"probabilities,omitempty"`
	Confidence    float64            `json:"confidence"`
	ConfidenceMethod ConfidenceMethod `json:"confidence_method"`
}

// Evidence binds the receipt to the language pipeline that produced it.
type Evidence struct {
	DeclarationDigest        string `json:"declaration_digest"`
	IRDigest                 string `json:"ir_digest"`
	GenerationDigest         string `json:"generation_digest"`
	ReverseObservationDigest string `json:"reverse_observation_digest"`
}

// Receipt is the smallest portable unit exchanged across the lab boundary.
type Receipt struct {
	Schema        string     `json:"schema"`
	WorkloadID    string     `json:"workload_id"`
	Model         string     `json:"model"`
	RequestDigest string     `json:"request_digest"`
	StateDigest   string     `json:"state_digest"`
	Questions     []Question `json:"questions"`
	Signals       []Signal   `json:"signals"`
	Evidence      Evidence   `json:"evidence"`
}

// Route is an observation-derived suggestion. It has no side effect.
type Route string

const (
	RouteAccept Route = "ACCEPT"
	RouteReview Route = "REVIEW"
	RouteReject Route = "REJECT"
)

// Validate enforces the fail-closed boundary for a receipt.
func (r Receipt) Validate() error {
	if r.Schema == "" || r.WorkloadID == "" || r.Model == "" {
		return errors.New("schema, workload_id, and model are required")
	}
	for name, digest := range map[string]string{
		"request_digest":             r.RequestDigest,
		"state_digest":               r.StateDigest,
		"declaration_digest":         r.Evidence.DeclarationDigest,
		"ir_digest":                  r.Evidence.IRDigest,
		"generation_digest":          r.Evidence.GenerationDigest,
		"reverse_observation_digest": r.Evidence.ReverseObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("%s must be a sha256 digest", name)
		}
	}
	questions := make(map[string]Question, len(r.Questions))
	for _, question := range r.Questions {
		if question.ID == "" || questions[question.ID].ID != "" {
			return errors.New("question ids must be unique and non-empty")
		}
		if question.Kind != KindChoice && question.Kind != KindScore && question.Kind != KindNoul {
			return fmt.Errorf("question %q has unsupported kind", question.ID)
		}
		if question.Kind == KindChoice && (len(question.Choices) == 0 || len(question.Choices) > 255) {
			return fmt.Errorf("question %q must have 1 to 255 choices", question.ID)
		}
		questions[question.ID] = question
	}
	if len(r.Signals) != len(r.Questions) {
		return errors.New("signals must match questions exactly")
	}
	seen := make(map[string]struct{}, len(r.Signals))
	for _, signal := range r.Signals {
		question, ok := questions[signal.QuestionID]
		if !ok || question.Kind != signal.Kind {
			return fmt.Errorf("signal %q does not match a question", signal.QuestionID)
		}
		if _, ok := seen[signal.QuestionID]; ok {
			return fmt.Errorf("duplicate signal %q", signal.QuestionID)
		}
		seen[signal.QuestionID] = struct{}{}
		if !finiteUnit(signal.Confidence) {
			return fmt.Errorf("signal %q has invalid confidence", signal.QuestionID)
		}
		if !validConfidenceMethod(signal.ConfidenceMethod) {
			return fmt.Errorf("signal %q has unsupported confidence method", signal.QuestionID)
		}
		switch signal.Kind {
		case KindChoice:
			if err := validateChoice(question, signal); err != nil {
				return err
			}
		case KindScore, KindNoul:
			if !finiteUnit(signal.Value) {
				return fmt.Errorf("signal %q has invalid value", signal.QuestionID)
			}
		}
	}
	return nil
}

func validateChoice(question Question, signal Signal) error {
	if len(signal.Probabilities) != len(question.Choices) {
		return fmt.Errorf("signal %q probability count does not match choices", signal.QuestionID)
	}
	choiceSet := make(map[string]struct{}, len(question.Choices))
	total := 0.0
	for _, choice := range question.Choices {
		choiceSet[choice] = struct{}{}
		probability, ok := signal.Probabilities[choice]
		if !ok || !finiteUnit(probability) {
			return fmt.Errorf("signal %q has invalid probability for %q", signal.QuestionID, choice)
		}
		total += probability
	}
	for choice := range signal.Probabilities {
		if _, ok := choiceSet[choice]; !ok {
			return fmt.Errorf("signal %q has unknown choice %q", signal.QuestionID, choice)
		}
	}
	if math.Abs(total-1) > 1e-9 {
		return fmt.Errorf("signal %q probabilities must sum to one", signal.QuestionID)
	}
	return nil
}

func finiteUnit(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0) && value >= 0 && value <= 1
}

func validConfidenceMethod(method ConfidenceMethod) bool {
	switch method {
	case "", ConfidenceMethodUnspecified, ConfidenceMethodCalibrated, ConfidenceMethodMaxProbability, ConfidenceMethodTopTwoMargin, ConfidenceMethodOneMinusEntropy:
		return true
	default:
		return false
	}
}

func normalizedConfidenceMethod(method ConfidenceMethod) ConfidenceMethod {
	if method == "" {
		return ConfidenceMethodUnspecified
	}
	return method
}

func validDigest(value string) bool {
	const prefix = "sha256:"
	if len(value) != len(prefix)+sha256.Size*2 || value[:len(prefix)] != prefix {
		return false
	}
	_, err := hex.DecodeString(value[len(prefix):])
	return err == nil
}

// EvidenceDigest gives the pipeline binding a stable, content-addressed identity.
func (r Receipt) EvidenceDigest() (string, error) {
	if err := r.Validate(); err != nil {
		return "", err
	}
	payload, err := json.Marshal(r.Evidence)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(payload)
	return "sha256:" + hex.EncodeToString(digest[:]), nil
}

// CanonicalQuestionIDs returns deterministic question identity ordering for logs.
func (r Receipt) CanonicalQuestionIDs() []string {
	ids := make([]string, 0, len(r.Questions))
	for _, question := range r.Questions {
		ids = append(ids, question.ID)
	}
	sort.Strings(ids)
	return ids
}

// RouteByConfidence returns a reviewable observation and never authorizes an action.
func RouteByConfidence(confidence, threshold float64) Route {
	if !finiteUnit(confidence) || !finiteUnit(threshold) {
		return RouteReview
	}
	if confidence >= threshold {
		return RouteAccept
	}
	return RouteReview
}

// RouteBySignal returns an observation-derived route only when the confidence
// method is explicit and present in the caller's allowlist. It never authorizes
// execution or a side effect.
func RouteBySignal(signal Signal, threshold float64, allowedMethods ...ConfidenceMethod) Route {
	if !finiteUnit(signal.Confidence) || !finiteUnit(threshold) || !validConfidenceMethod(signal.ConfidenceMethod) {
		return RouteReview
	}
	method := normalizedConfidenceMethod(signal.ConfidenceMethod)
	if method == ConfidenceMethodUnspecified {
		return RouteReview
	}
	for _, allowed := range allowedMethods {
		if normalizedConfidenceMethod(allowed) == method {
			return RouteByConfidence(signal.Confidence, threshold)
		}
	}
	return RouteReview
}
