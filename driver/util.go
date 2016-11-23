package main

import (
	"encoding/json"
)

// parseOptions attempts to decode the json string `opts` into
// a FlexOptions structure
func parseOptions(opts string) (FlexOptions, error) {
	options := FlexOptions{}

	err := json.Unmarshal([]byte(opts), &options)

	if err != nil {
		return options, err
	}

	return options, nil
}
