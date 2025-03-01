package rwfile

import (
	"os"
	"path/filepath"
)

type File struct {
	*os.File
}

func (f File) Close() {
	if f.File != nil {
		unlock(f.File)
		f.File.Close()
	}
}

func OpenRead(name string) (File, error) {
	f, err := os.Open(name)
	if err != nil {
		return File{}, err
	}
	if err = rlock(f); err != nil {
		f.Close()
		return File{}, err
	}
	return File{f}, nil
}

func OpenWrite(name string) (File, error) {
	if err := os.MkdirAll(filepath.Dir(name), os.ModePerm); err != nil {
		return File{}, err
	}
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return File{}, err
	}
	if err = wlock(f); err != nil {
		f.Close()
		return File{}, err
	}
	return File{f}, nil
}

func OpenReadWrite(name string) (File, error) {
	if err := os.MkdirAll(filepath.Dir(name), os.ModePerm); err != nil {
		return File{}, err
	}
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return File{}, err
	}
	if err = wlock(f); err != nil {
		f.Close()
		return File{}, err
	}
	return File{f}, nil
}
