package cli

//Test struct for executro
type FakeExecutor struct {
	traceLogging bool
	counter int
	argumentsOutput map[int][]string
}

func (tx *FakeExecutor) NewFakeExecutor() CfExecutor {
	return tx
}

func (tx *FakeExecutor) Execute(arguments []string) (err error) {
	tx.counter++
	size := len(tx.argumentsOutput)
	if tx.argumentsOutput == nil {
		//init lazy - only for fake executor
		tx.argumentsOutput = map[int][]string{0: arguments}
	} else {
		tx.argumentsOutput[size] = arguments
	}
	return nil
}
func (tx *FakeExecutor) ExecutorCallCount() int {
	return tx.counter
}

func (tx *FakeExecutor) ExecutorArgumentsOutput() map[int][]string {
	return tx.argumentsOutput
}

