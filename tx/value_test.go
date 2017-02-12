package tx

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHexData(t *testing.T) {
	assert := assert.New(t)

	orig := HexData("!@##sdfg")
	enc, err := json.Marshal(orig)
	assert.Nil(err)

	var parsed HexData
	err = json.Unmarshal(enc, &parsed)
	assert.Nil(err)

	assert.Equal(orig, parsed)
}

func TestRawValuer(t *testing.T) {
	assert := assert.New(t)

	orig := []byte("C9spw.e")
	val := NewValue(orig)
	enc, err := json.Marshal(val)
	assert.Nil(err)

	var parsed RawValue
	err = json.Unmarshal(enc, &parsed)
	assert.Nil(err)

	assert.EqualValues(orig, parsed.Value)
	assert.Equal(RawValueType, parsed.Type)
}
