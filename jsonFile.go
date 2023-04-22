package ebookdownloader

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// LoadBookJSONData 从文件中读取BookInfo信息，并返回
func LoadBookJSONData(filename string) (BookInfo, error) {
	fileAbs, _ := filepath.Abs(filename)
	filePtr, _ := os.Open(fileAbs)

	defer filePtr.Close()

	var bookinfo BookInfo

	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err := decoder.Decode(&bookinfo)
	if err != nil {
		return BookInfo{}, err
	}
	return bookinfo, nil
}
