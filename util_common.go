package ssql

import (
	"encoding/base64"
	"reflect"
	"strings"

	"github.com/samber/lo"
)

func base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func base64Decode(input string) (string, error) {
	str, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", nil
	}
	return string(str), nil
}

func reverse[T any](args *[]T) *[]T {
	output := lo.Reverse(*args)
	return &output
}

func popArray[T any](slice *[]T) *[]T {
	if slice == nil {
		return nil
	}
	if len(*slice) < 1 {
		return slice
	}
	arr := append((*slice)[:len(*slice)-1], (*slice)[len(*slice):]...)
	return &arr
}

func ExtractStructMappings(tags []string, s interface{}) (TagMappings, TagMappings) {
	t := reflect.TypeOf(s)
	mappingsByFields := TagMappings{}
	mappingsByTags := TagMappings{}
	tags = append(tags, "relation")
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		for _, tag := range tags {
			col := field.Tag.Get(tag)
			if mappingsByFields[tag] == nil {
				mappingsByFields[tag] = make(map[string]string)
			}
			if mappingsByTags[tag] == nil {
				mappingsByTags[tag] = make(map[string]string)
			}
			if col != "" {
				mappingsByFields[tag][field.Name] = strings.Split(col, ",")[0]
				mappingsByTags[tag][strings.Split(col, ",")[0]] = field.Name
			}
		}
	}
	return mappingsByTags, mappingsByFields
}

type TagMappings map[string]map[string]string // [tagType][structField]tagValue.

func (mappings TagMappings) GetTag(tag, field string) string {
	if tag == "" || field == "" || mappings[tag] == nil {
		return ""
	}
	return mappings[tag][field]
}
