package rwfile

import (
	"bytes"
	"os"
	"time"

	"github.com/nfam/pool/buffer"
)

func Write(name string, data []byte, mtime *time.Time) (int, error) {
	f, err := OpenReadWrite(name)
	if err != nil {
		return 0, err
	}
	n, err := write(f.File, data)
	f.Close()

	if mtime != nil {
		_ = os.Chtimes(name, *mtime, *mtime)
	}
	return n, err
}

func write(f *os.File, data []byte) (int, error) {
	b := buffer.Get()
	defer b.Close()

	// Read old bytes.
	if _, err := b.ReadFrom(f); err != nil {
		f.Close()
		return 0, err
	}
	old := b.Bytes()
	if bytes.Equal(data, old) {
		f.Close()
		return 0, nil
	}

	// Find diff.
	var doff []int
	var dlen []int
	{
		var (
			mlen = min(len(data), len(old))
			ioff int
			ilen int
		)
		for i := range mlen {
			if data[i] != old[i] {
				if ilen == 0 {
					ioff = i
					ilen = 1
				} else {
					ilen++
				}
			} else if ilen > 0 {
				doff = append(doff, ioff)
				dlen = append(dlen, ilen)
				ilen = 0
			}
		}
		if mlen < len(data) {
			if ilen == 0 {
				ioff = mlen
				ilen = len(data) - mlen
			} else {
				ilen += len(data) - mlen
			}
		}
		if ilen > 0 {
			doff = append(doff, ioff)
			dlen = append(dlen, ilen)
			ilen = 0
		}
	}

	// Patch diff.
	var count int
	for i := range doff {
		ioff := doff[i]
		ilen := dlen[i]

		n, err := f.WriteAt(data[ioff:ioff+ilen], int64(ioff))
		count += n
		if err != nil {
			return count, err
		}
	}
	if len(data) < len(old) {
		if err := f.Truncate(int64(len(data))); err != nil {
			return count, err
		}
	}
	err := f.Sync()
	return count, err
}
