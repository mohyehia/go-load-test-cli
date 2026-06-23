package config

import (
	"fmt"
	"strings"
)

/*
To make the standard flag package support this, we must implement Go's built-in flag.Value interface.
Any custom type that implements this interface can be treated as a repeatable flag!
The interface looks like this under the hood:
type Value interface {
    String() string
    Set(string) error
}
*/

type headerFlags []string

func (h *headerFlags) String() string {
	return fmt.Sprintf("%v", *h)
}

func (h *headerFlags) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func extractHeaderFlags(rawHeaders headerFlags) (map[string]string, error) {
	kvMap := make(map[string]string)
	// parse rawHeaders
	for _, header := range rawHeaders {
		kv := strings.SplitN(header, ":", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid header format '%s'. Must be 'Key: Value'", header)
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		kvMap[key] = value
	}
	return kvMap, nil
}
