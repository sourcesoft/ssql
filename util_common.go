package ssql

import (
	"encoding/base64"

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