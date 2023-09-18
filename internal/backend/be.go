package backend

import "github.com/fakeyanss/usg-go/internal/model"

type StorageBackend interface {
	CreateBucket(bucket string)
	DeleteBucket(bucket string)
	ListBuckets() []*model.UsgBucket
	HeadBucket(bucket string) bool

	ListObjects(bucket string, delimiter string, marker string, maxSize int, keyPrefix string)
}
