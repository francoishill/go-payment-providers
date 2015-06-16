package gopp

import (
	"fmt"
	"net/url"
	"strings"
)

type PostKeyValue struct {
	Key   string
	Value ValueString
}

func readKeyValuePairsInCorrectOrderFromPostBody(requestPostBody []byte) SliceOfPostKeyValue {
	finalKeyVals := []*PostKeyValue{}

	requestBodyString := string(requestPostBody)
	keyValPairs := strings.Split(requestBodyString, "&")
	for _, keyAndVal := range keyValPairs {
		tmpSplit := strings.Split(keyAndVal, "=")
		finalKeyVals = append(finalKeyVals, &PostKeyValue{
			Key:   tmpSplit[0],
			Value: ValueString(tmpSplit[1]),
		})
	}

	return SliceOfPostKeyValue(finalKeyVals)
}

type SliceOfPostKeyValue []*PostKeyValue

func (this SliceOfPostKeyValue) CombineIntoSingleString(mustEscape bool) string {
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
