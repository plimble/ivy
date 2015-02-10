package ivy

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"image"
	"image/png"
	"testing"
)

func TestGMProcess(t *testing.T) {
	buffer := new(bytes.Buffer)
	png.Encode(buffer, image.NewRGBA(image.Rect(0, 0, 200, 200)))

	gm := NewGMProcessor()
	params, _ := ParseParams("r_100x100,c_50x50,g_c,q_50")
	img, err := gm.Process(params, "text.png", buffer)
	assert.NoError(t, err)
	assert.NotNil(t, img)
	assert.True(t, img.Len() > 0)

	params, _ = ParseParams("r_100x0,c_50x50,g_c,q_50")
	img, err = gm.Process(params, "text.png", buffer)
	assert.NoError(t, err)
	assert.NotNil(t, img)
	assert.True(t, img.Len() > 0)

	params, _ = ParseParams("r_0x100,c_50x50,g_c,q_50")
	img, err = gm.Process(params, "text.png", buffer)
	assert.NoError(t, err)
	assert.NotNil(t, img)
	assert.True(t, img.Len() > 0)
}

func TestGMGetGravity(t *testing.T) {
	gm := NewGMProcessor()
	assert.Equal(t, "NorthWest", gm.getGravity("nw"))
	assert.Equal(t, "North", gm.getGravity("n"))
	assert.Equal(t, "NorthEast", gm.getGravity("ne"))
	assert.Equal(t, "West", gm.getGravity("w"))
	assert.Equal(t, "Center", gm.getGravity("c"))
	assert.Equal(t, "East", gm.getGravity("e"))
	assert.Equal(t, "SouthWest", gm.getGravity("sw"))
	assert.Equal(t, "South", gm.getGravity("s"))
	assert.Equal(t, "SouthEast", gm.getGravity("se"))
	assert.Equal(t, "NorthWest", gm.getGravity("xx"))
}
