package global

type ExecutionResult uint

const (
	Pending ExecutionResult = iota
	PendinRejudge
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
