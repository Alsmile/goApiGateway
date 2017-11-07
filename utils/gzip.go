package utils

import (
	"compress/gzip"
	"errors"
	"io"
	"os"
	"sync"

	"cloudtogo.cn/webserver/services"
)

type compressionPool struct {
	sync.Pool
	Level int
}

var gzipPool = &compressionPool{Level: gzip.BestCompression}

// Compress 使用gzip压缩成.gz
func Compress(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gzipWriter := getGzipWriter(dstFile)
	if gzipWriter == nil {
		return errors.New(services.ErrorProjectMaxSize)
	}
	defer gzipWriter.Close()

	_, err = io.Copy(gzipWriter, srcFile)

	if err != nil {
		releaseGzipWriter(gzipWriter)
		return err
	}
	err = gzipWriter.Flush()
	releaseGzipWriter(gzipWriter)

	return nil
}

func getGzipWriter(w io.Writer) *gzip.Writer {
	v := gzipPool.Get()
	if v == nil {
		gzipWriter, err := gzip.NewWriterLevel(w, gzipPool.Level)
		if err != nil {
			return nil
		}
		return gzipWriter
	}
	gzipWriter := v.(*gzip.Writer)
	gzipWriter.Reset(w)
	return gzipWriter
}

func releaseGzipWriter(gzipWriter *gzip.Writer) {
	gzipWriter.Close()
	gzipPool.Put(gzipWriter)
}
