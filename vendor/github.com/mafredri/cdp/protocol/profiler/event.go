// Code generated by cdpgen. DO NOT EDIT.

package profiler

import (
	"github.com/mafredri/cdp/protocol/debugger"
	"github.com/mafredri/cdp/rpcc"
)

// ConsoleProfileStartedClient is a client for ConsoleProfileStarted events. Sent when new profile recording is started using console.profile() call.
type ConsoleProfileStartedClient interface {
	// Recv calls RecvMsg on rpcc.Stream, blocks until the event is
	// triggered, context canceled or connection closed.
	Recv() (*ConsoleProfileStartedReply, error)
	rpcc.Stream
}

// ConsoleProfileStartedReply is the reply for ConsoleProfileStarted events.
type ConsoleProfileStartedReply struct {
	ID       string            `json:"id"`              // No description.
	Location debugger.Location `json:"location"`        // Location of console.profile().
	Title    *string           `json:"title,omitempty"` // Profile title passed as an argument to console.profile().
}

// ConsoleProfileFinishedClient is a client for ConsoleProfileFinished events.
type ConsoleProfileFinishedClient interface {
	// Recv calls RecvMsg on rpcc.Stream, blocks until the event is
	// triggered, context canceled or connection closed.
	Recv() (*ConsoleProfileFinishedReply, error)
	rpcc.Stream
}

// ConsoleProfileFinishedReply is the reply for ConsoleProfileFinished events.
type ConsoleProfileFinishedReply struct {
	ID       string            `json:"id"`              // No description.
	Location debugger.Location `json:"location"`        // Location of console.profileEnd().
	Profile  Profile           `json:"profile"`         // No description.
	Title    *string           `json:"title,omitempty"` // Profile title passed as an argument to console.profile().
}
