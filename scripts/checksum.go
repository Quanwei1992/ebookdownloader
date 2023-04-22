package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
)

// 对文件生成md5验证信息
func md5f(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

func sha1f(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()

	h := sha1.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

func sha256f(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()

	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

func sha384f(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()

	h := sha512.New384()
	_, err = io.Copy(h, f)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

func sha512f(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()

	h := sha512.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

// Output 把得到的验证信息，全部写入到文件中
func Output(infile, outfile string) {
	w, err := os.Create(outfile)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer w.Close()

	w.WriteString("MD5: " + md5f(infile) + "\n")
	w.WriteString("SHA1: " + sha1f(infile) + "\n")
	w.WriteString("SHA256: " + sha256f(infile) + "\n")
	w.WriteString("SHA384 " + sha384f(infile) + "\n")
	w.WriteString("SHA512: " + sha512f(infile) + "\n")
}
func main() {
	infile := flag.String("infile", "", "需要生成验证信息的文件名")
	outfile := flag.String("outfile", "", "输入验证信息到指定的文件名")
	flag.Parse()
	if *infile == "" {
		flag.Usage()
		os.Exit(0)
	}
	if *infile != "" && *outfile == "" {
		*outfile = *infile + ".hash"
	}
	Output(*infile, *outfile)

}
