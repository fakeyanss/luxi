package backend

import (
	"context"
	"errors"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/fakeyanss/usg-go/internal/model"
	"github.com/rs/zerolog/log"
)

var (
	s3Cli *s3.Client

	AmzResponseKey = "AmzResponse"
	UsgReqIdHeader = "usg-Req-Id"
	UsgReqKey      = "UsgRequest"
)

type (
	AmzResHeaderCtx     struct{}
	AmzResStatusCodeCtx struct{}
)

type S3Backend struct {
	client *s3.Client
}

func NewS3Backend() *S3Backend {
	sdkConf, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("access_key", "secret_key", "")),
		config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to init s3 client")
	}
	s3Cli = s3.NewFromConfig(sdkConf, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String("http://localhost:9000")
	}, s3.WithAPIOptions(passResquestHeaderFunc, storeResponseFunc))
	return &S3Backend{
		client: s3Cli,
	}
}

func passResquestHeaderFunc(stack *middleware.Stack) error {
	return stack.Build.Add(middleware.BuildMiddlewareFunc("PassRequestHeader",
		func(ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler) (middleware.BuildOutput, middleware.Metadata, error) {
			switch req := in.Request.(type) {
			case *smithyhttp.Request:
				header := ctx.Value(model.UsgReqHeaderCtx{}).(http.Header)
				for k, v := range header {
					for _, ele := range v {
						req.Header.Add(k, ele)
					}
				}
			}
			return next.HandleBuild(ctx, in)
		}), middleware.Before)
}

func storeResponseFunc(stack *middleware.Stack) error {
	return stack.Deserialize.Add(middleware.DeserializeMiddlewareFunc("StoreAmzResponse",
		func(ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler) (middleware.DeserializeOutput, middleware.Metadata, error) {
			out, meta, err := next.HandleDeserialize(ctx, in)
			if resp, ok := out.RawResponse.(*smithyhttp.Response); ok {
				meta.Set(AmzResHeaderCtx{}, resp.Header)
				meta.Set(AmzResStatusCodeCtx{}, resp.StatusCode)
			}
			if err != nil {
				log.Error().Err(err).Msgf("Fail to %s", stack.ID())
			}
			return out, meta, err
		}), middleware.Before)
}

func retrieveError(ctx context.Context, err error) *model.UsgResp {
	ur := &model.UsgResp{}
	body := &model.UsgErrorBody{}
	var e1 *awshttp.ResponseError
	if errors.As(err, &e1) {
		ur.StatusCode = e1.HTTPStatusCode()
		ur.Header = &e1.Response.Header
	}
	var e2 smithy.APIError
	if errors.As(err, &e2) {
		body.Code = e2.ErrorCode()
		body.Message = e2.ErrorMessage()
	}
	var e3 s3.ResponseError
	if errors.As(err, &e3) {
		body.RequestId = e3.ServiceRequestID()
	}
	body.Resource = ctx.Value(model.UsgReqURICtx{}).(string)
	ur.Body = body
	return ur
}

func retrieveResp(meta *middleware.Metadata, body any) *model.UsgResp {
	ur := &model.UsgResp{}
	ur.StatusCode = meta.Get(AmzResStatusCodeCtx{}).(int)
	header := meta.Get(AmzResHeaderCtx{}).(http.Header)
	ur.Header = &header
	ur.Body = &body
	return ur
}

func (s3b *S3Backend) CreateBucket(bucket string) *model.UsgResp {
	return nil
}

func (s3b *S3Backend) ListBuckets(ctx context.Context) *model.UsgResp {
	resp, err := s3b.client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return retrieveError(ctx, err)
	}
	return retrieveResp(&resp.ResultMetadata, &resp)
}

func (s3b *S3Backend) ListObjects(ctx context.Context, bucket, delimiter, encodingType, marker string, maxKeys int32, keyPrefix string) *model.UsgResp {
	resp, err := s3Cli.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket:       &bucket,
		Delimiter:    &delimiter,
		EncodingType: types.EncodingType(encodingType),
		Marker:       &marker,
		MaxKeys:      maxKeys,
		Prefix:       &keyPrefix,
	})
	if err != nil {
		return retrieveError(ctx, err)
	}
	return retrieveResp(&resp.ResultMetadata, &resp)
}
