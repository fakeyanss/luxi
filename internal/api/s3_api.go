package api

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var (
	s3Cli *s3.Client

	UsgReqIdHeader = "usg-Req-Id"
	UsgReqKey      = "UsgRequest"
)

type (
	AmzResHeaderCtx struct{}

	UsgReqHeaderCtx struct{}
)

func initS3() {
	sdkConf, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("guichen01", "guichen011", "")),
		config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to init s3 client")
	}
	s3Cli = s3.NewFromConfig(sdkConf, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String("http://guichen01.bcc-bdbl.baidu.com:8021")
	}, s3.WithAPIOptions(passResquestHeaderFunc, storeResponseHeaderFunc))
}

func passResquestHeaderFunc(stack *middleware.Stack) error {
	return stack.Build.Add(middleware.BuildMiddlewareFunc("PassRequestHeader",
		func(ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler) (middleware.BuildOutput, middleware.Metadata, error) {
			switch req := in.Request.(type) {
			case *smithyhttp.Request:
				header := ctx.Value(UsgReqHeaderCtx{}).(http.Header)
				for k, v := range header {
					for _, ele := range v {
						req.Header.Add(k, ele)
					}
				}
			}
			return next.HandleBuild(ctx, in)
		}), middleware.Before)
}

func storeResponseHeaderFunc(stack *middleware.Stack) error {
	return stack.Deserialize.Add(middleware.DeserializeMiddlewareFunc("StoreResponseHeader",
		func(ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler) (middleware.DeserializeOutput, middleware.Metadata, error) {
			out, meta, err := next.HandleDeserialize(ctx, in)
			if err != nil {
				return out, meta, err
			}
			switch res := out.RawResponse.(type) {
			case *smithyhttp.Response:
				meta.Set(AmzResHeaderCtx{}, res.Header)
			}

			return out, meta, nil
		}), middleware.Before)
}

func genReqID(c *gin.Context) {
	reqID := uuid.NewString()
	c.Header(UsgReqIdHeader, reqID)
}

func passResponseHeader(c *gin.Context, meta middleware.Metadata) {
	header := meta.Get(AmzResHeaderCtx{}).(http.Header)
	for k, v := range header {
		for _, ele := range v {
			c.Writer.Header().Add(k, ele)
		}
	}
	// remove content-length to avoid writing body more than the declared
	c.Header("content-length", "")
}

func RegisterS3Api(group *gin.RouterGroup) {
	initS3()

	group.Use(genReqID)

	group.GET("/", listBuckets)
	group.GET("/:bucket", listObjects)
}
