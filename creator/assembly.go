package creator

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gabriel-vasile/mimetype"
)

// addMedia 添加媒体资源
func (ep *EpubInfo) addMedia(source string, internalFilename string, mediaFileFormat string, mediaFolderName string, mediaMap map[string]string) (string, error) {
	// TODO :: 判断source是否真实存在(本地资源/网络资源)
	if internalFilename == "" {
		internalFilename = filepath.Base(source)
		_, ok := mediaMap[internalFilename]
		// 如果文件名太长、无效或已被使用，请尝试生成唯一的文件名
		// FIXME :: 如果文件已经被使用, 则不再重新生成
		if len(internalFilename) > 255 || !fs.ValidPath(internalFilename) || ok {
			internalFilename = fmt.Sprintf(
				mediaFileFormat,
				len(mediaMap)+1,
				strings.ToLower(filepath.Ext(source)),
			)
		}
	}

	// 判断最终的文件名是否被使用过
	if _, ok := mediaMap[internalFilename]; ok {
		return "", errors.New("文件名重复")
	}

	mediaMap[internalFilename] = source
	return path.Join(
		"..",
		mediaFolderName,
		internalFilename,
	), nil
}

// addSection 添加章节内容
func (ep *EpubInfo) addSection(parentFilename string, body string, sectionTitle string, internalFilename string, internalCSSPath string) (string, error) {
	parentIndex := -1

	// Generate a filename if one isn't provided
	if internalFilename == "" {
		index := 1
		for internalFilename == "" {
			internalFilename = fmt.Sprintf(sectionFileFormat, index)
			for item, section := range ep.sections {
				if section.filename == parentFilename {
					parentIndex = item
				}
				if section.filename == internalFilename {
					internalFilename, index = "", index+1
					if parentFilename == "" || parentIndex != -1 {
						break
					}
				}
				// Check for nested sections with the same filename to avoid duplicate entries
				if section.children != nil {
					for _, subsection := range *section.children {
						if subsection.filename == internalFilename {
							internalFilename, index = "", index+1
						}
					}
				}
			}
		}
	} else {
		for item, section := range ep.sections {
			if section.filename == parentFilename {
				parentIndex = item
			}
			if section.filename == internalFilename {
				return "", errors.New("章节名已存在")
			}
			if section.children != nil {
				for _, subsection := range *section.children {
					if subsection.filename == internalFilename {
						return "", errors.New("章节名已存在")
					}
				}
			}
		}
	}

	if parentFilename != "" && parentIndex == -1 {
		return "", errors.New("父级章节不存在")
	}

	x := newXhtml(body)
	x.setTitle(sectionTitle)
	//x.setXmlnsEpub(xmlnsEpub)

	if internalCSSPath != "" {
		x.setCSS(internalCSSPath)
	}

	s := section{
		filename: internalFilename,
		xhtml:    x,
		children: nil,
	}

	if parentIndex != -1 {
		if ep.sections[parentIndex].children == nil {
			var section []section
			ep.sections[parentIndex].children = &section
		}
		(*ep.sections[parentIndex].children) = append(*ep.sections[parentIndex].children, s)
	} else {
		ep.sections = append(ep.sections, s)
	}

	return internalFilename, nil
}

// 创建 mimetype 文件
func (ep *EpubInfo) writeMimetype(rootEpubDir string) {
	mimetypeFilePath := filepath.Join(rootEpubDir, mimetypeFilename)
	if err := os.WriteFile(mimetypeFilePath, []byte(mediaTypeEpubTemplate), filePermissions); err != nil {
		panic(fmt.Sprintf("Error writing mimetype file: %s", err))
	}
}

