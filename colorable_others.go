//go:build !windows
// +build !windows

package colorable

import (
	"io"
	"os"

	"github.com/mattn/go-isatty"
)

// NewColorable returns new instance of Writer which handles escape sequence.
func NewColorable(file io.Writer) io.Writer {
	if file == nil {
		panic("nil passed instead of io.Writer to NewColorable()")
	}

	return file
}

// NewColorableStdout returns new instance of Writer which handles escape sequence for stdout.
func NewColorableStdout() io.Writer {
	return os.Stdout
}

// NewColorableStderr returns new instance of Writer which handles escape sequence for stderr.
func NewColorableStderr() io.Writer {
	return os.Stderr
}

// EnableColorsStdout enable colors if possible.
func EnableColorsStdout(enabled *bool) func() {
	if enabled != nil {
		*enabled = true
	}
	return func() {}
}

type filelike interface {
	Fd() uintptr
}

func IsWriterTerminal(file io.Writer) bool {
	if f, ok := file.(filelike); ok {
		return IsTerminal(f.Fd())
	}
	return false
}

func IsReaderTerminal(file io.Reader) bool {
	if f, ok := file.(filelike); ok {
		return IsTerminal(f.Fd())
	}
	return false
}

func IsTerminal(fd uintptr) bool {
	return isatty.IsTerminal(fd)
}

func IsCygwinTerminal(fd uintptr) bool {
	return isatty.IsCygwinTerminal(fd)
}
