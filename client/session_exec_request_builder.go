package client

import (
	"fmt"
	"strings"

	"github.com/golang-devops/go-psexec/shared/dtos"
)

type SessionExecRequestBuilderBase interface {
	Winshell() SessionExecRequestBuilder
	Bash() SessionExecRequestBuilder
	None() SessionExecRequestBuilder
}

type SessionExecRequestBuilder interface {
	Exe(exe string) SessionExecRequestBuilder
	Args(args ...string) SessionExecRequestBuilder
	WorkingDir(workingDir string) SessionExecRequestBuilder
	Detached() SessionExecRequestBuilder

	BuildAndDoRequest() (*ExecResponse, error)
}

func NewSessionExecRequestBuilderBase(session *session) SessionExecRequestBuilderBase {
	return &sessionExecRequestBuilder{session: session, dto: &dtos.ExecDto{}}
}

type sessionExecRequestBuilder struct {
	session *session

	dto      *dtos.ExecDto
	detached bool
}

func (s *sessionExecRequestBuilder) Winshell() SessionExecRequestBuilder {
	s.dto.Executor = "winshell"
	return s
}

func (s *sessionExecRequestBuilder) Bash() SessionExecRequestBuilder {
	s.dto.Executor = "bash"
	return s
}

func (s *sessionExecRequestBuilder) None() SessionExecRequestBuilder {
	s.dto.Executor = "none"
	return s
}

func (s *sessionExecRequestBuilder) Exe(exe string) SessionExecRequestBuilder {
	s.dto.Exe = exe
	return s
}

func (s *sessionExecRequestBuilder) Args(args ...string) SessionExecRequestBuilder {
	s.dto.Args = args
	return s
}

func (s *sessionExecRequestBuilder) WorkingDir(workingDir string) SessionExecRequestBuilder {
	s.dto.WorkingDir = workingDir
	return s
}

func (s *sessionExecRequestBuilder) Detached() SessionExecRequestBuilder {
	s.detached = true
	return s
}

func (s *sessionExecRequestBuilder) BuildAndDoRequest() (*ExecResponse, error) {
	if strings.TrimSpace(s.dto.Executor) == "" {
		panic("SessionExecRequestBuilder requires Executor")
	}
	if strings.TrimSpace(s.dto.Exe) == "" {
		panic("SessionExecRequestBuilder requires Exe")
	}

	var relUrl string
	if s.detached {
		relUrl = "/auth/start"
	} else {
		relUrl = "/auth/stream"
	}

	resp, err := s.session.StreamEncryptedJsonRequest(relUrl, s.dto)
	if err != nil {
		return nil, fmt.Errorf("Unable make POST request to url '%s', error: %s", relUrl, err.Error())
	}

	return resp, nil
}
