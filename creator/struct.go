package creator

import (
	"sync"
)

// EpubInfo epub构建所需要的资源
type EpubInfo struct {
	sync.Mutex
	metaData
	media
	sections []section
	pkg      *pkg
	// Table of contents
	toc *toc
}

// 元数据
type metaData struct {
	title      string
	identifier string
	date       string
	publisher  string
	creator    string
	format     string
	source     string
	// type        string
	description string
	rights      string
	relation    string
	contributor string
	language    string
	subject     string
}

// 媒体资源
type media struct {
	css    map[string]string
	fonts  map[string]string
	images map[string]string
	videos map[string]string
	audios map[string]string
}

// 章节内容
type section struct {
	filename string
	xhtml    *xhtml
	children *[]section
}

type writeCounter struct {
	Total int64
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += int64(n)
	return n, nil
}
