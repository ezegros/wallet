package domain

type Wallet struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Seed   string `json:"seed"`
}
