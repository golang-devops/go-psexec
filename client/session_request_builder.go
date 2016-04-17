package client

import (
	"fmt"
	"strings"

	"github.com/golang-devops/go-psexec/shared"
)

type SessionRequestBuilderBase interface {
	Winshell() SessionRequestBuilder
	Bash() SessionRequestBuilder
}

type SessionRequestBuilder interface {
	Exe(exe string) SessionRequestBuilder
	Args(args ...string) SessionRequestBuilder
	WorkingDir(workingDir string) SessionRequestBuilder
	Detached() SessionRequestBuilder

	BuildAndDoRequest() (*RequestResponse, error)
}

func NewSessionRequestBuilder(session *Session) SessionRequestBuilderBase {
	return &sessionRequestBuilder{session: session, dto: &shared.ExecDto{}}
}

type sessionRequestBuilder struct {
	session *Session

	dto      *shared.ExecDto
	detached bool
}

func (s *sessionRequestBuilder) Winshell() SessionRequestBuilder {
	s.dto.Executor = "winshell"
	return s
}

func (s *sessionRequestBuilder) Bash() SessionRequestBuilder {
	s.dto.Executor = "bash"
	return s
}

func (s *sessionRequestBuilder) Exe(exe string) SessionRequestBuilder {
	s.dto.Exe = exe
	return s
}

func (s *sessionRequestBuilder) Args(args ...string) SessionRequestBuilder {
	s.dto.Args = args
	return s
}

func (s *sessionRequestBuilder) WorkingDir(workingDir string) SessionRequestBuilder {
	s.dto.WorkingDir = workingDir
	return s
}

func (s *sessionRequestBuilder) Detached() SessionRequestBuilder {
	s.detached = true
	return s
}

func (s *sessionRequestBuilder) BuildAndDoRequest() (*RequestResponse, error) {
	if strings.TrimSpace(s.dto.Executor) == "" {
		panic("SessionRequestBuilder requires Executor")
	}
	if strings.TrimSpace(s.dto.Exe) == "" {
		panic("SessionRequestBuilder requires Exe")
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
