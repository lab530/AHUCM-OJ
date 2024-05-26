package global

type ExecutionResult uint

const (
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

// 为了方便前端理解，可以定义一个映射来将 ExecutionResult 转换为字符串描述
var ExecutionResultToString = map[ExecutionResult]string{
	Pending:             "Pending",
	PendingRejudge:      "Pending Rejudge",
	Compiling:           "Compiling",
	Running:             "Running",
	Accepted:            "Accepted",
	PresentationError:   "Presentation Error",
	WrongAnswer:         "Wrong Answer",
	TimeLimitExceeded:   "Time Limit Exceeded",
	MemoryLimitExceeded: "Memory Limit Exceeded",
	OutputLimitExceeded: "Output Limit Exceeded",
	RuntimeError:        "Runtime Error",
	CompileError:        "Compile Error",
	UnknownError:        "Unknown Error",
}
