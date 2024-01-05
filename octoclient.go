package octoclient

import (
	"bytes"
	"context"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
)

const (
	apiEndpoint     = "/service/invoke"
	apiEndpointFile = "/service/invoke-file"
	method          = "POST"
	contentType     = "application/json"
)

/*
Usage:
  - Create object of OctoConfig with clientID/access-token, baseURL of Octopus & others
  - Create instance of Octo-Client once using this object
  - call the service-invoke using the payload.
  - The other features like pathParams will be included in payload
*/
type OctoQueryParams struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type OctoHeaders struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type OctoPayload struct {
	ServiceID      string                 `json:"serviceID"`
	QueryParams    []OctoQueryParams      `json:"queryParameters"`
	DynamicHeaders []OctoHeaders          `json:"dynamicHeaders"`
	Data           map[string]interface{} `json:"data"`
	RequestID      string                 `json:"requestID"` // Acts as unique identifier for each request.
}

type OctoFileField struct {
	FieldName string
	FilePath  string
}
type OctoTextField struct {
	FieldName  string
	FieldValue string
}
type OctoPayloadForm struct {
	ServiceID  string          `json:"serviceID"`
	TextFields []OctoTextField `json:"textFields"`
	FileFields []OctoFileField `json:"fileFields"`
}

type OctoResponse struct {
	Message   string                 `json:"msg"`
	RequestID uuid.UUID              `json:"requestId"`
	Data      map[string]interface{} `json:"data"`
}

type OctoClient struct {
	HTTPClient    *http.Client
	baseURL       string
	token         string
	authorization string
}

type Options struct {
	// Other options in http.Client will be added here e.g, custom timeout
	BaseURL       string
	Token         string // Use AccessToken in place of clientID
	Authorization string // Auth Token
}

func New(options Options) *OctoClient {
	baseURL := trimTrailingSlash(options.BaseURL)
	return &OctoClient{
		HTTPClient:    &http.Client{},
		baseURL:       baseURL,
		token:         options.Token,
		authorization: options.Authorization,
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
	req.Header.Add("Authorization", o.authorization)

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

func (o *OctoClient) ServiceInvokeForm(ctx context.Context, payload OctoPayloadForm) (*OctoResponse, error) {
	callingUrl := o.baseURL + apiEndpointFile

	var requestBody bytes.Buffer

	multiPartWriter := multipart.NewWriter(&requestBody)
	err := multiPartWriter.WriteField("serviceID", payload.ServiceID)
	if err != nil {
		return nil, err
	}

	err = processTextFields(payload.TextFields, multiPartWriter)
	if err != nil {
		return nil, err
	}

	err = processFileFields(payload.FileFields, multiPartWriter)
	if err != nil {
		return nil, err
	}

	err = multiPartWriter.Close()
	if err != nil {
		return nil, err
	}

	var response OctoResponse
	var req *http.Request

	if ctx == nil {
		req, err = http.NewRequestWithContext(context.TODO(), method, callingUrl, &requestBody)
	} else {
		req, err = http.NewRequestWithContext(ctx, method, callingUrl, &requestBody)
	}

	if err != nil {
		return nil, err
	}
	req.Header.Add("clientId", o.token)
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+multiPartWriter.Boundary())
	req.Header.Add("Authorization", o.authorization)

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
