// +build windows

package template

import "os"

func Chown(dst *os.File, src os.FileInfo) error {
	return nil
}

func Chmod(dst *os.File, src os.FileInfo) error {
	return nil
}
