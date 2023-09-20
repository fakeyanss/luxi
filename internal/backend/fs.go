package backend

import (
	"context"
	"io/fs"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/fakeyanss/usg-go/internal/model"
	"github.com/rs/zerolog/log"
)

type FsBackend struct {
	rootPath string
}

func NewFsBackend() *FsBackend {
	fsb := &FsBackend{
		rootPath: "data/fs",
	}
	initRootPath(fsb.rootPath)
	return fsb
}

func initRootPath(rootPath string) {
	err := os.MkdirAll(rootPath, fs.ModePerm)
	log.Fatal().Err(err).Msg("Fail to init fs backend, exit.")
}

func (fsb *FsBackend) CreateBucket(bucket string) *model.UsgResp {
	return nil
}

func (fsb *FsBackend) ListBuckets(ctx context.Context) *model.UsgResp {
	files, err := os.ReadDir(fsb.rootPath)
	if err != nil {
		log.Error().Err(err).Msgf("Fail to read directory=%s", fsb.rootPath)
	}
	output := &s3.ListBucketsOutput{
		Buckets: []types.Bucket{},
	}

	var ownerName string
	var ownerID string
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		info, err := os.Stat(fsb.rootPath + "/" + f.Name())
		var creationDate time.Time
		if err != nil {
			// todo
		}
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			if ownerName == "" {
				ownerName = strconv.Itoa(int(stat.Uid))
				ownerID = ownerName
			}
			cTimespec := stat.Birthtimespec
			// mills := cTimespec.Nano() / 1000000
			creationDate = time.Unix(cTimespec.Sec, cTimespec.Nano())
		}

		output.Buckets = append(output.Buckets, types.Bucket{
			Name:         aws.String(f.Name()),
			CreationDate: &creationDate,
		})
	}
	output.Owner = &types.Owner{
		DisplayName: &ownerName,
		ID:          &ownerID,
	}
	ur := &model.UsgResp{}
	ur.StatusCode = 200
	// ur.Header = &header
	ur.Body = output
	return ur
}

func (fsb *FsBackend) ListObjects(ctx context.Context, bucket, delimiter, encodingType, marker string, maxKeys int32, keyPrefix string) *model.UsgResp {
	path := fsb.rootPath + "/" + bucket
	files, err := os.ReadDir(fsb.rootPath)
	if err != nil {

	}
	output := &s3.ListObjectsOutput{
		Contents: []types.Object{},
		MaxKeys:  1000,
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		info, err := f.Info()
		var creationDate time.Time
		if err != nil {
			creationDate = info.ModTime()
		}

		output.Contents = append(output.Contents, types.Object{
			ETag:         aws.String(""),
			Key:          aws.String(path + "/" + f.Name()),
			LastModified: &creationDate,
			Owner: &types.Owner{
				DisplayName: aws.String(info.Name()),
				ID:          aws.String(info.Name()),
			},
			Size:         info.Size(),
			StorageClass: types.ObjectStorageClassStandard,
		})
	}
	ur := &model.UsgResp{}
	ur.StatusCode = 200
	// ur.Header = &header
	ur.Body = output
	return ur
}
