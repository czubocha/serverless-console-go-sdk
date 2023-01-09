package lambda

import (
	"context"
	schema "go.buf.build/protocolbuffers/go/serverless/sdk-schema/serverless/instrumentation/v1"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_wrapper_requestResponse(t *testing.T) {
	requests := 0
	responses := 0
	server := testServer(t, &requests, &responses)
	defer server.Close()
	w := wrapper{
		organizationID:       "example-id",
		lambdaName:           "example-lambda-name",
		httpClient:           server.Client(),
		externalExtensionURL: server.URL,
	}
	handler := bytesHandlerFunc(func(ctx context.Context, payload []byte) ([]byte, error) {
		return []byte(`example output`), nil
	})
	handler = w.requestResponse(handler)
	_, _ = handler(context.Background(), []byte(`example payload`))
	if requests != 1 {
		t.Errorf("expected 1 request, got %d", requests)
	}
	if responses != 1 {
		t.Errorf("expected 1 response, got %d", responses)
	}
}

func testServer(t *testing.T, requests, responses *int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/request-response" {
			t.Errorf("expected to request '/request-response', got: %s", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/x-protobuf" {
			t.Errorf("expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error("unable to read request body")
		}
		reqRes := schema.RequestResponse{}
		if err := proto.Unmarshal(bytes, &reqRes); err != nil {
			t.Error("unable to unmarshal request body")
		}
		switch reqRes.Origin {
		case schema.RequestResponse_ORIGIN_REQUEST:
			*requests++
			if *reqRes.Body != `example payload` {
				t.Errorf("expected example payload, got: %s", *reqRes.Body)
			}
		case schema.RequestResponse_ORIGIN_RESPONSE:
			*responses++
			if *reqRes.Body != `example output` {
				t.Errorf("expected example output, got: %s", *reqRes.Body)
			}
		}
	}))
}
