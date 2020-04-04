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

func TestInitBoltdb(t *testing.T) {
	var metainfo Meta
	boltdb, err := InitBoltDB("./ebookdownloader_t.db")
	if err != nil {
		t.Logf("初始化数据库测试失败,原因%s", err.Error())
	}
	err = boltdb.Save(testMeta)
	if err != nil {
		t.Logf("保存数据到boltdb->ebookdownloader.db失败，原因:%s", err.Error())
	}
	metainfo, err = boltdb.FindOneByUUID("0000-0000-0000-0000")
	assert.Equal(t, metainfo.BookUUID, testMeta.BookUUID)

}
