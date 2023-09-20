package model

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

type (
	UsgReqHeaderCtx struct{}
	UsgReqURICtx    struct{}
	UsgRespCtx      struct{}
)

var (
	UsgReqIdHeaderKey = "usg-Req-Id"
)

type UsgResp struct {
	StatusCode int
	Header     *http.Header
	Body       any
}

type UsgErrorBody struct {
	XMLName   xml.Name `xml:"Error" json:"-"` // 指定根节点名称
	Code      string
	Message   string
	Resource  string
	RequestId string
}

func (e *UsgErrorBody) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
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
