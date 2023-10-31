package xebula

import (
	"xebula/network"
)

type Xebula struct {
	config network.Config
	tron   network.TronNetwork
}

func NewXebula(cfg network.Config) *Xebula {
	return &Xebula{
		config: cfg,
		tron:   network.TronNetwork{Network: network.Network{Config: cfg}},
	}
}
