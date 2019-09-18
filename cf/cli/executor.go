package cli

import (
	"os"
	"os/exec"
)

type Executor interface {
	Execute(arguments []string) (err error)
}

//HttpConnection
type executor struct {
	traceLogging bool
}

func NewExecutor(traceLogging bool) Executor {
	return executor{
		traceLogging: traceLogging,
	}
}

func (ec executor) Execute(arguments []string) (err error) {
	cfCmdToolPath, err := exec.LookPath("cf")
	if err != nil {
		return err
	}

	outChannel := os.NewFile(0, os.DevNull)
	if ec.traceLogging {
		outChannel = os.Stdout
	}

	cmd := exec.Cmd{
		Path:   cfCmdToolPath,
		Args:   append([]string{cfCmdToolPath}, arguments...),
		Stdout: outChannel,
		Stderr: os.Stderr,
	}

	return cmd.Run()
}
