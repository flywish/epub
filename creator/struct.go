package creator

import (
	"sync"
)

// epub构建所需要的资源
type EpubInfo struct {
	sync.Mutex
	MetaData
}

// 元数据
type MetaData struct {
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
