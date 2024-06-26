package main

import (
	"errors"
	"io"
	"io/ioutil" //nolint:all
	"os"

	"github.com/cheggaaa/pb/v3" //nolint:all
)

const tmpFilenamePattern = "temp.*"

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

type params struct {
	offset int64
	limit  int64
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	params, err := validateArguments(fromPath, offset, limit)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fromPath, os.O_RDONLY, 0o666)
	if err != nil {
		return err
	}
	defer file.Close()

	if offset > 0 {
		_, err = file.Seek(params.offset, io.SeekStart)
		if err != nil {
			return err
		}
	}

	tmpFile, err := ioutil.TempFile("", tmpFilenamePattern)
	if err != nil {
		return err
	}

	err = copyFile(file, tmpFile, params.limit)
	if err != nil {
		return err
	}

	err = os.Rename(tmpFile.Name(), toPath)
	if err != nil {
		return err
	}

	return nil
}

func validateArguments(filePath string, offset int64, limit int64) (*params, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	if !fileInfo.Mode().IsRegular() {
		return nil, ErrUnsupportedFile
	}

	if fileInfo.Size() < offset {
		return nil, ErrOffsetExceedsFileSize
	}

	if limit > (fileInfo.Size() - offset) {
		limit = fileInfo.Size() - offset
	}

	return &params{
		offset: offset,
		limit:  limit,
	}, nil
}

func copyFile(src io.Reader, dst io.Writer, limit int64) error {
	pBar := pb.Start64(limit)
	defer pBar.Finish()

	var total int64
	for {
		n, err := io.CopyN(dst, src, 1)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return err
		}

		total += n
		if total == limit {
			break
		}
		pBar.Increment()
	}

	return nil
}
