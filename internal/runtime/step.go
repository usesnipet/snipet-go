package runtime

type StepResult int

const (
	StepContinue StepResult = iota
	StepCancel
	StepMaxTurnsReached
	StepFinish
)
