package main

import (
	"context"
	"fmt"

	"github.com/finbox-in/octoclient"
)

var options = octoclient.Options{
	BaseURL:       "", // Octopus URL ( UAT or Prod provided )
	Token:         "", // Token or ClientID provided
	Authorization: "", //Auth token provided by Octopus
}

func main() {
	serviceID := "" // serviceID provided
	requestID := "" // requestID gotten from an external service (if any)

	//OctoClient: Create a sample payload
	var payload = octoclient.OctoPayload{
		ServiceID: serviceID,
		Data: map[string]interface{}{
			"key": "value",
		},
		RequestID: requestID,
	}

	var octoClient = octoclient.New(options)

	response, err := octoClient.ServiceInvoke(context.TODO(), payload)
	if err != nil {
		fmt.Println("err", err)
	}

	fmt.Println(response.Message, response.Data, response.RequestID)
}
