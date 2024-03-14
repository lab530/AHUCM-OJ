package global

type ExecutionResult struct {
	Halt     bool
	Accepted struct {
		Value1 int
		Value2 uint32
		Value3 uint32
	}
	WrongAnswer struct {
		Value1 int
		Value2 int
	}
	PresentationError bool
	CompilationError  string
	RuntimeError      string
	TimeLimitExceeded struct {
		Value1 uint32
		Value2 uint32
	}
	MemoLimitExceeded struct {
		Value1 uint32
		Value2 uint32
	}
	UnknownError string
}
