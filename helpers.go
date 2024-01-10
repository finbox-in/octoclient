package octoclient

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
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
	forwardSlash := "/"
	return strings.TrimSuffix(url, forwardSlash)
}

func processTextFields(textFields []OctoTextField, multiPartWriter *multipart.Writer) error {
	var err error
	for _, textField := range textFields {
		err = multiPartWriter.WriteField(textField.FieldName, textField.FieldValue)
		if err != nil {
			return err
		}
	}
	return err
}

func processFileFields(fileFields []OctoFileField, multiPartWriter *multipart.Writer) error {
	totalFileFields := len(fileFields)
	errChannels := make([]chan error, totalFileFields)
	for i := range errChannels {
		errChannels[i] = make(chan error)
	}

	for iter, fileField := range fileFields {
		errChan := errChannels[iter]
		go func(fieldName string, filePath string, errChan chan error) {
			file, err := os.Open(filePath)
			if err != nil {
				errChan <- err
				return
			}
			defer file.Close()
			fileName := filepath.Base(filePath)

			fileWriter, err := multiPartWriter.CreateFormFile(fieldName, fileName)
			if err != nil {
				errChan <- err
				return
			}

			_, err = io.Copy(fileWriter, file)
			if err != nil {
				errChan <- err
				return
			}
			errChan <- nil
		}(fileField.FieldName, fileField.FilePath, errChan)
	}
	for i := 0; i < totalFileFields; i++ {
		err := <-errChannels[i]
		if err != nil {
			return err
		}
	}
	return nil
}
