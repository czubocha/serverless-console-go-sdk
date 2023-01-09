package lambda

import (
	"context"
	"net/http"
	"os"
)

type (
	wrapper struct {
		organizationID       string
		lambdaName           string
		externalExtensionURL string
		httpClient           doer
	}
	doer interface {
		Do(req *http.Request) (*http.Response, error)
	}
)

func newWrapper() *wrapper {
	return &wrapper{
		organizationID:       os.Getenv("SLS_DEV_MODE_ORG_ID"),
		lambdaName:           os.Getenv("AWS_LAMBDA_FUNCTION_NAME"),
		httpClient:           &http.Client{},
		externalExtensionURL: externalExtensionHttpServerHost,
	}
}

type bytesHandlerFunc func(context.Context, []byte) ([]byte, error)

func (f bytesHandlerFunc) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	return f(ctx, payload)
}

func isTrue(envVarName string) bool {
	v := os.Getenv(envVarName)
	return v == "TRUE" || v == "1"
}
