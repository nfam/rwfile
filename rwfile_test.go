package rwfile

import (
	"bytes"
	"os"
	"testing"
)

func TestReadWrite(t *testing.T) {
	const name = "rwfile_test_gen.txt"

	f, err := OpenReadWrite(name)
	if err != nil {
		t.Error(err)
		return
	}
	if _, err = f.WriteString("read-write"); err != nil {
		f.Close()
		t.Error(err)
		return
	}
	f.Close()

	f, err = OpenRead(name)
	if err != nil {
		t.Error(err)
		return
	}
	var b bytes.Buffer
	if _, err = b.ReadFrom(f); err != nil {
		f.Close()
		t.Error(err)
		return
	}
	f.Close()
	if b.String() != "read-write" {
		t.Error("reault " + b.String())
	}

	_ = os.Remove(name)
}
