package ebookdownloader

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

// CreateMD5 根据filename生成md5信息
func CreateMD5(filename string) (md5Str string, err error) {

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Open", err)
		return "", err
	}

	defer f.Close()

	md5hash := md5.New()
	if _, err = io.Copy(md5hash, f); err != nil {
		fmt.Println("Copy", err)
		return "", err
	}

	md5hash.Sum(nil)
	md5Str = fmt.Sprintf("%x", md5hash.Sum(nil))
	return md5Str, nil
}
