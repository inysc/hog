//go:build !linux
// +build !linux

package hog

import (
	"os"
)

func chown(_ string, _ os.FileInfo) error {
	return nil
}
