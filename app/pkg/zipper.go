package pkg

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
)

func ZipFolder(folderPath string) ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip folders themselves (we'll infer structure from file paths)
		if info.IsDir() {
			return nil
		}

		// Create zip path relative to the base folder
		relPath, err := filepath.Rel(folderPath, path)
		if err != nil {
			return err
		}

		fileInZip, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		fileOnDisk, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fileOnDisk.Close()

		_, err = io.Copy(fileInZip, fileOnDisk)
		return err
	})

	if err != nil {
		return nil, err
	}

	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
