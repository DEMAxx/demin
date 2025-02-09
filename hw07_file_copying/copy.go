package main

import (
	"errors"
	"fmt"
	"github.com/cheggaaa/pb"
	"io"
	"os"
	"time"
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

	// create dest
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
	//if offset > 0 {
	//	buf := make([]byte, limit)
	//	_, err = CopyBufferAt(dest, reader, offset, buf)
	//
	//	if err != nil {
	//		return err
	//	}
	//} else {
	//	if limit > 0 {
	//		io.CopyN(dest, reader, limit)
	//	} else {
	//		io.Copy(dest, reader)
	//	}
	//}

	bar.Finish()

	return nil
}

func CopyBufferAt(dest io.WriterAt, src io.Reader, off int64, buf []byte) (written int64, err error) {
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dest.WriteAt(buf[0:nr], off+written)
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errors.New("invalid write result")
				}
			}
			written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
