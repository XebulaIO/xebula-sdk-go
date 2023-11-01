package network

type Network struct {
	Config Config
	Curve  any

	Wallet Wallet
}

type INetwork interface {
	genereatePrivateKey() []uint8
}
