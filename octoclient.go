package octoclient

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
)

const (
	apiEndpoint = "/service/invoke"
	method      = "POST"
	contentType = "application/json"
)

/*
Usage:
  - Create object of OctoConfig with clientID/access-token, baseURL of Octopus & others
  - Create instance of Octo-Client once using this object
  - call the service-invoke using the payload.
  - The other features like pathParams will be included in payload
*/
type OctoPayload struct {
	ServiceID   string                 `json:"serviceID`
	QueryParams map[string]interface{} `json:"queryParameters"`
	Data        map[string]interface{} `json:"data"`
}

type OctoResponse struct {
	Message   string                 `json:"msg"`
	RequestID uuid.UUID              `json:"requestId"`
	Data      map[string]interface{} `json:"data"`
}

type OctoClient struct {
	HTTPClient *http.Client
	baseURL    string
	token      string
}

type Options struct {
	// Other options in http.Client will be added here e.g, custom timeout
	BaseURL string
	Token   string // Use AccessToken in place of clientID
}

func New(options Options) *OctoClient {
	baseURL := trimTrailingSlash(options.BaseURL)
	return &OctoClient{
		HTTPClient: &http.Client{},
		baseURL:    baseURL,
		token:      options.Token,
	}
}

func (o *OctoClient) getHttpClient() http.Client {
	return *o.HTTPClient
}

func (o *OctoClient) ServiceInvoke(ctx context.Context, payload OctoPayload) (*OctoResponse, error) {

	callingUrl := o.baseURL + apiEndpoint
	var response OctoResponse

	finalPayload, err := ConvertStructToJSON(payload)
	if err != nil {
		return nil, err
	}
	var req *http.Request
	if ctx == nil {
		req, err = http.NewRequestWithContext(context.TODO(), method, callingUrl, finalPayload)
	} else {
		req, err = http.NewRequestWithContext(ctx, method, callingUrl, finalPayload)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Add("clientId", o.token)
	req.Header.Add("Content-Type", contentType)

	res, err := o.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// TODO: Handling if return type !json
	response, err = ConvertByteToStruct(body)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
