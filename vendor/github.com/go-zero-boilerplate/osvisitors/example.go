package osvisitors

import (
	"fmt"
)

func printLinuxPathVariableCommand() {
	printPathVariableCommand(LinuxOs)
}

func printPathVariableCommand(osType OsType) {
	visitor := &exampleEchoPathVariableVisitor{}
	osType.Accept(visitor)
	fmt.Println(fmt.Sprintf("The echo command for osType '%#T' would be: %s", osType, visitor.answer))
}

type exampleEchoPathVariableVisitor struct{ answer string }

func (e *exampleEchoPathVariableVisitor) VisitWindows() { e.answer = "echo %PATH%" }
func (e *exampleEchoPathVariableVisitor) VisitLinux()   { e.answer = "echo $PATH" }
func (e *exampleEchoPathVariableVisitor) VisitDarwin()  { e.answer = "echo $PATH" }
