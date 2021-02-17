package agent

import (
	"testing"
)

func TestJsonUnMarshal(t *testing.T) {
	// arrange
	var f []byte = nil

	// act
	_, err := JsonUnmarshal(f)

	// assert
	if err == nil {
		t.Fatal("Unmarshal of the nil type should fail.")
	}
}
