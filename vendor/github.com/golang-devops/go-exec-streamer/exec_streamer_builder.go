package execstreamer

import (
	"errors"
	"io"
)

//NewExecStreamerBuilder will create a builder to elegantly create a new ExecStreamer
func NewExecStreamerBuilder() ExecStreamerBuilder {
	return &execStreamerBuilder{&execStreamer{}}
}

//ExecStreamerBuilder is the builder interface
type ExecStreamerBuilder interface {
	ExecutorName(executorName string) ExecStreamerBuilder
	Exe(exe string) ExecStreamerBuilder
	Args(args ...string) ExecStreamerBuilder
	Dir(dir string) ExecStreamerBuilder
	Env(env ...string) ExecStreamerBuilder

	Writers(writers io.Writer) ExecStreamerBuilder

	StdoutWriter(writer io.Writer) ExecStreamerBuilder
	StdoutPrefix(prefix string) ExecStreamerBuilder

	StderrWriter(writer io.Writer) ExecStreamerBuilder
	StderrPrefix(prefix string) ExecStreamerBuilder

	AutoFlush() ExecStreamerBuilder

	DebugInfo(s string) ExecStreamerBuilder

	Build() (ExecStreamer, error)
}

type execStreamerBuilder struct {
	d *execStreamer
}

//ExecutorName sets the ExecutorName
func (e *execStreamerBuilder) ExecutorName(executorName string) ExecStreamerBuilder {
	e.d.ExecutorName = executorName
	return e
}

//Exe sets the Exe
func (e *execStreamerBuilder) Exe(exe string) ExecStreamerBuilder {
	e.d.Exe = exe
	return e
}

//Args sets the Args
func (e *execStreamerBuilder) Args(args ...string) ExecStreamerBuilder {
	e.d.Args = args
	return e
}

//Dir sets the Dir
func (e *execStreamerBuilder) Dir(dir string) ExecStreamerBuilder {
	e.d.Dir = dir
	return e
}

//Env sets the Env
func (e *execStreamerBuilder) Env(env ...string) ExecStreamerBuilder {
	e.d.Env = env
	return e
}

//Writers sets the Writers
func (e *execStreamerBuilder) Writers(writers io.Writer) ExecStreamerBuilder {
	e.d.StdoutWriter = writers
	e.d.StderrWriter = writers
	return e
}

//StdoutWriter sets the StdoutWriter
func (e *execStreamerBuilder) StdoutWriter(writer io.Writer) ExecStreamerBuilder {
	e.d.StdoutWriter = writer
	return e
}

//StdoutPrefix sets the StdoutPrefix
func (e *execStreamerBuilder) StdoutPrefix(prefix string) ExecStreamerBuilder {
	e.d.StdoutPrefix = prefix
	return e
}

//StderrWriter sets the StderrWriter
func (e *execStreamerBuilder) StderrWriter(writer io.Writer) ExecStreamerBuilder {
	e.d.StderrWriter = writer
	return e
}

//StderrPrefix sets the StderrPrefix
func (e *execStreamerBuilder) StderrPrefix(prefix string) ExecStreamerBuilder {
	e.d.StderrPrefix = prefix
	return e
}

//AutoFlush enables the AutoFlush
func (e *execStreamerBuilder) AutoFlush() ExecStreamerBuilder {
	e.d.AutoFlush = true
	return e
}

//DebugInfo adds debug info to be printed out when errors occur
func (e *execStreamerBuilder) DebugInfo(s string) ExecStreamerBuilder {
	e.d.DebugInfo = s
	return e
}

//Build will validate the set properties and return a ExecStreamer
func (e *execStreamerBuilder) Build() (ExecStreamer, error) {
	if e.d.ExecutorName == "" {
		return nil, errors.New("ExecStreamerBuilder requires ExecutorName to be non-empty")
	}
	if e.d.Exe == "" {
		return nil, errors.New("ExecStreamerBuilder requires Exe to be non-empty")
	}
	if e.d.StdoutWriter == nil {
		return nil, errors.New("ExecStreamerBuilder requires StdoutWriter to be set")
	}
	if e.d.StderrWriter == nil {
		return nil, errors.New("ExecStreamerBuilder requires StderrWriter to be set")
	}

	return e.d, nil
}
