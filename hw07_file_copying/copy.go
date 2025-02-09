package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb" //nolint
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	var source io.Reader
	var sourceSize int64

	s, err := os.Open(fromPath)
	if err != nil {
		fmt.Printf("Can't open %s: %v\n", fromPath, err)
		return ErrUnsupportedFile
	}

	defer func(s *os.File) {
		err := s.Close()
		if err != nil {
			fmt.Printf("Can't close %s: %v\n", fromPath, err)
		}
	}(s)

	sourceStat, err := s.Stat()
	if err != nil {
		fmt.Printf("Can't stat %s: %v\n", fromPath, err)
		return ErrUnsupportedFile
	}

	sourceSize = sourceStat.Size()

	if offset > 0 {
		if offset > sourceSize {
			return ErrOffsetExceedsFileSize
		}
		_, err = s.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
	}

	source = s

	dest, err := os.Create(toPath)
	if err != nil {
		fmt.Printf("Can't create %s: %v\n", toPath, err)
		return ErrUnsupportedFile
	}

	defer func(dest *os.File) {
		err := dest.Close()
		if err != nil {
			fmt.Printf("Can't close %s: %v\n", toPath, err)
		}
	}(dest)

	copySize := int(limit)

	if limit == 0 {
		copySize = int(sourceSize)
	}

	bar := pb.New(copySize)
	bar.SetUnits(pb.U_BYTES)
	bar.SetRefreshRate(time.Millisecond * 10)
	bar.ShowSpeed = true
	bar.Start()

	if limit > 0 {
		source = io.LimitReader(source, limit)
	}

	reader := bar.NewProxyReader(source)

	io.Copy(dest, reader)

	bar.Finish()

	return nil
}
