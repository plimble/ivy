package ivy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