// 创建Epub目录
func (ep *EpubInfo) createEpubFolders(epubRootSourceDir string) {
	if err := os.MkdirAll(
		filepath.Join(
			epubRootSourceDir,
			contentFolderName,
		),
		dirPermissions); err != nil {
		// No reason this should happen if tempDir creation was successful
		panic(fmt.Sprintf("Error creating EPUB subdirectory: %s", err))
	}

	if err := os.MkdirAll(
		filepath.Join(
			epubRootSourceDir,
			contentFolderName,
			textFolderName,
		),
		dirPermissions); err != nil {
		panic(fmt.Sprintf("Error creating xhtml subdirectory: %s", err))
	}

	if err := os.MkdirAll(
		filepath.Join(
			epubRootSourceDir,
			metaInfFolderName,
		),
		dirPermissions); err != nil {
		panic(fmt.Sprintf("Error creating META-INF subdirectory: %s", err))
	}
}

// 创建固定的 container.xml 文件
func (ep *EpubInfo) writeContainerFile(rootEpubDir string) {
	containerFilePath := filepath.Join(rootEpubDir, metaInfFolderName, containerFilename)
	if err := os.WriteFile(
		containerFilePath,
		[]byte(
			fmt.Sprintf(
				containerFileTemplate,
				contentFolderName,
				contentFilename,
			),
		),
		filePermissions,
	); err != nil {
		panic(fmt.Sprintf("Error writing container file: %s", err))
	}
}

// 写入css文件
func (ep *EpubInfo) writeCSSFiles(rootEpubDir string) error {
	err := ep.writeMedia(rootEpubDir, ep.media.css, cssFolderName)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// Clean up the cover temp file if one was created
	//os.Remove(ep.cover.cssTempFile)

	return nil
}

// 写入字体文件
func (ep *EpubInfo) writeFonts(rootEpubDir string) error {
	return ep.writeMedia(rootEpubDir, ep.media.fonts, fontFolderName)
}

// 写入图片文件
func (ep *EpubInfo) writeImages(rootEpubDir string) error {
	return ep.writeMedia(rootEpubDir, ep.media.images, imageFolderName)
}

// 写入视频文件
func (ep *EpubInfo) writeVideos(rootEpubDir string) error {
	return ep.writeMedia(rootEpubDir, ep.media.videos, videoFolderName)
}

// 写入音频文件
func (ep *EpubInfo) writeAudios(rootEpubDir string) error {
	return ep.writeMedia(rootEpubDir, ep.media.audios, audioFolderName)
}

// writeMedia 通用的添加媒体资源方式
func (ep *EpubInfo) writeMedia(rootEpubDir string, mediaMap map[string]string, mediaFolderName string) error {

	log.Print(mediaMap)

	if len(mediaMap) > 0 {
		mediaFolderPath := filepath.Join(rootEpubDir, contentFolderName, mediaFolderName)
		if err := os.Mkdir(mediaFolderPath, dirPermissions); err != nil {
			return fmt.Errorf("unable to create directory: %s", err)
		}

		// TODO :: 获取文件内容, 需支持远程文件
		for mediaFilename, mediaSource := range mediaMap {
			mediaType, err := ep.fetchMedia(mediaSource, mediaFolderPath, mediaFilename)
			if err != nil {
				return err
			}

			// The cover image has a special value for the properties attribute
			mediaProperties := ""
			//if mediaFilename == ep.cover.imageFilename {
			//	mediaProperties = coverImageProperties
			//}
			//
			//// Add the file to the OPF manifest
			ep.pkg.addToManifest(fixXMLId(mediaFilename), filepath.Join(mediaFolderName, mediaFilename), mediaType, mediaProperties)
		}
	}
	return nil
}

