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
	mediaTypeCSS      = "text/css"

	// 生成的默认文件名规则
	cssFileFormat     = "css_%04d%s"
	fontFileFormat    = "font_%04d%s"
	imageFileFormat   = "image_%04d%s"
	videoFileFormat   = "video_%04d%s"
	audioFileFormat   = "audio_%04d%s"
	sectionFileFormat = "section_%04d.xhtml"

	mediaTypeXhtml = "application/xhtml+xml"
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
	xhtmlTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
  <head>
    <title></title>
  </head>
  <body></body>
</html>
`
)

// 定义其他常量
const (
	tmpDir = "tmp" // k8s需要挂载到其他目录

	dirPermissions  = 0755
	filePermissions = 0644
)
