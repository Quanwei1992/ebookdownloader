package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testbi = BookInfo{
	Name:     "我是谁",
	Author:   "sndnvaps",
	Chapters: C,
}

var C = []Chapter{
	{
		Title:   "第一章",
		Content: "这是第一章\r\n内容测试\r\n",
		Link:    "https://github.com/sndnvaps/ebookdownloader",
	},
	{
		Title:   "第二章",
		Content: "这是第二章\r\n内容测试\r\n",
		Link:    "https://github.com/sndnvaps/ebookdownloader",
	},
	{
		Title:   "第三章",
		Content: "这是第三章\r\n内容测试\r\n",
		Link:    "https://github.com/sndnvaps/ebookdownloader",
	},
}

func TestBookInfo(t *testing.T) {
	bookname := "我是谁"
	assert.Equal(t, bookname, testbi.Name)

	author := "sndnvaps"
	assert.Equal(t, author, testbi.Author)

	Size := len(testbi.Chapters)
	assert.Equal(t, 3, Size)

	link := "https://github.com/sndnvaps/ebookdownloader"
	assert.Equal(t, link, testbi.Chapters[0].Link)
}

func TestGenerateTxt(t *testing.T) {
	testbi.GenerateTxt()
	savename := "./outputs/" + testbi.Name + "-" + testbi.Author + ".txt"

	assert.True(t, true, isExist(savename))
	os.RemoveAll(savename)

}

func TestGenerateMobi(t *testing.T) {
	testbi.GenerateMobi()
	savename := "./outputs/" + testbi.Name + "-" + testbi.Author + ".mobi"

	assert.True(t, true, isExist(savename))
	//os.RemoveAll(savename)
}

// IsExist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
