package domain

type Wallet struct {
	Seed    string `json:"seed,omitempty"`
	Address string `json:"address"`
	Index   int    `json:"index,omitempty"` // Index of the address in the wallet
}
