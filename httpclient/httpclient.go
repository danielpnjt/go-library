package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/danielpnjt/go-library/contextwrap"
	"github.com/danielpnjt/go-library/log"

	"go.elastic.co/apm/module/apmhttp"
)

type contextKey struct {
	name string
}

var (
	client   *http.Client
	traceKey = contextKey{"trace"}
)

type TraceHttp struct {
	Request  interface{} `json:"request"`
	Response interface{} `json:"response"`
	Url      string      `json:"url"`
	Elapsed  string      `json:"elapsed"`
}

func Init() {
	client = &http.Client{
		Timeout: 20 * time.Second,
	}

	client = apmhttp.WrapClient(client)
}

func InitWithParam(c *http.Client) {
	client = apmhttp.WrapClient(c)
}

func Call(ctx context.Context, requestBody map[string]interface{}, header http.Header, endpoint string) (context.Context, []byte, http.Header, error) {
	start := time.Now()
	jsonRequest, _ := json.Marshal(requestBody)

	var payload *bytes.Reader

	if _, ok := header[http.CanonicalHeaderKey("X-CLIENT-ID")]; ok {
		var param = url.Values{}
		param.Set("request", string(jsonRequest))
		payload = bytes.NewReader([]byte(param.Encode()))
	} else {
		payload = bytes.NewReader(jsonRequest)
	}

	currentTrace := contextwrap.GetTraceFromContext(ctx)

	request, err := http.NewRequest("POST", endpoint, payload)
	if err != nil {
		return ctx, nil, nil, err
	}

	request.Header = header

	response, err := client.Do(request.WithContext(ctx))
	if err != nil {
		return ctx, nil, nil, err
	}

	defer response.Body.Close()

	responseByte, err := io.ReadAll(response.Body)
	if err != nil {
		return ctx, nil, nil, err
	}

	elapsed := time.Since(start).String()

	tr := &TraceHttp{
		Url:     endpoint,
		Request: log.Minify(requestBody),
		Elapsed: elapsed,
	}

	currentTrace = append(currentTrace, tr)

	ctx = context.WithValue(ctx, traceKey, currentTrace)

	var js map[string]interface{}
	err = json.Unmarshal(responseByte, &js)
	if err != nil {
		return ctx, nil, nil, err
	}

	tr.Response = log.Minify(js)

	return ctx, responseByte, response.Header, nil
}
