package ebookdownloader

import (
	"os"
	"testing"

	. "gopkg.in/check.v1"
)

var testbi = BookInfo{
	Name:           "我是谁",
	Author:         "sndnvaps",
	Description:    "这是我随便写的测试内容简介！",
	Volumes:        V, //分卷信息
	Chapters:       Chapters,
	DlCoverFromWeb: false, //使用直接生成的封面
}

var V = []Volume{
	{
		PrevChapterID: 0,
		CurrentVolume: "第一卷", //插入位置，第一章前面
		NextChapterID: 2,
	},
	{
		PrevChapterID: 2,
		CurrentVolume: "第二卷", //插入位置，第三章前面
		NextChapterID: 3,
	},
	{
		PrevChapterID: 5,
		CurrentVolume: "第三卷", //插入位置，第六章前面
		NextChapterID: 6,
	},
}
var Chapters = []Chapter{
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

var savePath = "./outputs/" + testbi.Name + "-" + testbi.Author

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestBookInfo(c *C) {

	bookname := "我是谁"

	c.Assert(bookname, Equals, testbi.Name)
	author := "sndnvaps"
	c.Assert(author, Equals, testbi.Author)
	size := len(testbi.Chapters)
	c.Assert(size, Equals, 10)

	link := "https://github.com/sndnvaps/ebookdownloader"
	c.Assert(link, Equals, testbi.Chapters[0].Link)
}

func (s *MySuite) TestGenerateTxt(c *C) {
	testbi.ChangeVolumeState(true /* hasVolume */)
	testbi.GenerateTxt()
	savename := savePath + "/" + testbi.Name + "-" + testbi.Author + ".txt"
	c.Assert(isExist(savename), Equals, true)
	os.RemoveAll(savePath)
}

func (s *MySuite) TestGenerateMobi(c *C) {
	testbi.ChangeVolumeState(true /* hasVolume */)
	testbi.SetKindleEbookType(true /* isMobi */, false /* isAwz3 */)
	testbi.GenerateISBN() //先生成ISBN码
	testbi.GenerateMobi()
	savename := savePath + "/" + testbi.Name + "-" + testbi.Author + ".mobi"
	c.Assert(isExist(savename), Equals, true)
	//os.RemoveAll(savePath)
}

func (s *MySuite) TestGenerateAzw3(c *C) {
	testbi.ChangeVolumeState(true /* hasVolume */)
	testbi.SetKindleEbookType(false /* isMobi */, true /* isAzw3 */)
	testbi.GenerateISBN() //先生成ISBN码
	testbi.GenerateMobi()
	savename := savePath + "/" + testbi.Name + "-" + testbi.Author + ".azw3"
	c.Assert(isExist(savename), Equals, true)
	os.RemoveAll(savePath)
}
func (s *MySuite) TestGenerateEPUB(c *C) {
	testbi.GenerateISBN() //先生成ISBN码
	testbi.GenerateEPUB()
	savename := savePath + "/" + testbi.Name + "-" + testbi.Author + ".epub"
	c.Assert(isExist(savename), Equals, true)
	os.RemoveAll(savePath)
}

// IsExist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
