package ebookdownloader

import (
	"encoding/json"
	"os"
	"path/filepath"
)

//Meta 小说的metainfo
type Meta struct {
	Ebhost      string `json:"ebhost"`
	Bookid      string `json:"bookid"`
	BookName    string `json:"bookname"`
	Author      string `json:"author"`
	CoverUrl    string `json:"cover_url"`
	Description string `json:"description"`
	TxtUrlPath  string `json:"txt_url_path"`
	MobiUrlPath string `json:"mobi_url_path"`
}

//WriteFile 把json数据写入 filename定义的文件中
func (this Meta) WriteFile(filename string) error {
	// 创建文件
	fileAbs, _ := filepath.Abs(filename) //使用相对路径
	filePtr, err := os.Create(fileAbs)
	if err != nil {
		return err
	}
	defer filePtr.Close()

	// 带JSON缩进格式写文件
	data, err := json.MarshalIndent(this, "", "  ")
	if err != nil {
		return err
	}

	//写入文件中
	filePtr.Write(data)

	return nil
}

//GetMetaData 从文件中读取meta信息，并返回
func GetMetaData(filename string) (Meta, error) {
	fileAbs, _ := filepath.Abs(filename)
	filePtr, _ := os.Open(fileAbs)

	defer filePtr.Close()

	var metainfo Meta

	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err := decoder.Decode(&metainfo)
	if err != nil {
		return Meta{}, err
	}
	return metainfo, nil
}
