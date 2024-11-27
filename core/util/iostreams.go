package util

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/mattn/go-isatty"
)

type IOStreams struct {
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer

	stdinIsTTY  bool
	stdoutIsTTY bool
	stderrIsTTY bool

	colorEnabled             bool
	progressMu               sync.Mutex
	progressActive           *spinner.Spinner
	progressIndicatorEnabled bool
}

// NewIOStreams initializes IOStreams with default terminal settings.
func NewIOStreams() *IOStreams {
	stdinIsTTY := isTerminal(os.Stdin)
	stdoutIsTTY := isTerminal(os.Stdout)
	stderrIsTTY := isTerminal(os.Stderr)
	progressIndicator := stdoutIsTTY && stderrIsTTY

	return &IOStreams{
		In:                       os.Stdin,
		Out:                      os.Stdout,
		ErrOut:                   os.Stderr,
		stdinIsTTY:               stdinIsTTY,
		stdoutIsTTY:              stdoutIsTTY,
		stderrIsTTY:              stderrIsTTY,
		colorEnabled:             stdoutIsTTY,
		progressIndicatorEnabled: progressIndicator,
	}
}

func isTerminal(f *os.File) bool {
	return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
}

func (s *IOStreams) IsStdinTTY() bool  { return s.stdinIsTTY }
func (s *IOStreams) IsStdoutTTY() bool { return s.stdoutIsTTY }
func (s *IOStreams) IsStderrTTY() bool { return s.stderrIsTTY }

func (s *IOStreams) ColorEnabled() bool { return s.colorEnabled }
func (s *IOStreams) SetColorEnabled(enabled bool) {
	s.colorEnabled = enabled
}

// Progress Indicator
func (s *IOStreams) StartProgressIndicator(label string) {
	if !s.progressIndicatorEnabled {
		return
	}
	s.progressMu.Lock()
	defer s.progressMu.Unlock()

	if s.progressActive != nil {
		s.progressActive.Prefix = label + " "
		return
	}

	// List of CharSets can be found here -> https://github.com/briandowns/spinner#available-character-sets
	sp := spinner.New(spinner.CharSets[27], 100*time.Millisecond, spinner.WithWriter(s.ErrOut), spinner.WithColor("fgCyan"))
	if label != "" {
		sp.Prefix = label + " "
	}
	sp.Start()
	s.progressActive = sp
}

func (s *IOStreams) StopProgressIndicator() {
	s.progressMu.Lock()
	defer s.progressMu.Unlock()

	if s.progressActive != nil {
		s.progressActive.Stop()
		s.progressActive = nil
	}
}

func (s *IOStreams) RunWithProgress(label string, task func() error) error {
	s.StartProgressIndicator(label)
	defer s.StopProgressIndicator()
	return task()
}
