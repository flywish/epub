package creator

// 定义各个文件夹 及 必要文件的名称
const (
	metaInfFolderName = "META-INF"
	contentFolderName = "EPUB"
	textFolderName    = "Text"
	cssFolderName     = "Styles"
	fontFolderName    = "Fonts"
	imageFolderName   = "Images"
	videoFolderName   = "Videos"
	audioFolderName   = "Audios"
	miscFolderName    = "Misc"

	contentFilename   = "content.opf"
	containerFilename = "container.xml"
	mimetypeFilename  = "mimetype"
)

// 定义固定文件中的内容模板
const (
	mediaTypeEpubTemplate = "application/epub+zip"
	containerFileTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
  <rootfiles>
    <rootfile full-path="%s/%s" media-type="application/oebps-package+xml" />
  </rootfiles>
</container>
`
)

// 定义其他常量
const (
	tmpDir = "tmp" // k8s需要挂载到其他目录

	dirPermissions  = 0755
	filePermissions = 0644
)
