package domain

type Wallet struct {
	UserID int64  `json:"user_id"`
	Seed   string `json:"seed"`
}
