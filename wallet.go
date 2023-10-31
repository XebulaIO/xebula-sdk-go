package xebula

import "xebula/network"

type Wallet struct {
	Config  network.Config
	Network network.INetwork
}

func (w *Wallet) Create() {

}

func NewWallet(n network.INetwork) *Wallet {
	return &Wallet{
		Network: n,
	}
}
