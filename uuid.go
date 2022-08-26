package ebookdownloader

import (
	uuid "github.com/satori/go.uuid"
)

// GenerateUUID 根据小说的作者名和小说名 生成uuid码，使用UUID_V5格式
func (this *BookInfo) GenerateUUID() {
	u := uuid.NewV5(uuid.NamespaceOID, this.Name+"-"+this.Author)
	uuidStr := u.String()
	this.BookUUID = uuidStr
}

// UUID 返回小说对应的uuid信息
func (this BookInfo) UUID() string {
	return this.BookUUID
}
