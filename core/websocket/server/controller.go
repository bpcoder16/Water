package server

import "context"

type Controller interface {
	Init(ctx context.Context)
	Prepare(ctx context.Context)
	Finish(ctx context.Context)
	Begin(ctx context.Context)
	End(ctx context.Context)
	Merge(ctx context.Context, base Controller)

	DoNotLog(ctx context.Context)
	DoNotLogRequest(ctx context.Context)
	DoNotLogResponse(ctx context.Context)

	SampleVariable(ctx context.Context, name string) (string, bool)

	// The controller NEEDN'T implement the following interface methods, they SHOULD
	// embedded the BaseEventController.
	Respond(ctx context.Context)
	SetEventType(ctx context.Context, typ string)

	// The controller SHOULD implement the following interface methods
	ParseEvent(ctx context.Context, eventStr []byte) (interface{}, error)
	ProcessEvent(ctx context.Context)
}
