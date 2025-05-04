package fileutils

import "os"

func SaveToFile(filename string, bytes []byte) error {
	err := os.WriteFile(filename, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReadFile(filename string) ([]byte, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
