package process

import "context"

// ============================================================================
// PROCESS INTERFACE (Concurrent — goroutine-based)
// ============================================================================

// Process represents a participant in a distributed computation.
//
// PARALLEL VERSION:
//
//	Instead of the synchronous Handle(msg) []Message interface,
//	each process runs as a long-lived goroutine. It reads messages
//	from its inbox channel and sends outgoing messages via a send function.
//
// From Cachin et al., Section 2.1:
//
//	"We consider a distributed system as a collection of processes
//	 that communicate by exchanging messages."
//
// The Run method is launched in its own goroutine by the runtime.
// It must return when ctx is cancelled.
type Process interface {
	// ID returns the unique identifier of this process.
	ID() ProcessID

	// Run is executed as a goroutine. It should:
	//   - Read messages from inbox
	//   - Call send(msg) to emit outgoing messages
	//   - Return when ctx.Done() fires
	Run(ctx context.Context, inbox <-chan Message, send func(Message))
}

// ============================================================================
// HANDLER INTERFACE (Synchronous adapter)
// ============================================================================

// Handler is the synchronous message-handler interface, identical to the
// original single-threaded Process interface. It processes one message
// at a time and returns outgoing messages.
//
// Use WrapHandler to convert a Handler into a concurrent Process.
type Handler interface {
	// ID returns the unique identifier of this process.
	ID() ProcessID

	// Handle processes a single incoming message and returns
	// zero or more outgoing messages (one atomic step).
	Handle(msg Message) []Message
}

// WrapHandler adapts a synchronous Handler into a concurrent Process.
// The resulting Process reads from its inbox in a loop and calls
// h.Handle for each message, forwarding replies via send.
func WrapHandler(h Handler) Process {
	return &handlerAdapter{handler: h}
}

// handlerAdapter bridges Handler → Process.
type handlerAdapter struct {
	handler Handler
}

func (a *handlerAdapter) ID() ProcessID { return a.handler.ID() }

func (a *handlerAdapter) Run(ctx context.Context, inbox <-chan Message, send func(Message)) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-inbox:
			if !ok {
				return // inbox closed
			}
			replies := a.handler.Handle(msg)
			for _, r := range replies {
				send(r)
			}
		}
	}
}
