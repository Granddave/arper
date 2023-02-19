package utils

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestSerializeDeserialize(t *testing.T) {
	// Arrange
	tempFile, err := ioutil.TempFile("", "serializer")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	type S struct {
		S string
		I int
	}

	orig := []S{
		{S: "a", I: 1},
		{S: "b", I: 2},
	}

	// Act
	err = Serialize(orig, tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	serialized := []S{}
	if err := Deserialize(&serialized, tempFile.Name()); err != nil {
		t.Fatal(err)
	}

	// Assert
	if len(orig) != len(serialized) {
		t.Errorf("Slice lengths do not match: %d != %d", len(orig), len(serialized))
	}

	for i := range orig {
		o := &orig[i]
		s := &serialized[i]
		if o.S != s.S {
			t.Errorf("Strings do not match: %s != %s", o.S, s.S)
		}
		if o.I != s.I {
			t.Errorf("Ints do not match: %d != %d", o.I, s.I)
		}
	}
}
