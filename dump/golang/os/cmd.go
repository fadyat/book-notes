package main

import (
	"fmt"
	"os/exec"
)

func mbDie(out []byte, err error) {
	switch e := err.(type) {
	case *exec.ExitError:
		fmt.Println("Exit code:", e.ExitCode())
	case *exec.Error:
		fmt.Println("Exec error:", e)
	case nil:
	default:
		fmt.Println("Unknown error:", e)
	}

	fmt.Print(string(out))
}

func main() {
	lsCommand := exec.Command("ls", "-l")
	out, err := lsCommand.Output()
	mbDie(out, err)

	dateCommand := exec.Command("date")
	out, err = dateCommand.Output()
	mbDie(out, err)

	llCommand := exec.Command("ll")
	out, err = llCommand.Output()
	mbDie(out, err)
}
