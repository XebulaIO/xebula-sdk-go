package network

type Network struct {
	Config Config
	Curve  any

	Wallet Wallet
}

type Wallet struct {
}

type INetwork interface {
	genereatePrivateKey() []uint8
}
