package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	err := runJob("1")
	if err != nil {
		msg := "This is a bug!"
		if _, ok := err.(IntermediateErr); ok {
			msg = err.Error()
		}
		handleError(1, err, msg)
	}
}

func handleError(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logID: %v]: ", key))
	log.Printf("%#v", err)
	fmt.Printf("[%v] %v\n", key, message)
}

type IntermediateErr struct {
	error
}

func runJob(id string) error {
	const jobBinPath = "/invalid/binary"
	isExecutable, err := lowLevelIsGloballyExec(jobBinPath)
	if err != nil {
		return IntermediateErr{wrapError(err, "cannot run job %q", id)}
	} else if !isExecutable {
		return wrapError(nil, "job binary is not executable")
	}
	return exec.Command(jobBinPath, "--id="+id).Run()
}

type LowLevelErr struct {
	error
}

func lowLevelIsGloballyExec(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, LowLevelErr{(wrapError(err, err.Error()))}
	}
	return info.Mode().Perm()&0100 == 0100, nil
}

type MyError struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]interface{}
}

func wrapError(err error, messagef string, msgArgs ...interface{}) MyError {
	return MyError{
		Inner:      err,
		Message:    fmt.Sprintf(messagef, msgArgs...),
		StackTrace: string(debug.Stack()),
		Misc:       make(map[string]interface{}),
	}
}

func (e MyError) Error() string {
	return e.Message
}
