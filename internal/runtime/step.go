package runtime

// StepResult is the outcome of a single engine turn (see Engine.step). It is
// only meaningful when step returns a nil error: on error the caller inspects
// the error and ignores the result, which will be StepResultInvalid.
type StepResult string

const (
	// StepResultInvalid: the step result is invalid, the error is not nil.
	StepResultInvalid StepResult = "invalid"
	// StepContinue: the turn produced output, run another turn.
	StepContinue StepResult = "continue"
	// StepFinish: the assistant produced a final message, finish the execution.
	StepFinish StepResult = "finish"
	// StepCancel: the context was cancelled, stop the execution.
	StepCancel StepResult = "cancel"
	// StepMaxTurnsReached: the turn budget was exhausted before this turn ran.
	StepMaxTurnsReached StepResult = "max_turns_reached"
)
