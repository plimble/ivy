package fileproxy

// import (
// 	"github.com/stretchr/testify/assert"
// 	"testing"
// )

// func TestParseParams(t *testing.T) {
// 	assert := assert.New(t)

// 	params, err := parseParams("")
// 	assert.NoError(err)
// 	assert.Equal(&Params{0, 0, DefaultCropMode, DefaultCropPos, DefaultScale, DefaultQuality, true, "", true}, params)

// 	params, err = parseParams("_")
// 	assert.NoError(err)
// 	assert.Equal(&Params{0, 0, DefaultCropMode, DefaultCropPos, DefaultScale, DefaultQuality, true, "", true}, params)

// 	params, err = parseParams("w_100,h_200,c_e,p_tc,s_2")
// 	assert.NoError(err)
// 	assert.Equal(&Params{100, 200, "e", "tc", 2, 100, false, "", false}, params)

// 	params, err = parseParams("w_100,h_200")
// 	assert.NoError(err)
// 	assert.Equal(&Params{100, 200, "", "", 1, 100, false, "", false}, params)

// 	params, err = parseParams("w_0,h_200,c_tc,s_2")
// 	assert.Error(err)
// 	assert.Nil(params)

// 	params, err = parseParams("w_100,h_0,c_tc,s_2")
// 	assert.Error(err)
// 	assert.Nil(params)

// 	params, err = parseParams("w_100,h_0,c_xc,s_2")
// 	assert.Error(err)
// 	assert.Nil(params)

// 	params, err = parseParams("w_100,h_0,c_tc,s_0")
// 	assert.Error(err)
// 	assert.Nil(params)

// 	params, err = parseParams("te")
// 	assert.Error(err)
// 	assert.Nil(params)

// 	params, err = parseParams("t_est")
// 	assert.Error(err)
// 	assert.Nil(params)
// }

// func BenchmarkParseParams(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		parseParams("w_100,h_200,c_e,p_tc,s_2")
// 	}
// }

// func BenchmarkParamsToString(b *testing.B) {
// 	params, _ := parseParams("w_100,h_200,c_e,p_tc,s_2")
// 	for i := 0; i < b.N; i++ {
// 		params.String()
// 	}
// }
