package pkg

import (
	"fmt"
	"os"
	"path/filepath"
)

func LoadFile(source, fileName string, target any, decode func(data []byte, v any) error) error {
	if source == "" {
		return fmt.Errorf("source is empty : %s", source)
	} else if fileName == "" {
		return fmt.Errorf("file is empty : %s", fileName)
	}
	// check if src exits
	fullPath := filepath.Join(source, fileName)
	if _, err := os.Stat(fullPath); err != nil && os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist at path %s, err : %w", fileName, fullPath, err)
	} else if err != nil {
		return err
	}

	// read from file
	dataBytes, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}
	if err = decode(dataBytes, target); err != nil {
		return fmt.Errorf("decode file %s to target err: %w", fullPath, err)
	}
	return nil
}

