package shorten

import (
	"strings"
)

const (
	dictionary    = "abcdefghjkmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
	lenDictionary = uint32(len(dictionary))
)

func GenerateShortenUrl(id uint32) string {
	var (
		builder strings.Builder
		nums    []uint32
	)

	for id > 0 {
		nums = append(nums, id%lenDictionary)
		id /= lenDictionary
	}

	for _, num := range nums {
		builder.WriteByte(dictionary[num])
	}

	return builder.String()
}
