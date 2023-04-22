package ebook

import (
	edl "github.com/sndnvaps/ebookdownloader"
)

// AsycChapter 同步下载章节的content内容
func AsycChapter(ResultChan chan chan edl.Chapter, chapter chan edl.Chapter) {
	for {
		c := <-ResultChan
		tmp := <-c
		chapter <- tmp
	}

}
