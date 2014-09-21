package fileproxy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseParams(t *testing.T) {
	assert := assert.New(t)

	params, err := ParseParams("")
	assert.NoError(err)
	assert.Equal(&Params{0, 0, DefaultCropMode, DefaultCropPos, DefaultScale, DefaultQuality, true}, params)

	params, err = ParseParams("_")
	assert.NoError(err)
	assert.Equal(&Params{0, 0, DefaultCropMode, DefaultCropPos, DefaultScale, DefaultQuality, true}, params)

	params, err = ParseParams("w_100,h_200,c_e,p_tc,s_2")
	assert.NoError(err)
	assert.Equal(&Params{100, 200, "e", "tc", 2, 100, false}, params)

	params, err = ParseParams("w_100,h_200")
	assert.NoError(err)
	assert.Equal(&Params{100, 200, "", "", 1, 100, false}, params)

	params, err = ParseParams("w_0,h_200,c_tc,s_2")
	assert.Error(err)
	assert.Nil(params)

	params, err = ParseParams("w_100,h_0,c_tc,s_2")
	assert.Error(err)
	assert.Nil(params)

	params, err = ParseParams("w_100,h_0,c_xc,s_2")
	assert.Error(err)
	assert.Nil(params)

	params, err = ParseParams("w_100,h_0,c_tc,s_0")
	assert.Error(err)
	assert.Nil(params)

	params, err = ParseParams("te")
	assert.Error(err)
	assert.Nil(params)

	params, err = ParseParams("t_est")
	assert.Error(err)
	assert.Nil(params)
}

func BenchmarkParseParams(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseParams("w_100,h_200,c_e,p_tc,s_2")
	}
}
