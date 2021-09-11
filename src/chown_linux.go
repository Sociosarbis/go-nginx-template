// +build linux

package template

import (
	"os"
	"syscall"
)

func Chown(dst *os.File, src os.FileInfo) error {
	return dst.Chown(int(src.Sys().(*syscall.Stat_t).Uid), int(src.Sys().(*syscall.Stat_t).Gid))
}

func Chown(dst *os.File, src os.FileInfo) error {
	return dst.Chmod(src.Mode())
}
