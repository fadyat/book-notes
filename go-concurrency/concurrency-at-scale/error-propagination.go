package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
)

type MyError struct {

	// Storing low-level error, we always want to be able to get back
	Inner   error
	Message string

	// taking note of a stack trace, when the error was created
	StackTrace string

	// storage for catch-all miscellaneous information,
	//  e.g a hash of the stack trace
	Misc map[string]interface{}
}

func (e MyError) Error() string {
	return e.Message
}

func wrapError(err error, messagef string, msgArgs ...interface{}) MyError {
	return MyError{
		Inner:      err,
		Message:    fmt.Sprintf(messagef, msgArgs...),
		StackTrace: string(debug.Stack()),
		Misc:       make(map[string]interface{}),
	}
}

type LowLevelErr struct {
	error
}

func isGloballyExec(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		// here we wrap the raw error with customized error
		// in that case we ok with message and don't mask it
		return false, LowLevelErr{wrapError(err, err.Error())}
	}

	return info.Mode().Perm()&0100 == 0100, nil
}

type IntermediateErr struct {
	error
}

func runJob(id string) error {
	const jobBinPath = "/bad/job/bin"

	isExecutable, err := isGloballyExec(jobBinPath)
	if err != nil {
		// here we are customizing the error with a crafted message.
		// we want to obfuscate the low-level details of why the job isn’t running
		// because we feel it’s not important information to consumers of our module.
		return IntermediateErr{wrapError(
			err,
			"can't run job %q: requisite binaries are not available",
			id,
		)}
	}

	if !isExecutable {
		return wrapError(
			nil,
			"can't run job %q: requisite binaries are not executable",
			id,
		)
	}

	return exec.Command(jobBinPath, "--id"+id).Run()
}

func handleError(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logID: %v", key))

	// log out the full error in case someone needs to dig into what happend
	log.Printf("%#v", err)
	log.Printf("[%v] %v", key, message)
}

func main() {
	// all functions are written in one file,
	//	but they can be broken into separate modules

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	if err := runJob("1"); err != nil {
		msg := "There was an unexpected issue; please report this as a bug."

		// checking if is there a well-crafted error, and we can pass it to user
		if _, ok := err.(IntermediateErr); ok {
			msg = err.Error()
		}

		// bind the log and the error message together with an ID of 1, can use guid
		handleError(1, err, msg)
	}
}
