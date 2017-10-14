package aero

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
)

// BodyReader represents a request body.
type BodyReader struct {
	reader io.ReadCloser
}

// JSON parses the body as a JSON object.
func (body BodyReader) JSON() (interface{}, error) {
	decoder := json.NewDecoder(body.reader)
	defer body.reader.Close()

	var data interface{}
	err := decoder.Decode(&data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

// JSONObject parses the body as a JSON object and returns a map[string]interface{}.
func (body BodyReader) JSONObject() (map[string]interface{}, error) {
	json, err := body.JSON()

	if err != nil {
		return nil, err
	}

	data, ok := json.(map[string]interface{})

	if !ok {
		return nil, errors.New("Invalid format: Expected JSON object")
	}

	return data, nil
}

// Bytes returns a slice of bytes containing the request body.
func (body BodyReader) Bytes() ([]byte, error) {
	data, err := ioutil.ReadAll(body.reader)
	defer body.reader.Close()

	if err != nil {
		return nil, err
	}

	return data, nil
}

// String returns a string containing the request body.
func (body BodyReader) String() (string, error) {
	bytes, err := body.Bytes()

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
