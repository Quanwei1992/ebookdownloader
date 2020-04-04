package ebookdownloader

import (
	"errors"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/codec/gob"
	"go.etcd.io/bbolt"
)

//Boltdb 定义boltdb的接口
type Boltdb struct {
	db *storm.DB
}

//InitBoltDB 初始化boltdb数据库，根据dbname
func InitBoltDB(dbname string) (Boltdb, error) {
	db, err := storm.Open(dbname, storm.Codec(gob.Codec), storm.BoltOptions(0600, &bbolt.Options{Timeout: 1 * time.Second}))
	boltdb := Boltdb{
		db: db,
	}
	return boltdb, err
}

//Drop 删除boltdb中的bucket；不要轻易使用，除非是想删除数据库中所有数据
func (this Boltdb) Drop() error {
	err := this.db.Drop(&Meta{})
	return err
}

//Close 关闭boltdb数据库
func (this *Boltdb) Close() error {
	err := this.db.Close()
	return err
}

//Save 保存metainfo数据到boltdb数据库中
func (this *Boltdb) Save(metainfo Meta) error {
	err := this.db.Save(&metainfo)
	return err
}

//FindOneByUUID 通过uuid查询boltdb中的单一条数据
func (this Boltdb) FindOneByUUID(uuid string) (Meta, error) {
	return this.FindOneByFieldName("BookUUID", uuid)
}

//FindOneByName 通过Bookname来查询boltdb中的单一条数据
func (this Boltdb) FindOneByBookName(bookname string) (Meta, error) {
	return this.FindOneByFieldName("BookName", bookname)
}

//FindOneByAuthor 通过author来查询boltdb中的单一条数据
func (this Boltdb) FindOneByAuthor(author string) (Meta, error) {
	return this.FindOneByFieldName("Author", author)
}

//FindOneByFieldName 根据fieldname中定义的项目，进行查询，只返回一条结果
func (this Boltdb) FindOneByFieldName(fieldname string, val string) (Meta, error) {
	var metainfo Meta
	fn := fieldname
	switch fn {
	case "BookUUID":
		err := this.db.One(fn, val, &metainfo)
		return metainfo, err
	case "BookName":
		err := this.db.One(fn, val, &metainfo)
		return metainfo, err
	case "Author":
		err := this.db.One(fn, val, &metainfo)
		return metainfo, err
	default:
		return metainfo, errors.New("查询条件出错，你不应该出现在这里")
	}
}

//FindAllByAuthor 查询当前作者author在blotdb中所有的作品
func (this Boltdb) FindAllByAuthor(author string) ([]Meta, error) {
	var metainfo []Meta
	err := this.db.Find("Author", author, &metainfo)
	return metainfo, err
}

//FindAll 查询boltdb中所有的数据
func (this Boltdb) FindAll() ([]Meta, error) {
	var metainfo []Meta
	err := this.db.All(&metainfo)
	return metainfo, err
}

func (this Boltdb) UpdateEpubInfo(uuid string, epubpath string, epubmd5 string) error {
	err := this.db.Update(&Meta{BookUUID: uuid, EPUBURLPath: epubpath, EPUBMD5: epubmd5})
	return err
}

func (this Boltdb) UpdateMobiInfo(uuid string, mobipath string, mobimd5 string) error {
	err := this.db.Update(&Meta{BookUUID: uuid, MobiURLPath: mobipath, MobiMD5: mobimd5})
	return err
}

func (this Boltdb) UpdateAzw3Info(uuid string, azw3path string, azw3md5 string) error {
	err := this.db.Update(&Meta{BookUUID: uuid, AZW3URLPath: azw3path, AZW3MD5: azw3md5})
	return err
}
