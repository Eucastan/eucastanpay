package kafka

import "encoding/json"

func Decode[T any](data []byte) (*T, error) {
	var event T

	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}

	return &event, nil
}
