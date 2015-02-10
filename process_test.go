package ivy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetGravity(t *testing.T) {
	assert.Equal(t, "NorthWest", getGravity("nw"))
	assert.Equal(t, "North", getGravity("n"))
	assert.Equal(t, "NorthEast", getGravity("ne"))
	assert.Equal(t, "West", getGravity("w"))
	assert.Equal(t, "Center", getGravity("c"))
	assert.Equal(t, "East", getGravity("e"))
	assert.Equal(t, "SouthWest", getGravity("sw"))
	assert.Equal(t, "South", getGravity("s"))
	assert.Equal(t, "SouthEast", getGravity("se"))
	assert.Equal(t, "NorthWest", getGravity("xx"))
}
