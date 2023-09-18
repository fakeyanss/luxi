package api

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func listObjects(c *gin.Context) {
	bucket := c.Param("bucket")
	ctx := context.WithValue(context.Background(), UsgReqHeaderCtx{}, c.Request.Header)
	res, err := s3Cli.ListObjects(ctx, &s3.ListObjectsInput{Bucket: aws.String(bucket)})
	if err != nil {
		log.Error().Err(err).Msg("Fail to list objects")
		c.XML(http.StatusBadRequest, nil)
	}
	if res != nil {
		passResponseHeader(c, res.ResultMetadata)
		c.XML(http.StatusOK, &res)
	}
}
