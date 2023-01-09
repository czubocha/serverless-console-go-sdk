package lambda

import (
	"bytes"
	"context"
	tagsv1 "go.buf.build/protocolbuffers/go/serverless/sdk-schema/serverless/instrumentation/tags/v1"
	schema "go.buf.build/protocolbuffers/go/serverless/sdk-schema/serverless/instrumentation/v1"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
	"time"
)

const disableReqResEnvVar = "SLS_DISABLE_REQUEST_RESPONSE_MONITORING"

func (w wrapper) requestResponse(handler bytesHandlerFunc) bytesHandlerFunc {
	if isTrue(disableReqResEnvVar) {
		return handler
	}
	return func(ctx context.Context, payload []byte) ([]byte, error) {
		reqID := requestID(ctx)
		w.sendReq(ctx, payload, reqID)
		output, err := handler(ctx, payload)
		w.sendRes(ctx, output, reqID)
		return output, err
	}
}

func (w wrapper) sendReq(ctx context.Context, payload []byte, reqID string) {
	w.sendReqRes(ctx, payload, reqID, schema.RequestResponse_ORIGIN_REQUEST)
}

func (w wrapper) sendRes(ctx context.Context, payload []byte, reqID string) {
	w.sendReqRes(ctx, payload, reqID, schema.RequestResponse_ORIGIN_RESPONSE)
}

func (w wrapper) sendReqRes(ctx context.Context, payload []byte, reqID string, origin schema.RequestResponse_Origin) {
	payloadString := string(payload)
	msg := w.reqResMsg(&reqID, &payloadString, origin)
	url := w.externalExtensionURL + "/request-response"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(msg))
	if err != nil {
		log.Print(err)
	}
	req.Header.Set("Content-Type", protobufContentType)

	response, err := w.httpClient.Do(req)
	if err != nil {
		log.Print(err)
	}
	if response.StatusCode != http.StatusOK {
		log.Printf("%s status code %d", url, response.StatusCode)
	}
}

func (w wrapper) reqResMsg(reqID, payloadString *string, origin schema.RequestResponse_Origin) []byte {
	timestamp := uint64(time.Now().UnixNano())
	msg := &schema.RequestResponse{
		SlsTags: &tagsv1.SlsTags{
			OrgId:   w.organizationID,
			Service: w.lambdaName,
			Sdk: &tagsv1.SdkTags{
				Name:    sdkName,
				Version: sdkVersion,
			},
		},
		RequestId: reqID,
		Body:      payloadString,
		Origin:    origin,
		Timestamp: &timestamp,
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		log.Print(err)
	}
	return msgBytes
}
