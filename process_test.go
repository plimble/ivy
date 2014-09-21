package fileproxy

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"image"
	// "image/gif"
	// "image/jpeg"
	"image/png"
	"testing"
)

func TestProcessResize(t *testing.T) {
	params := &Params{
		Width:    100,
		Height:   100,
		CropMode: "",
		CropPos:  "",
		Scale:    1,
		Quality:  70,
	}

	buffer := new(bytes.Buffer)

	png.Encode(buffer, image.NewRGBA(image.Rect(0, 0, 200, 200)))

	result, err := process(params, "test.png", ".png", buffer)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestProcessCropExact(t *testing.T) {
	params := &Params{
		Width:    100,
		Height:   100,
		CropMode: "e",
		CropPos:  "tc",
		Scale:    1,
		Quality:  70,
	}

	buffer := new(bytes.Buffer)

	png.Encode(buffer, image.NewRGBA(image.Rect(0, 0, 200, 200)))

	result, err := process(params, "test.png", ".png", buffer)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestProcessCropScale(t *testing.T) {
	params := &Params{
		Width:    100,
		Height:   100,
		CropMode: "s",
		CropPos:  "tc",
		Scale:    1,
		Quality:  70,
	}

	buffer := new(bytes.Buffer)

	png.Encode(buffer, image.NewRGBA(image.Rect(0, 0, 200, 200)))

	result, err := process(params, "test.png", ".png", buffer)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func BenchmarkResize(b *testing.B) {
	params := &Params{
		Width:    500,
		Height:   500,
		CropMode: "",
		CropPos:  "",
		Scale:    1,
		Quality:  70,
	}

	buffer := new(bytes.Buffer)

	png.Encode(buffer, image.NewRGBA(image.Rect(0, 0, 700, 700)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		process(params, "test.png", ".png", buffer)
	}
}
