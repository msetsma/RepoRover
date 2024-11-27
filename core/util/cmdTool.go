package util

import (
	"github.com/msetsma/RepoRover/core/config"
)

type CmdTool struct {
	IOStreams *IOStreams
	Config    func() (*config.Manifest, error)
}

func NewCmdTool() *CmdTool {
	cfg := func() (*config.Manifest, error) {
		return config.Load()
	}

	// At some point we might need to use the cfg to generate the io streams.
	io := NewIOStreams()

	return &CmdTool{
		IOStreams: io,
		Config:    cfg,
	}
}
