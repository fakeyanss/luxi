package model

import (
	"io"
	"net/http"
	"time"
)

type UsgResponse struct {
	StatusCode int
	Header     http.Header
	Body       io.ReadCloser
}

type UsgBucket struct {
	Name         string
	CreationDate time.Time
}

type UsgBucketNameSorter []*UsgBucket

func (u UsgBucketNameSorter) Len() int {
	return len(u)
}

func (u UsgBucketNameSorter) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u UsgBucketNameSorter) Less(i, j int) bool {
	return u[i].Name < u[j].Name
}

type UsgObject struct {
	ETag         string
	Key          string
	LastModified time.Time
	Size         uint64
	StorageClass string
	UsgOwner
}

type UsgObjectNameSorter []*UsgObject

func (u UsgObjectNameSorter) Len() int {
	return len(u)
}

func (u UsgObjectNameSorter) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u UsgObjectNameSorter) Less(i, j int) bool {
	return u[i].Key < u[j].Key
}

type UsgOwner struct {
	DisplayName string
	ID          string
}
