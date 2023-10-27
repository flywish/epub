package reader

import (
	"archive/zip"
	"os"
)

// epub 格式中, 这个文件是固定的
const containerPath = "META-INF/container.xml"

// epub 文件查找路径: Container -> package -> metadata/manifest/spine/guide
type EpubInfo struct {
	Container
	files map[string]*zip.File
}

type EpubCloser struct {
	EpubInfo
	f *os.File
}

type Container struct {
	RootFiles []*RootFile `xml:"rootfiles>rootfile"`
}

type RootFile struct {
	FullPath string `xml:"full-path,attr"`
	Package
}

type Package struct {
	Metadata
	Manifest
	Spine
}

type Metadata struct {
	Title       string `xml:"metadata>title"`
	Language    string `xml:"metadata>language"`
	Identifier  string `xml:"metadata>idenifier"`
	Creator     string `xml:"metadata>creator"`
	Contributor string `xml:"metadata>contributor"`
	Publisher   string `xml:"metadata>publisher"`
	Subject     string `xml:"metadata>subject"`
	Description string `xml:"metadata>description"`
	Event       []struct {
		Name string `xml:"event,attr"`
		Date string `xml:",innerxml"`
	} `xml:"metadata>date"`
	Type     string `xml:"metadata>type"`
	Format   string `xml:"metadata>format"`
	Source   string `xml:"metadata>source"`
	Relation string `xml:"metadata>relation"`
	Coverage string `xml:"metadata>coverage"`
	Rights   string `xml:"metadata>rights"`
}

type Manifest struct {
	Items []Item `xml:"manifest>item"`
}

type Item struct {
	Id        string `xml:"id,attr"`
	Href      string `xml:"href,attr"`
	MediaType string `xml:"media-type,attr"`
	f         *zip.File
}

type Spine struct {
	Itemrefs []Itemref `xml:"spine>itemref"`
}

type Itemref struct {
	Idref string `xml:"idref,attr"`
	*Item
}
