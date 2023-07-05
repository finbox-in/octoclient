package octoclient

import (
	"errors"
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
	baseURL    string
	clientID   uuid.UUID
}

type OctoConfig struct {
	// Other options in http.Client will be added here e.g, custom timeout
	BaseURL  string
	ClientID string // Use AccessToken in place of clientID
}

func New(options OctoConfig) *OctoClient {
	baseURL := trimTrailingSlash(options.BaseURL)
	clientID, _ := uuid.Parse(options.ClientID)
	return &OctoClient{
		HTTPClient: &http.Client{},
		baseURL:    baseURL,
		clientID:   clientID,
	}
}

func (o *OctoClient) getHttpClient() http.Client {
	return *o.HTTPClient
}

func (o *OctoClient) ServiceInvoke(payload OctoPayload) (OctoResponse, error) {

	callingUrl := o.baseURL + apiEndpoint
	var response OctoResponse

	flag := IsValidID(payload.ServiceID.String()) && IsValidID(o.clientID.String())
	if !flag {
		return OctoResponse{}, errors.New("invalid id entered")
	}

	finalPayload, err := ConvertStructToJSON(payload)
	if err != nil {
		return OctoResponse{}, err
	}

	req, err := http.NewRequest(method, callingUrl, finalPayload)

	if err != nil {
		return response, err
	}
	req.Header.Add("clientId", o.clientID.String())
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
		return OctoResponse{}, err
	}

	return response, nil
}
