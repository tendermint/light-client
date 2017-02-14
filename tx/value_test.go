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

	bin := HexData{79, 3, 72, 4}
	be, err := json.Marshal(bin)
	assert.Nil(err)
	// make sure proper hex
	assert.Equal(be[0], byte('"'))
	assert.Equal(be[1], byte('4'))
	assert.Equal(be[2], byte('f'))

	err = json.Unmarshal(be, &parsed)
	assert.Nil(err)

	assert.Equal(bin, parsed)
}

// TODO: make better test
func TestBase58Data(t *testing.T) {
	assert := assert.New(t)

	orig := B58Data("!@##sdfg")
	enc, err := json.Marshal(orig)
	assert.Nil(err)

	var parsed B58Data
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
