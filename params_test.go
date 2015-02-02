package fileproxy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseParams(t *testing.T) {
	assert := assert.New(t)

	params, err := parseParams("")
	assert.NoError(err)
	assert.Equal(newParams(), params)

	params, err = parseParams("_")
	assert.NoError(err)
	assert.Equal(newParams(), params)

	params, err = parseParams("r_100x200,c_100x100,g_n,q_100")
	assert.NoError(err)
	assert.Equal(&Params{100, 200, 100, 100, "n", 100, true, true, true, false, ""}, params)
}

func BenchmarkParseParams(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseParams("r_100x200,c_100x100,g_n,q_100")
	}
}

func BenchmarkParamsToString(b *testing.B) {
	params, _ := parseParams("r_100x200,c_100x100,g_n,q_100")
	for i := 0; i < b.N; i++ {
		params.String()
	}
}
