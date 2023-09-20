package api

import (
	"context"
	"strconv"

	"github.com/fakeyanss/usg-go/internal/backend"
	"github.com/fakeyanss/usg-go/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterS3Api(group *gin.RouterGroup) {
	s3ctl := &s3Ctl{
		backend: backend.NewS3Backend(),
		// backend: backend.NewFsBackend(),
	}

	group.Use(genReqID, errHandler)

	group.GET("/", s3ctl.listBuckets)
	group.GET("/:bucket", s3ctl.listObjects)
}

func genReqID(c *gin.Context) {
	var reqID string
	reqID = c.Request.Header.Get(model.UsgReqIdHeaderKey)
	if reqID == "" {
		reqID = uuid.NewString()
		c.Request.Header.Add(model.UsgReqIdHeaderKey, reqID)
	}
	c.Header(model.UsgReqIdHeaderKey, reqID)

	c.Next()
}

func errHandler(c *gin.Context) {
	c.Next()

	if n := len(c.Errors); n > 0 {
		e := c.Errors[n-1]
		if err := e.Err; err != nil {
			if usgErr, ok := err.(*model.UsgError); ok {
				writeResponseByContentType(c, usgErr.HttpCode, usgErr.ToUsgErrorBody(c))
			}
		}
		return
	}
}

func passResponse(c *gin.Context, ur *model.UsgResp) {
	if ur.Header != nil {
		for k, v := range *ur.Header {
			for _, ele := range v {
				c.Writer.Header().Add(k, ele)
			}
		}
	}
	// remove content-length to avoid writing body more than the declared
	c.Header("Content-Length", "")
	// remove content-type
	c.Header("Content-Type", "")

	writeResponseByContentType(c, ur.StatusCode, ur.Body)
}

func writeResponseByContentType(c *gin.Context, code int, body any) {
	if c.Request.Header.Get("Content-Type") == "application/json" {
		c.JSON(code, body)
	} else {
		c.XML(code, body)
	}
}

type s3Ctl struct {
	backend backend.StorageBackend
}

func (s3 *s3Ctl) listBuckets(c *gin.Context) {
	ctx := context.WithValue(context.Background(), model.UsgReqURICtx{}, c.Request.URL.Path)
	ctx = context.WithValue(ctx, model.UsgReqHeaderCtx{}, c.Request.Header)
	ur := s3.backend.ListBuckets(ctx)
	passResponse(c, ur)
}

func (s3 *s3Ctl) listObjects(c *gin.Context) {
	ctx := context.WithValue(context.Background(), model.UsgReqURICtx{}, c.Request.URL.Path)
	ctx = context.WithValue(ctx, model.UsgReqHeaderCtx{}, c.Request.Header)

	bucket := c.Param("bucket")
	delimiter := c.Query("delimiter")
	encodingType := c.Query("encoding-type")
	marker := c.Query("marker")
	maxKeysStr := c.Query("max-keys")
	var maxKeys int32 = 1000
	if maxKeysStr != "" {
		maxKeys64, err := strconv.ParseInt(maxKeysStr, 10, 32)
		if err != nil {
			_ = c.Error(model.NewInvalidArgument("Invalid argument with max-keys"))
			return
		}
		maxKeys = int32(maxKeys64)
	}
	prefix := c.Query("prefix")
	ur := s3.backend.ListObjects(ctx, bucket, delimiter, encodingType, marker, maxKeys, prefix)
	passResponse(c, ur)
}
