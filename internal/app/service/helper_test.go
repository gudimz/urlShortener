package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testID         = 1114539422
	testShortenURL = "Q6T8Ib"
)

func TestHelper_generateShortenURL(t *testing.T) {
	testCases := []struct {
		name string
		id   uint32
		want string
	}{
		{
			name: "return an alphanumeric identifier",
			id:   testID,
			want: testShortenURL,
		},
		{
			name: "return empty identifier",
			id:   0,
			want: "",
		},
	}

	for _, tt := range testCases {
		got := generateShortenURL(tt.id)
		assert.Equal(t, tt.want, got)
	}
}
