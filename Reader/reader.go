package Reader

import (
	"errors"
	"os"
)

func ReadFile(filePath string) ([]byte, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.New("can't read a file")
	}
	return file, err
}
