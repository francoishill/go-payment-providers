package PayfastUIConfig

import (
	"fmt"
	"net/url"
	"strings"
)

type keyValPair struct {
	Key   string
	Value string
}

func appendKeyValToStringSlice(strSlice []string, key, val string) []string {
	valEncoded := url.QueryEscape(strings.Trim(val, " "))
	if strings.Trim(valEncoded, " ") == "" {
		return strSlice
	}

	return append(strSlice, fmt.Sprintf("%s=%s", key, valEncoded))
}

func (this *keyValPair) appendKeyValToStringSlice(strSlice []string) []string {
	return appendKeyValToStringSlice(strSlice, this.Key, this.Value)
}
