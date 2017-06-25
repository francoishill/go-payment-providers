package payfast

import (
	"fmt"
	"net/url"
	"strings"
)

type postKeyValue struct {
	Key   string
	Value string
}

func readOrderedKeyValuePairsFromPostBody(requestPostBody []byte) postKeyValueSlice {
	finalKeyVals := []*postKeyValue{}

	requestBodyString := string(requestPostBody)
	keyValPairs := strings.Split(requestBodyString, "&")
	for _, keyAndVal := range keyValPairs {
		tmpSplit := strings.Split(keyAndVal, "=")
		finalKeyVals = append(finalKeyVals, &postKeyValue{
			Key:   tmpSplit[0],
			Value: strings.TrimSpace(tmpSplit[1]),
		})
	}

	return postKeyValueSlice(finalKeyVals)
}

type postKeyValueSlice []*postKeyValue

func (p postKeyValueSlice) Combine(mustEscape bool) string {
	keyValCombinedList := []string{}
	for _, keyVal := range p {
		var value string
		if mustEscape {
			value = url.QueryEscape(string(keyVal.Value))
		} else {
			value = string(keyVal.Value)
		}
		keyValCombinedList = append(keyValCombinedList, fmt.Sprintf("%s=%s", keyVal.Key, value)) //No need to Escape again like payfast example, they are still escaped
	}

	return strings.Join(keyValCombinedList, "&")
}
