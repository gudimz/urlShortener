package shorten

import (
	"github.com/gofiber/fiber/v2/utils"
	"testing"
)

type testCase struct {
	id       uint32
	expected string
	log      string
}

func Test_CreateShortenUrl(t *testing.T) {

	t.Run("Returned short url", func(t *testing.T) {
		testCases := []testCase{
			{
				id:       0,
				expected: "",
			},
			{
				id:       2420807732,
				expected: "HJbWyd",
			},
			{
				id:       lenDictionary,
				expected: "ab",
			},
		}

		for i, tCase := range testCases {
			t.Logf("Test #%d\nid:%d, expected:|%v|", i+1, tCase.id, tCase.expected)
			actual := GenerateShortenUrl(tCase.id)
			utils.AssertEqual(t, actual, tCase.expected)

		}

	})

}
