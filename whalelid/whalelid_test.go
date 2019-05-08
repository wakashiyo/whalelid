package whalelid

import (
	"fmt"
	"testing"
	"unsafe"
)

type testCmd struct {
	command string
}

type testError struct {
	message string
}

func (te *testError) Error() string {
	return fmt.Sprintf("error %v", te.message)
}

func (tc *testCmd) Run() error {
	if tc.command == "" {
		return &testError{"testCmd command is empty"}
	}
	fmt.Println(tc.command)
	return nil
}

func (tc *testCmd) Output() ([]byte, error) {
	if tc.command == "" {
		return nil, &testError{"testCmd command is empty"}
	}
	return *(*[]byte)(unsafe.Pointer(&tc.command)), nil
}

func TestExecuteCommandOutputIsSuccess(t *testing.T) {
	tc := testCmd{
		command: "test1",
	}
	b := []byte{}
	err := output(&tc, &b)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Println(b)
	if len(b) == 0 {
		t.Fatalf("error length %#v", len(b))
	}
}

func TestExecuteCommandOutputIsFailed(t *testing.T) {
	tc := testCmd{
		command: "",
	}
	b := []byte{}
	err := output(&tc, &b)
	if err == nil {
		t.Fatal("TestExecuteCommandOutputIsFailed is error")
	}
	//fmt.Println(b)
	if len(b) != 0 {
		t.Fatalf("error length %#v", len(b))
	}
}
