package api

type Error struct {
	Status    int    `json:"status"`
	ErrorCode int    `json:"error-code"`
	Reason    string `json:"description"`
}

type CreateBucketRequest struct {
	Name string `json:"name"`
}

type BucketResponse struct {
	Name string `json:"name"`
}
