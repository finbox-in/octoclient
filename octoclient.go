package octoclient

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
)

/*
Usage:
  - Create instance of Octo-Client once
  - call the service-invoke using the clientID, payload.
  - The other features like pathParams will be included in payload
*/
type OctoPayload struct {
	ServiceID uuid.UUID              `json:"serviceID`
	Data      map[string]interface{} `json:"data"`
}

type OctoResponse struct {
	Message   string                 `json:"msg"`
	RequestID uuid.UUID              `json:"requestId"`
	Data      map[string]interface{} `json:"data"`
}

type OctoClient struct {
	HTTPClient *http.Client
	BaseURL    string
}

func New(baseUrl string) *OctoClient {
	baseUrl = trimTrailingSlash(baseUrl)
	return &OctoClient{
		HTTPClient: &http.Client{},
		BaseURL:    baseUrl,
	}
}

func (o *OctoClient) getHttpClient() http.Client {
	return *o.HTTPClient
}

func (o *OctoClient) ServiceInvoke(clientID string, payload OctoPayload) (OctoResponse, error) {

	// TODO: clientID will be replaced with token in coming future
	apiEndpoint := "/service/invoke"
	callingUrl := o.BaseURL + apiEndpoint
	method := "POST"
	contentType := "application/json"
	var response OctoResponse

	flag := IsValidID(payload.ServiceID.String())&& IsValidID(clientID)
	if !flag {
		return response, errors.New("invalid id entered")
	}

	finalPayload, err := ConvertStructToJSON(payload)
	if err != nil {
		return response, err
	}

	req, err := http.NewRequest(method, callingUrl, finalPayload)

	if err != nil {
		return response, err
	}
	req.Header.Add("clientId", clientID)
	req.Header.Add("Content-Type", contentType)

	res, err := o.HTTPClient.Do(req)
	if err != nil {
		return response, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, err
	}
	// TODO: Handling if return type !json
	response, err = ConvertByteToStruct(body)
	if err != nil {
		return response, err
	}

	return response, nil
}