// 获取媒体资源并写入到对应的文件中
func (ep *EpubInfo) fetchMedia(mediaSource, mediaFolderPath, mediaFilename string) (mediaType string, err error) {
	mediaFilePath := filepath.Join(
		mediaFolderPath,
		mediaFilename,
	)
	// failfast, create the output file handler at the begining, if we cannot write the file, bail out
	w, err := os.Create(mediaFilePath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败 %s: %s", mediaFilePath, err)
	}
	defer w.Close()
	var source io.ReadCloser
	source, err = os.Open(mediaSource)

	if source == nil {
		return "", fmt.Errorf("打开文件失败 %s: %s", mediaFilePath, err)
	}
	defer source.Close()

	_, err = io.Copy(w, source)
	if err != nil {
		return "", fmt.Errorf("写入文件失败 %s: %s", mediaFilePath, err)
	}

	// Detect the mediaType
	r, err := os.Open(mediaFilePath)
	if err != nil {
		return "", err
	}
	defer r.Close()
	mime, err := mimetype.DetectReader(r)
	if err != nil {
		panic(err)
	}

	// Is it CSS?
	mtype := mime.String()
	if mime.Is("text/plain") {
		if filepath.Ext(mediaSource) == ".css" || filepath.Ext(mediaFilename) == ".css" {
			mtype = "text/css"
		}
	}
	return mtype, nil
}

// 生成package 中的id
func fixXMLId(id string) string {
	if len(id) == 0 {
		panic("No id given")
	}
	fixedId := []rune{}
	for i := 0; len(id) > 0; i++ {
		r, size := utf8.DecodeRuneInString(id)
		if i == 0 {
			// The new id should be prefixed with 'id' if an invalid
			// starting character is found
			// this is not 100% accurate, but a better check than no check
			if unicode.IsNumber(r) || unicode.IsPunct(r) || unicode.IsSymbol(r) {
				fixedId = append(fixedId, []rune("id")...)
			}
		}
		if !unicode.IsSpace(r) && r != ':' {
			fixedId = append(fixedId, r)
		}
		id = id[size:]
	}
	return string(fixedId)
}

func (ep *EpubInfo) writeSections(rootEpubDir string) {
	var index int

	if len(ep.sections) > 0 {
		// If a cover was set, add it to the package spine first so it shows up
		// first in the reading order
		//if ep.cover.xhtmlFilename != "" {
		//	ep.pkg.addToSpine(ep.cover.xhtmlFilename)
		//}

		for _, section := range ep.sections {
			// Set the title of the cover page XHTML to the title of the EPUB
			//if section.filename == ep.cover.xhtmlFilename {
			//	section.xhtml.setTitle(ep.Title())
			//}

			sectionFilePath := filepath.Join(rootEpubDir, contentFolderName, textFolderName, section.filename)
			section.xhtml.write(sectionFilePath)
			relativePath := filepath.Join(textFolderName, section.filename)

			// The cover page should have already been added to the spine first
			//if section.filename != ep.cover.xhtmlFilename {
			ep.pkg.addToSpine(section.filename)
			//}
			ep.pkg.addToManifest(section.filename, relativePath, mediaTypeXhtml, "")

			// Don't add pages without titles or the cover to the TOC
			if section.xhtml.Title() != "" {
				ep.toc.addSection(index, section.xhtml.Title(), relativePath)

				// Add subsections
				if section.children != nil {
					for _, child := range *section.children {
						index += 1
						relativeSubPath := filepath.Join(textFolderName, child.filename)
						ep.toc.addSubSection(relativePath, index, child.xhtml.Title(), relativeSubPath)

						subSectionFilePath := filepath.Join(rootEpubDir, contentFolderName, textFolderName, child.filename)
						child.xhtml.write(subSectionFilePath)

						// Add subsection to spine
						ep.pkg.addToSpine(child.filename)
						ep.pkg.addToManifest(child.filename, relativeSubPath, mediaTypeXhtml, "")
					}
				}
			}

			index += 1
		}
	}
}

func (ep *EpubInfo) writeNav(rootEpubDir string) {
	ep.pkg.addToManifest(tocNavItemID, filepath.Join(textFolderName, tocNavFilename), mediaTypeXhtml, tocNavItemProperties)
	// ep.pkg.addToManifest(tocNcxItemID, tocNcxFilename, mediaTypeNcx, "")

	ep.toc.write(rootEpubDir)
}

// 写入
func (ep *EpubInfo) writePackageFile(rootEpubDir string) {
	ep.pkg.write(rootEpubDir)
}
