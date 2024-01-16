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
	for _, fileField := range fileFields {
		file, err := os.Open(fileField.FilePath)
		if err != nil {
			return err
		}
		defer file.Close()
		fileName := filepath.Base(fileField.FilePath)

		fileWriter, err := multiPartWriter.CreateFormFile(fileField.FieldName, fileName)
		if err != nil {

			return err
		}

		_, err = io.Copy(fileWriter, file)
		if err != nil {
			return err
		}
	}
	return nil
}
