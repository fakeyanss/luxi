package api

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func listBuckets(c *gin.Context) {
	ctx := context.WithValue(context.Background(), UsgReqHeaderCtx{}, c.Request.Header)
	res, err := s3Cli.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		log.Error().Err(err).Msg("Fail to list buckets")
		c.XML(http.StatusBadRequest, nil)
	}
	if res != nil {
		passResponseHeader(c, res.ResultMetadata)
		c.XML(http.StatusOK, res)
	}
}
