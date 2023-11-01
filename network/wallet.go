package network

type Wallet struct {
	Config  Config
	Network INetwork
}

func (w *Wallet) Create() {}

func NewWallet(n INetwork) *Wallet {
	return &Wallet{
		Network: n,
	}
}
