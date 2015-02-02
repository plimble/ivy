package fileproxy

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"image"
	"image/png"
	"testing"
)

func TestProcessResize(t *testing.T) {
	params, _ := parseParams("r_100x100")

	buffer := new(bytes.Buffer)

	png.Encode(buffer, image.NewRGBA(image.Rect(0, 0, 200, 200)))

	result, err := process(params, "test.png", buffer)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestProcessCrop(t *testing.T) {
	params, _ := parseParams("c_100x100,g_c")

	buffer := new(bytes.Buffer)

	png.Encode(buffer, image.NewRGBA(image.Rect(0, 0, 200, 200)))

	result, err := process(params, "test.png", buffer)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestProcessResizeCrop(t *testing.T) {
	params, _ := parseParams("r_150x150,c_100x100,g_c")

	buffer := new(bytes.Buffer)

	png.Encode(buffer, image.NewRGBA(image.Rect(0, 0, 200, 200)))

	result, err := process(params, "test.png", buffer)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func BenchmarkResizeCrop(b *testing.B) {
	params, _ := parseParams("r_350x350,c_200x200,g_c")

	buffer := new(bytes.Buffer)

	png.Encode(buffer, image.NewRGBA(image.Rect(0, 0, 700, 700)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		process(params, "test.png", buffer)
	}
}
