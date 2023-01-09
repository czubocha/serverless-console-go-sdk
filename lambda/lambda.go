package lambda

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

const (
	sdkName                         = "aws-lambda-sdk-go"
	sdkVersion                      = "0.0.1"
	externalExtensionHttpServerHost = "http://localhost:2773"
	protobufContentType             = "application/x-protobuf"
)

// Start wraps handler function with serverless instrumentation wrappers
// and starts lambda using standard lambda.Start from aws-lambda-go AWS library.
func Start(handler any) {
	h := lambda.NewHandler(handler)
	w := newWrapper()
	h = w.requestResponse(h.Invoke)
	lambda.Start(h)
}

func requestID(ctx context.Context) string {
	if lambdaContext, ok := lambdacontext.FromContext(ctx); ok {
		return lambdaContext.AwsRequestID
	}
	return ""
}
