package autoupdate

import (
	"os/exec"
	"os"
	"fmt"
)

type command struct {
	*exec.Cmd
	wasStoppedByUs bool
}

func createCommand(commandName string) *command {
	command := &command{
		Cmd:            exec.Command(commandName),
		wasStoppedByUs: false,
	}
	command.Cmd.Stdout = os.Stdout
	command.Cmd.Stderr = os.Stderr

	return command
}

func (c *command) stop() {
	c.wasStoppedByUs = true
	c.Process.Kill()
	c.Process.Wait()
}

func (c *command) listenForStop() {
	go func() {
		cmd.Wait()
		if !cmd.wasStoppedByUs && cmd.ProcessState != nil && cmd.ProcessState.Exited() {
			exitStatus := 1
			if cmd.ProcessState.Success() {
				exitStatus = 0
			}
			fmt.Println("Closing with", exitStatus)
			os.Exit(exitStatus)
		}
	}()

}
