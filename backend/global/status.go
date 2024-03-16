package global

type ExecutionResult uint

const (
	// iota begin from 0
	Pending ExecutionResult = iota
	PendingRejudge
	Compiling
	Running
	Accepted
	PresentationError
	WrongAnswer
	TimeLimitExceeded
	MemoryLimitExceeded
	OutputLimitExceeded
	RuntimeError
	CompileError
	UnknownError
)
