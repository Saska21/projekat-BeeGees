package runtime

import "fmt"

// Checker verifies a property of the execution.
//
// From Cachin et al., Section 2.1:
//
// Distributed algorithm properties are classified as:
//
//	SAFETY: "Something bad NEVER happens."
//	  - Checked after every step.
//	  - If violated at any step, the property does not hold.
//	  - Example: "No message is delivered more than once." (PL2)
//
//	LIVENESS: "Something good EVENTUALLY happens."
//	  - Checked at the END of execution.
//	  - Requires the execution to be "long enough" or "fair."
//	  - Example: "If p sends m to q, then q eventually delivers m." (PL1)
type Checker struct {
	Name  string
	Check func(events []Event) string // "" = OK, non-empty = violation
	AtEnd bool                        // true = liveness (check at end only)
}

// NoCreationChecker verifies no message is delivered without being sent first.
func NoCreationChecker() Checker {
	return Checker{
		Name:  "No creation",
		AtEnd: false,
		Check: func(events []Event) string {
			sent := make(map[string]bool)
			for _, e := range events {
				if e.Type == EventSend && e.Message != nil {
					sent[e.Message.String()] = true
				}
				if e.Type == EventDeliver && e.Message != nil {
					key := e.Message.String()
					if !sent[key] {
						return "No-creation violation: delivered message was never sent: " + key
					}
				}
			}
			return ""
		},
	}
}

// NoDuplicationChecker verifies no message is delivered more than once.
func NoDuplicationChecker() Checker {
	return Checker{
		Name:  "No duplication",
		AtEnd: false,
		Check: func(events []Event) string {
			delivered := make(map[string]int)
			for _, e := range events {
				if e.Type == EventDeliver && e.Message != nil {
					key := e.Message.String()
					delivered[key]++
					if delivered[key] > 1 {
						return fmt.Sprintf("No-duplication violation: message delivered %d times: %s", delivered[key], key)
					}
				}
			}
			return ""
		},
	}
}

// ReliableDeliveryChecker verifies all sent messages are eventually delivered.
// This is a LIVENESS property, checked only at the end.
func ReliableDeliveryChecker() Checker {
	return Checker{
		Name:  "Reliable delivery",
		AtEnd: true,
		Check: func(events []Event) string {
			sent := make(map[string]bool)
			delivered := make(map[string]bool)
			for _, e := range events {
				if e.Type == EventSend && e.Message != nil {
					sent[e.Message.String()] = true
				}
				if e.Type == EventDeliver && e.Message != nil {
					delivered[e.Message.String()] = true
				}
			}
			for key := range sent {
				if !delivered[key] {
					return "Reliable-delivery concern: sent message not delivered: " + key
				}
			}
			return ""
		},
	}
}
