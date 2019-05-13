package os

import "os/exec"

type OsCommands struct {
	Command string
	Args    []string
}

func (oc *OsCommands) Execute(b *[]byte) error {
	cmd := createCommnad(oc.Command, oc.Args...)
	result, err := cmd.Output()
	if err != nil {
		return err
	}
	b = &result
	return nil
}

func (oc *OsCommands) ExecuteWithOutput() error {
	cmd := createCommnad(oc.Command, oc.Args...)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func createCommnad(command string, args ...string) *exec.Cmd {
	cmd := exec.Command(command, args...)
	return cmd
}
