package rwfile

import (
	"bytes"
	_ "embed"
	"os"
	"strconv"
	"testing"
)

//go:embed write_test.txt
var write_test_txt []byte

func TestWrite(t *testing.T) {
	name := "write_test_gen.txt"

	lines := bytes.Split(write_test_txt, []byte{'\n'})
	for i, line := range lines {
		index := bytes.Index(line, []byte{' '})
		if index <= 0 {
			t.Error("missing count at " + strconv.Itoa(i+1))
			t.FailNow()
		}
		count, err := strconv.Atoi(string(line[:index]))
		if err != nil {
			t.Error("invalid count at " + strconv.Itoa(i+1))
			t.FailNow()
		}
		content := line[index+1:]
		if n, err := Write(name, content, nil); err != nil {
			t.Error(err)
			t.FailNow()
		} else if n != count {
			t.Error("count " + strconv.Itoa(count) + " != " + strconv.Itoa(n) + " at " + strconv.Itoa(i+1))
			t.FailNow()
		}
		data, _ := os.ReadFile(name)
		if !bytes.Equal(data, content) {
			t.Error("not equal at " + strconv.Itoa(i+1))
			t.FailNow()
		}
	}
	_ = os.Remove(name)
}
