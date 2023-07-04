# OctoClient

The OctoClient is a simple, lightweight Go library that is used to interact with external services via HTTP requests. The client encapsulates the complexity of HTTP requests and provides a neat, clean, and easy-to-use interface to users.

## Installation

Use the go get command to install the OctoClient.

```bash
go get github.com/finbox-in/octoclient
```

## Usage
Below is a brief explanation of the functions provided by the OctoClient.

Creating a new client
To create a new OctoClient, use the New function. This function requires the baseURL of the service you will interact with.

```bash
baseUrl := "http://localhost:8000"
client := octoclient.New(baseUrl)
```

### Invoking a Service
The ServiceInvoke function is used to invoke a service. It requires a clientID and payload which includes serviceID and the data to be sent. It returns an OctoResponse which includes a message, requestID and the data returned by the service.

```go
payload := octoclient.OctoPayload{
    ServiceID: uuid.New(), // Replace with your service ID
    Data: map[string]interface{}{
        "key1": "value1",
        "key2": "value2",
    },
}

clientID := "your-client-id"
response, err := client.ServiceInvoke(clientID, payload)

if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Printf("Message: %s\nRequestID: %s\nData: %v\n", 
    response.Message, response.RequestID, response.Data)
```

License
MIT
