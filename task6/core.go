package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type writerWithLog struct {
	target   io.Writer
	total    uint64
	progress uint64
}

func (o *writerWithLog) Write(p []byte) (int, error) {
	n, err := o.target.Write(p)
	if err == nil && n > 0 {
		o.progress += uint64(n)
		fmt.Printf("\rProgress: %d%% ", (o.progress*100)/o.total)
	}
	return n, err
}

func (o *writerWithLog) FinishProgress() {
	fmt.Println()
}

func fileCopy(src, dst string, offset, limit int64) error {
	var err error

	var srcFile *os.File
	srcFile, err = os.Open(src)
	if err != nil {
		log.Printf("Fail to open source file: %s\n", err.Error())
		return err
	}
	defer srcFile.Close()

	var srcFileSize uint64

	// stat of src
	var srcFileInfo os.FileInfo
	srcFileInfo, err = srcFile.Stat()
	if err == nil {
		srcFileSize = uint64(srcFileInfo.Size())
	} else {
		log.Printf("Fail to get stat of source file: %s\n", err.Error())
	}

	var dstFile *os.File
	dstFile, err = os.Create(dst)
	if err != nil {
		log.Printf("Fail to open destination file: %s\n", err.Error())
		return err
	}
	defer dstFile.Close()

	dstWrapper := &writerWithLog{target: dstFile, total: srcFileSize}

	if offset > 0 {
		_, err = srcFile.Seek(offset, io.SeekStart)
		if err != nil {
			log.Printf("Fail to seek source file: %s\n", err.Error())
			return err
		}
	}

	if limit <= 0 {
		_, err = io.Copy(dstWrapper, srcFile)
	} else {
		_, err = io.CopyN(dstWrapper, srcFile, limit)
		if err == io.EOF { // if out of the file size
			err = nil
		}
	}
	if err != nil {
		log.Printf("Fail to copy file: %s\n", err.Error())
		return err
	}

	dstWrapper.FinishProgress()

	return nil
}
