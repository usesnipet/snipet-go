package tool

type Call struct {
	Key   string
	Input any
}

type Result struct {
	Key     string
	Success bool
	Output  any
	Error   error
}
