package ebookdownloader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testMeta = Meta{
	BookUUID: "0000-0000-0000-0000",
	BookISBN: "955-955-955-955",
	BookName: "我是谁-sndnvaps",
	Bookid:   "0000-0001",
	Author:   "sndnvaps",
}

var testMeta01 = Meta{
	BookUUID: "0000-0000-0001-0001",
	BookISBN: "955-955-955-966",
	BookName: "我不是我-sndnvaps",
	Bookid:   "0000-0002",
	Author:   "sndnvaps",
}

func TestBoltdb(t *testing.T) {
	var metainfo Meta
	var metainfos []Meta
	boltdb, err := InitBoltDB("./ebookdownloader_t.db")
	if err != nil {
		t.Logf("初始化数据库测试失败,原因%s", err.Error())
	}
	err = boltdb.Save(testMeta)
	if err != nil {
		t.Logf("保存数据到boltdb->ebookdownloader_t.db失败，原因:%s", err.Error())
	}
	err = boltdb.Save(testMeta01)
	if err != nil {
		t.Logf("保存数据到boltdb->ebookdownloader_t.db失败，原因:%s", err.Error())
	}
	metainfo, err = boltdb.FindOneByUUID("0000-0000-0000-0000")
	assert.Equal(t, metainfo.BookUUID, testMeta.BookUUID)
	metainfo, err = boltdb.FindOneByAuthor("sndnvaps")
	assert.Equal(t, metainfo.Author, testMeta.Author)

	metainfos, err = boltdb.FindAllByAuthor("sndnvaps")
	assert.Equal(t, 2, len(metainfos))

	err = boltdb.UpdateTXTInfo("0000-0000-0000-0000", "public/我是谁-sndnvaps/我是谁-sndnvaps.txt", "this-is-md5-string-for-test")
	if err != nil {
		t.Logf("更新数据到boltdb->ebookdownloader_t.db失败，原因:%s", err.Error())
	}
	metainfo, err = boltdb.FindOneByUUID("0000-0000-0000-0000")
	assert.Equal(t, "this-is-md5-string-for-test", metainfo.TxtMD5)

}
