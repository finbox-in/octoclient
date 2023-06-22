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
  - call the service-invoke using the clientID, serviceID, payload & other params and feature
*/
type OctoPayload struct {
	ClientID  uuid.UUID              `json:"client_id`
	ServiceID uuid.UUID              `json:"service_id`
	Data      map[string]interface{} `json:"data"`
}

type OctoResponse struct {
	Message   string                 `json:"message"`
	RequestID uuid.UUID              `json:"request_id"`
	Data      map[string]interface{} `json:"data"`
	MetaData  map[string]interface{} `json:"meta_data"`
}

type OctoClient struct {
	HTTPClient *http.Client
	BaseURL    string
}

func New(baseUrl string) *OctoClient {
	//TODO: trim the base url right so that it can be used with other endpoints
	baseUrl = trimTrailingSlash(baseUrl)
	return &OctoClient{
		HTTPClient: &http.Client{},
		BaseURL:    baseUrl,
	}
}

func (o *OctoClient) getHttpClient() http.Client {
	return *o.HTTPClient
}

func (o *OctoClient) ServiceInvoke(serviceID, clientID string, payload OctoPayload) (OctoResponse, error) {

	// TODO: clientID will be replaced with token in coming future
	apiEndpoint := "/service/invoke"
	callingUrl := o.BaseURL + apiEndpoint
	method := "POST"
	var response OctoResponse

	flag := IsValidID(serviceID) && IsValidID(clientID)
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
	req.Header.Add("Content-Type", "application/json")

	res, err := o.HTTPClient.Do(req)
	if err != nil {
		return response, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, err
	}
	// Handle if return type is HTML
	response, err = ConvertByteToStruct(body)
	if err != nil {
		return response, err
	}

	return response, nil
}
