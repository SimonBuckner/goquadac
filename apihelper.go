package goquadac

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ApiHelper struct {
	baseUrl string
	headers http.Header
}

func NewApiHelper(baseUrl string) *ApiHelper {

	return &ApiHelper{
		baseUrl: baseUrl,
		headers: make(http.Header),
	}
}

func (api *ApiHelper) SetDefaultHeader(key, value string) *ApiHelper {
	api.headers.Set(key, value)
	return api
}

func (api *ApiHelper) SetAuthHeader(key, value string) *ApiHelper {
	api.headers.Set(key, value)
	return api
}

func (api *ApiHelper) NewGetQuery(endpoint string) *ApiQuery {
	query := api.newApiQuery(http.MethodGet, endpoint)
	return query
}

func (api *ApiHelper) NewPostQuery(endpoint string) *ApiQuery {
	query := api.newApiQuery(http.MethodPost, endpoint)
	return query
}

type ApiQuery struct {
	request          *http.Request
	response         *http.Response
	urlQuery         url.Values
	dumpRequest      bool
	dumpResponse     bool
	dumpResponseBody bool
}

func (api *ApiHelper) newApiQuery(method, endpoint string) *ApiQuery {
	urlPath, err := url.JoinPath(api.baseUrl, endpoint)
	PanicOnError("Cannout build query URL", err)

	req, err := http.NewRequest(method, urlPath, nil)
	PanicOnError("Cannot build request", err)

	query := ApiQuery{
		request:  req,
		urlQuery: url.Values{},
	}
	query.request.Header = api.headers
	return &query
}

func (query *ApiQuery) AddUrlQuery(key, value string) *ApiQuery {
	query.urlQuery.Add(key, value)
	return query
}

func (query *ApiQuery) SetDumpRequest(state bool) *ApiQuery {
	query.dumpRequest = state
	return query
}

func (query *ApiQuery) SetDumpResponse(state bool) *ApiQuery {
	query.dumpResponse = state
	return query
}

func (query *ApiQuery) SetDumpResponseBody(state bool) *ApiQuery {
	query.dumpResponseBody = state
	return query
}

func (query *ApiQuery) Call() (*ApiQuery, error) {
	if query.request == nil {
		return nil, fmt.Errorf("please build the request before calling it")
	}
	if query.dumpRequest {
		query.DumpRequest()
	}
	res, err := http.DefaultClient.Do(query.request)
	query.response = res
	if query.dumpResponse {
		query.DumpRespone(query.dumpResponseBody)
	}

	return query, err
}

func (query *ApiQuery) DecodeJsonBody(out any) error {

	content, err := io.ReadAll(query.response.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, out)
}

func (query *ApiQuery) DumpRequest() {
	out, err := httputil.DumpRequestOut(query.request, false)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	fmt.Println()
}

func (query *ApiQuery) DumpRespone(dumpBody bool) {
	out, err := httputil.DumpResponse(query.response, dumpBody)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	fmt.Println()
}

func (query *ApiQuery) ResponsOK() bool {
	return query.response != nil && query.response.StatusCode == http.StatusOK
}

func (query *ApiQuery) Get(out any) error {
	_, err := query.Call()
	if !query.ResponsOK() {
		return err
	}
	return query.DecodeJsonBody(&out)
}
