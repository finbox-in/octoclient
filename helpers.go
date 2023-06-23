package octoclient

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"
)

func ConvertStructToJSON(obj OctoPayload) (*strings.Reader, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return strings.NewReader(string(jsonData)), nil
}


func ConvertByteToStruct(body []byte) (OctoResponse, error) {
	var response OctoResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		return OctoResponse{}, err
	}
	return response, nil
}

func trimTrailingSlash(url string) string {
	forwardSlash :="/"
	return strings.TrimSuffix(url, forwardSlash)
}

func IsValidID(clientID string) bool {
	_, err := uuid.Parse(clientID)
	return err == nil
}