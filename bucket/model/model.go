package model

type Bucket struct {
	Name string `json:"name"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
