package buckets

type Repository interface {
	GetAllBuckets() ([]*Bucket, error)
	CreateBucket(string) (*Bucket, error)
	GetBucket(string) (*Bucket, error)
	DropBucket(string) error
	//Store(*Bucket, string, any) error
	//Read(*Bucket, string) (any, error)
	//Delete(*Bucket, string) error
}
