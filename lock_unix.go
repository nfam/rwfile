//go:build !windows

package rwfile

import (
	"os"
	"syscall"
)

func rlock(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_SH)
}

func wlock(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
}

func unlock(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
}
