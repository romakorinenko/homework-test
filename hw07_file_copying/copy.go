package main

import (
	"errors"
	"io"
	"os"

	"github.com/schollz/progressbar/v3" //nolint:depguard
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

var bufSize = 1024 * 256

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileForReading, fileSize, err := getInputFile(fromPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = fileForReading.Close()
	}()

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	_, err = fileForReading.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	remain := fileSize - offset
	toCopy := remain
	if limit > 0 && limit < remain {
		toCopy = limit
	}

	fileForWriting, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = fileForWriting.Close()
	}()

	buf := make([]byte, bufSize)

	progressBar := createProgressBar(fileSize)

	var copied int64
	for copied < toCopy {
		readSize := bufSize
		left := toCopy - copied

		if left < int64(readSize) {
			readSize = int(left)
		}

		n, readErr := fileForReading.Read(buf[:readSize])
		if readErr != nil {
			return readErr
		}

		if n > 0 {
			written, writeErr := fileForWriting.Write(buf[:n])
			if writeErr != nil {
				return writeErr
			}

			if written != n {
				return io.ErrShortWrite
			}
			copied += int64(written)

			err := progressBar.Set64(copied)
			if err != nil {
				return err
			}
		}
	}

	progressBar.Describe("Copied successfully")
	err = progressBar.Set64(fileSize)
	if err != nil {
		return err
	}

	return nil
}

func getInputFile(fromPath string) (*os.File, int64, error) {
	file, err := os.Open(fromPath)
	if err != nil {
		return nil, 0, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, 0, err
	}

	if !fileInfo.Mode().IsRegular() {
		return nil, 0, ErrUnsupportedFile
	}

	return file, fileInfo.Size(), nil
}

func createProgressBar(fileSize int64) *progressbar.ProgressBar {
	progressBar := progressbar.DefaultBytes(fileSize, "Copying...")
	return progressBar
}
