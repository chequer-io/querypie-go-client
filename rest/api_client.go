package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
	AuthToken  string
}

func NewAPIClient(baseURL, authToken string) *APIClient {
	return &APIClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
		AuthToken:  authToken,
	}
}

func (client *APIClient) doRequest(method, endpoint string, body interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s%s", client.BaseURL, endpoint)

	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("error marshaling request body: %v", err))
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error creating request: %v", err))
	}

	req.Header.Set("Authorization", "Bearer "+client.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error sending request: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("error: received non-200 status code %d", resp.StatusCode))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error reading response body: %v", err))
	}

	return respBody, nil
}

func (client *APIClient) GetData(endpoint string) (map[string]interface{}, error) {
	respBody, err := client.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, errors.New(fmt.Sprintf("error unmarshaling JSON response: %v", err))
	}

	return result, nil
}

func (client *APIClient) PostData(endpoint string, data map[string]interface{}) (map[string]interface{}, error) {
	respBody, err := client.doRequest("POST", endpoint, data)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, errors.New(fmt.Sprintf("error unmarshaling JSON response: %v", err))
	}

	return result, nil
}
