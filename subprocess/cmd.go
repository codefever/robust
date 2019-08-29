package subprocess

import (
	"bufio"
	"context"
	"io"
	"log"
	"os/exec"
)

func forwardLinesToLog(name string, reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		log.Printf("%v=> %v", name, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Fail to read %v: %v", name, err)
	}
}

// RunCommand runs the provided command and returns its error channel with its
// cancellation function.
func RunCommand(command string) (*exec.Cmd, <-chan error, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "sh", "-c", "exec "+command)
	maybeSetUpPdeathsig(cmd)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	errc := make(chan error, 1)
	go func() {
		errc <- cmd.Wait()
		close(errc)
	}()

	go forwardLinesToLog("STDOUT", stdout)
	go forwardLinesToLog("STDINT", stderr)

	return cmd, errc, cancel
}
