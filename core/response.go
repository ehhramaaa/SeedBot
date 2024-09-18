package core

import (
	"encoding/json"
	"fmt"

)

func handleResponseMap(respBody []byte) (map[string]interface{}, error) {
	// Mengurai JSON ke dalam map[string]interface{}
	var result map[string]interface{}
	err := json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	return result, nil
}

func handleResponseArray(respBody []byte) ([]map[string]interface{}, error) {
	// Variabel untuk menampung hasil unmarshalling
	var result []map[string]interface{}

	// Mengurai JSON ke dalam variabel 'result'
	err := json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	return result, nil
}
