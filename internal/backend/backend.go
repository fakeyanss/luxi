package backend

import (
	"context"

	"github.com/fakeyanss/usg-go/internal/model"
)

type StorageBackend interface {
	CreateBucket(bucket string) *model.UsgResp
	// DeleteBucket(bucket string)
	ListBuckets(ctx context.Context) *model.UsgResp
	// HeadBucket(bucket string) bool

	ListObjects(ctx context.Context, bucket, delimiter, encodingType, marker string, maxKeys int32, keyPrefix string) *model.UsgResp
}
