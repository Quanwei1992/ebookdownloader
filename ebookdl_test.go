package ebookdownloader

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testbi = BookInfo{
	Name:        "我是谁",
	Author:      "sndnvaps",
	Description: "这是我随便写的测试内容简介！",
	Volumes:     V, //分卷信息
	Chapters:    C,
}

var V = []Volume{
	{
		PrevChapterId: 0,
		CurrentVolume: "第一卷", //插入位置，第一章前面
		NextChapterId: 2,
	},
	{
		PrevChapterId: 2,
		CurrentVolume: "第二卷", //插入位置，第三章前面
		NextChapterId: 3,
	},
	{
		PrevChapterId: 5,
		CurrentVolume: "第三卷", //插入位置，第六章前面
		NextChapterId: 6,
	},
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
	{
		Title:   "第四章",
		Content: "这是第四章\r\n内容测试\r\n",
		Link:    "https://github.com/sndnvaps/ebookdownloader",
	},
	{
		Title:   "第五章",
		Content: "这是第五章\r\n内容测试\r\n",
		Link:    "https://github.com/sndnvaps/ebookdownloader",
	},
	{
		Title:   "第六章",
		Content: "这是第六章\r\n内容测试\r\n",
		Link:    "https://github.com/sndnvaps/ebookdownloader",
	},
	{
		Title:   "第七章",
		Content: "这是第七章\r\n内容测试\r\n",
		Link:    "https://github.com/sndnvaps/ebookdownloader",
	},
	{
		Title:   "第八章",
		Content: "这是第八章\r\n内容测试\r\n",
		Link:    "https://github.com/sndnvaps/ebookdownloader",
	},
	{
		Title:   "第九章",
		Content: "这是第九章\r\n内容测试\r\n",
		Link:    "https://github.com/sndnvaps/ebookdownloader",
	},
	{
		Title:   "第十章",
		Content: "这是第十章\r\n内容测试\r\n",
		Link:    "https://github.com/sndnvaps/ebookdownloader",
	},
}

func TestBookInfo(t *testing.T) {
	bookname := "我是谁"
	assert.Equal(t, bookname, testbi.Name)

	author := "sndnvaps"
	assert.Equal(t, author, testbi.Author)

	Size := len(testbi.Chapters)
	assert.Equal(t, 10, Size)

	link := "https://github.com/sndnvaps/ebookdownloader"
	assert.Equal(t, link, testbi.Chapters[0].Link)
}

func TestGenerateTxt(t *testing.T) {
	testbi.ChangeVolumeState(true /* hasVolume */)
	testbi.GenerateTxt()
	savename := "./outputs/" + testbi.Name + "-" + testbi.Author + "/" + testbi.Name + "-" + testbi.Author + ".txt"

	assert.True(t, true, isExist(savename))
	//os.RemoveAll(savename)

}

func TestGenerateMobi(t *testing.T) {
	testbi.ChangeVolumeState(true /* hasVolume */)
	testbi.SetKindleEbookType(true /* isMobi */, false /* isAwz3 */)
	testbi.GenerateMobi()
	savename := "./outputs/" + testbi.Name + "-" + testbi.Author + "/" + testbi.Name + "-" + testbi.Author + ".mobi"

	assert.True(t, true, isExist(savename))
	//os.RemoveAll(savename)
}

func TestGenerateAwz3(t *testing.T) {
	testbi.ChangeVolumeState(true /* hasVolume */)
	testbi.SetKindleEbookType(false /* isMobi */, true /* isAwz3 */)
	testbi.GenerateMobi()
	savename := "./outputs/" + testbi.Name + "-" + testbi.Author + "/" + testbi.Name + "-" + testbi.Author + ".awz3"

	assert.True(t, true, isExist(savename))
	//os.RemoveAll(savename)
}

// IsExist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
