package algo

import (
	"testing"
)

func TestAwsSign(t *testing.T) {
	// cfg, err := config.LoadDefaultConfig(context.TODO(),
	// 	config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("guichen01", "guichen01", "")),
	// 	config.WithRegion("us-east-1"),
	// )
	// if err != nil {
	// 	fmt.Println("Error loading AWS configuration:", err)
	// 	return
	// }

	// signer := v4.NewSigner()

	// req, _ := http.NewRequest("GET", "http://guichen01.bcc-bdbl.baidu.com:8021/", nil)

	// err = signer.SignHTTP(context.TODO(), cfg.Credentials, req, "", "s3", "us-east-1", time.Now())

	// actualAuth := req.Header.Get("Authorization")
	// expectAuth := "AWS4-HMAC-SHA256 Credential=guichen01/20230918/us-east-1/s3/aws4_request, SignedHeaders=accept-encoding;amz-sdk-invocation-id;amz-sdk-request;host;x-amz-content-sha256;x-amz-date, Signature=a6e9a2c161b4663a55cb4e1943725c9f1e08c12864aa38726676655a194ac1bf"
	// assert.Equal(t, expectAuth, actualAuth)
}
