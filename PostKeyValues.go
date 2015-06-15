package gopp

import (
	"fmt"
	"net/url"
	"strings"
)

type PostKeyValues struct {
	Key   string
	Value ValueString
}

func readKeyValuePairsInCorrectOrderFromPostBody(requestPostBody []byte) SliceOfPostKeyValues {
	finalKeyVals := []*PostKeyValues{}

	requestBodyString := string(requestPostBody)
	keyValPairs := strings.Split(requestBodyString, "&")
	for _, keyAndVal := range keyValPairs {
		tmpSplit := strings.Split(keyAndVal, "=")
		finalKeyVals = append(finalKeyVals, &PostKeyValues{
			Key:   tmpSplit[0],
			Value: ValueString(tmpSplit[1]),
		})
	}

	return SliceOfPostKeyValues(finalKeyVals)
}

type SliceOfPostKeyValues []*PostKeyValues

func (this SliceOfPostKeyValues) CombineIntoSingleString(mustEscape bool) string {
	keyValCombinedList := []string{}
	for _, keyVal := range this {
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
