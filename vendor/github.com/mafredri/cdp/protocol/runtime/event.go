// Code generated by cdpgen. DO NOT EDIT.

package runtime

import (
	"encoding/json"

	"github.com/mafredri/cdp/rpcc"
)

// ExecutionContextCreatedClient is a client for ExecutionContextCreated events. Issued when new execution context is created.
type ExecutionContextCreatedClient interface {
	// Recv calls RecvMsg on rpcc.Stream, blocks until the event is
	// triggered, context canceled or connection closed.
	Recv() (*ExecutionContextCreatedReply, error)
	rpcc.Stream
}

// ExecutionContextCreatedReply is the reply for ExecutionContextCreated events.
type ExecutionContextCreatedReply struct {
	Context ExecutionContextDescription `json:"context"` // A newly created execution context.
}

// ExecutionContextDestroyedClient is a client for ExecutionContextDestroyed events. Issued when execution context is destroyed.
type ExecutionContextDestroyedClient interface {
	// Recv calls RecvMsg on rpcc.Stream, blocks until the event is
	// triggered, context canceled or connection closed.
	Recv() (*ExecutionContextDestroyedReply, error)
	rpcc.Stream
}

// ExecutionContextDestroyedReply is the reply for ExecutionContextDestroyed events.
type ExecutionContextDestroyedReply struct {
	ExecutionContextID ExecutionContextID `json:"executionContextId"` // Id of the destroyed context
}

// ExecutionContextsClearedClient is a client for ExecutionContextsCleared events. Issued when all executionContexts were cleared in browser
type ExecutionContextsClearedClient interface {
	// Recv calls RecvMsg on rpcc.Stream, blocks until the event is
	// triggered, context canceled or connection closed.
	Recv() (*ExecutionContextsClearedReply, error)
	rpcc.Stream
}

// ExecutionContextsClearedReply is the reply for ExecutionContextsCleared events.
type ExecutionContextsClearedReply struct{}

// ExceptionThrownClient is a client for ExceptionThrown events. Issued when exception was thrown and unhandled.
type ExceptionThrownClient interface {
	// Recv calls RecvMsg on rpcc.Stream, blocks until the event is
	// triggered, context canceled or connection closed.
	Recv() (*ExceptionThrownReply, error)
	rpcc.Stream
}

// ExceptionThrownReply is the reply for ExceptionThrown events.
type ExceptionThrownReply struct {
	Timestamp        Timestamp        `json:"timestamp"`        // Timestamp of the exception.
	ExceptionDetails ExceptionDetails `json:"exceptionDetails"` // No description.
}

// ExceptionRevokedClient is a client for ExceptionRevoked events. Issued when unhandled exception was revoked.
type ExceptionRevokedClient interface {
	// Recv calls RecvMsg on rpcc.Stream, blocks until the event is
	// triggered, context canceled or connection closed.
	Recv() (*ExceptionRevokedReply, error)
	rpcc.Stream
}

// ExceptionRevokedReply is the reply for ExceptionRevoked events.
type ExceptionRevokedReply struct {
	Reason      string `json:"reason"`      // Reason describing why exception was revoked.
	ExceptionID int    `json:"exceptionId"` // The id of revoked exception, as reported in exceptionUnhandled.
}

// ConsoleAPICalledClient is a client for ConsoleAPICalled events. Issued when console API was called.
type ConsoleAPICalledClient interface {
	// Recv calls RecvMsg on rpcc.Stream, blocks until the event is
	// triggered, context canceled or connection closed.
	Recv() (*ConsoleAPICalledReply, error)
	rpcc.Stream
}

// ConsoleAPICalledReply is the reply for ConsoleAPICalled events.
type ConsoleAPICalledReply struct { // Type Type of the call.
	//
	// Values: "log", "debug", "info", "error", "warning", "dir", "dirxml", "table", "trace", "clear", "startGroup", "startGroupCollapsed", "endGroup", "assert", "profile", "profileEnd", "count", "timeEnd".
	Type               string             `json:"type"`
	Args               []RemoteObject     `json:"args"`                 // Call arguments.
	ExecutionContextID ExecutionContextID `json:"executionContextId"`   // Identifier of the context where the call was made.
	Timestamp          Timestamp          `json:"timestamp"`            // Call timestamp.
	StackTrace         *StackTrace        `json:"stackTrace,omitempty"` // Stack trace captured when the call was made.
	// Context Console context descriptor for calls on non-default console context (not console.*): 'anonymous#unique-logger-id' for call on unnamed context, 'name#unique-logger-id' for call on named context.
	//
	// Note: This property is experimental.
	Context *string `json:"context,omitempty"`
}

// InspectRequestedClient is a client for InspectRequested events. Issued when object should be inspected (for example, as a result of inspect() command line API call).
type InspectRequestedClient interface {
	// Recv calls RecvMsg on rpcc.Stream, blocks until the event is
	// triggered, context canceled or connection closed.
	Recv() (*InspectRequestedReply, error)
	rpcc.Stream
}

// InspectRequestedReply is the reply for InspectRequested events.
type InspectRequestedReply struct {
	Object RemoteObject    `json:"object"` // No description.
	Hints  json.RawMessage `json:"hints"`  // No description.
}