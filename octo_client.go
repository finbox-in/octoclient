package octoclient

import (
	"net/http"

	"github.com/google/uuid"
)

type OctoPayload struct {
	ClientID	uuid.UUID				`json:"client_id`
	ServiceID	uuid.UUID				`json:"service_id`
	Data		map[string]interface{}	`json:"data"`
}

type OctoResponse struct {
	Message		string					`json:"message"`
	RequestID	uuid.UUID				`json:"request_id"`
	Data 		map[string]interface{}	`json:"data"`
	MetaData	map[string]interface{}	`json:"meta_data"`
}

type OctoClient struct {
	HTTPClient *http.Client
	BaseURL			string
}

func New( baseUrl string) *OctoClient{
	//TODO: trim the base url right so that it can be used with other endpoints 
	return &OctoClient{
		HTTPClient: &http.Client{},
		BaseURL: baseUrl,
	}
}

func (o *OctoClient) getHttpClient() http.Client {
	return *o.HTTPClient
}

func (o *OctoClient) ServiceInvoke(payload OctoPayload) {

	// take the /service/invoke endpoint and make new URL

	// have the method defined 

	// convert to json 

	// verify the payload if needed 

	// call the url 
	
	// make the OctoResponse 

	// return the same 

}

// TODO: func convert to JSON

// TODO: func JSON to payload 