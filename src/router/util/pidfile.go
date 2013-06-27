package util

import (
	"fmt"
	"os"
)

func WritePidFile(p string) error {
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = fmt.Fprintf(f, "%d", os.Getpid())
	if err != nil {
		return err
	}

	return nil
}
