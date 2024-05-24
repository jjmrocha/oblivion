package api

type Error struct {
	Status    int    `json:"status"`
	ErrorCode int    `json:"error-code"`
	Reason    string `json:"description"`
}
