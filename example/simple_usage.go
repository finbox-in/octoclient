package main

import (
	"fmt"

	"github.com/finbox-in/octoclient"
)

var octoConfig = octoclient.Options{
	BaseURL: "", // Octopus URL ( UAT or Prod provided )
	Token:   "", // Token or ClientID provided
}

func main() {
	serviceID := "" // serviceID provided

	//OctoClient: Create a sample payload
	var payload = octoclient.OctoPayload{
		ServiceID: serviceID,
		Data: map[string]interface{}{
			"key": "value",
		},
	}

	var octoClient = octoclient.New(octoConfig)

	response, err := octoClient.ServiceInvoke(nil, payload)
	if err != nil {
		fmt.Println("err", err)
	}

	fmt.Println(response.Message, response.Data, response.RequestID)
}
