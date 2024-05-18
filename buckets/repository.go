package buckets

type Repository interface {
	CreateBucket(string) (*Bucket, error)
	GetBucket(string) (*Bucket, error)
	//DropBucket(string) error
	//GetAllBuckets() ([]*Bucket, error)
	//Store(*Bucket, string, any) error
	//Read(*Bucket, string) (any, error)
	//Delete(*Bucket, string) error
}
