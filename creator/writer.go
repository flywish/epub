package creator

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

// NewEpub 创建一个epub对象, 并设置默认的信息
func NewEpub(title string) *EpubInfo {
	epub := &EpubInfo{}

	epub.SetTitle(title)
	epub.SetLanguage("zh-cn")
	epub.SetDate(time.Now().Format("2006-01-02 15:04:05"))

	epub.media.css = make(map[string]string)
	epub.media.fonts = make(map[string]string)
	epub.media.images = make(map[string]string)
	epub.media.videos = make(map[string]string)
	epub.media.audios = make(map[string]string)

	epub.pkg = newPackage()
	epub.toc = newToc()

	return epub
}

// SetTitle 书名
func (ep *EpubInfo) SetTitle(title string) {
	ep.Lock()
	defer ep.Unlock()
	ep.metaData.title = title
}

// SetIdentifier Identifier
func (ep *EpubInfo) SetIdentifier(identifier string) {
	ep.Lock()
	defer ep.Unlock()
	ep.metaData.identifier = identifier
}

// SetDate 设置创建日期
func (ep *EpubInfo) SetDate(date string) {
	ep.Lock()
	defer ep.Unlock()
	ep.metaData.date = date
}

// SetPublisher 出版商
func (ep *EpubInfo) SetPublisher(publisher string) {
	ep.Lock()
	defer ep.Unlock()
	ep.metaData.publisher = publisher
}

// SetCreator 作者
func (ep *EpubInfo) SetCreator(creator string) {
	ep.Lock()
	defer ep.Unlock()
	ep.metaData.creator = creator
}

// SetFormat 格式
func (ep *EpubInfo) SetFormat(format string) {
	ep.Lock()
	defer ep.Unlock()
	ep.metaData.format = format
}

// SetSource 来源
func (ep *EpubInfo) SetSource(source string) {
	ep.Lock()
	defer ep.Unlock()
	ep.metaData.source = source
}

// SetType 类型
//func (ep *EpubInfo) SetType(t string) {
//	ep.Lock()
//	defer ep.Unlock()
//	ep.metaData.Type = t
//}

// SetDescription 设置描述
func (ep *EpubInfo) SetDescription(description string) {
	ep.Lock()
	defer ep.Unlock()
	ep.metaData.description = description
}

// SetRights 权利
func (ep *EpubInfo) SetRights(rights string) {
	ep.Lock()
	defer ep.Unlock()
	ep.metaData.rights = rights
}

// SetRelation 涉及
func (ep *EpubInfo) SetRelation(relation string) {
	ep.Lock()
	defer ep.Unlock()
	ep.metaData.relation = relation
}

// SetContributor 投稿者
func (ep *EpubInfo) SetContributor(contributor string) {
	ep.Lock()
	defer ep.Unlock()
	ep.metaData.contributor = contributor
}

// SetLanguage 语言
func (ep *EpubInfo) SetLanguage(language string) {
	ep.Lock()
	defer ep.Unlock()
	ep.metaData.language = language
}

// SetSubject 主题
func (ep *EpubInfo) SetSubject(subject string) {
	ep.Lock()
	defer ep.Unlock()
	ep.metaData.subject = subject
}

// AddCSS 添加样式文件
func (ep *EpubInfo) AddCSS(source string, internalFilename string) (string, error) {
	ep.Lock()
	defer ep.Unlock()
	return ep.addMedia(source, internalFilename, cssFileFormat, cssFolderName, ep.media.css)
}

// AddFont 字体
func (ep *EpubInfo) AddFont(source string, internalFilename string) (string, error) {
	ep.Lock()
	defer ep.Unlock()
	return ep.addMedia(source, internalFilename, fontFileFormat, fontFolderName, ep.media.fonts)
}

// AddImage 图片
func (ep *EpubInfo) AddImage(source string, imageFilename string) (string, error) {
	ep.Lock()
	defer ep.Unlock()
	return ep.addMedia(source, imageFilename, imageFileFormat, imageFolderName, ep.media.images)
}

// AddVideo 视频
func (ep *EpubInfo) AddVideo(source string, videoFilename string) (string, error) {
	ep.Lock()
	defer ep.Unlock()
	return ep.addMedia(source, videoFilename, videoFileFormat, videoFolderName, ep.media.videos)
}

// AddAudio 音频
func (ep *EpubInfo) AddAudio(source string, audioFilename string) (string, error) {
	ep.Lock()
	defer ep.Unlock()
	return ep.addMedia(source, audioFilename, audioFileFormat, audioFolderName, ep.media.audios)
}

// AddSection 章节
func (ep *EpubInfo) AddSection(body string, sectionTitle string, internalFilename string, internalCSSPath string) (string, error) {
	ep.Lock()
	defer ep.Unlock()
	return ep.addSection("", body, sectionTitle, internalFilename, internalCSSPath)
}

// AddSubSection 子章节
func (ep *EpubInfo) AddSubSection(parentFilename string, body string, sectionTitle string, internalFilename string, internalCSSPath string) (string, error) {
	ep.Lock()
	defer ep.Unlock()
	return ep.addSection(parentFilename, body, sectionTitle, internalFilename, internalCSSPath)
}

// Write 解析epub内设置的元素, 进行打包
func (ep *EpubInfo) Write(outFilePath string) error {
	// 创建一个空的zip资源包
	outFile, err := os.Create(outFilePath)
	if err != nil {
		return errors.New("创建目标文件失败")
	}
	defer outFile.Close()

	err = ep.Assemble(outFile)
	return err
}

// Assemble 根据设置的所有数据, 构建zip包内的资源
func (ep *EpubInfo) Assemble(outFile io.Writer) error {
	ep.Lock()
	defer ep.Unlock()

	// 创建一个随机的目录, 用于暂存epub资源,
	// TODO :: 建议存放在 os.TempDir() 下面
	uuidDirName := uuid.Must(uuid.NewUUID()).String()
	epubRootTempDir := path.Join(tmpDir, uuidDirName)
	absEpubRootTempDir, _ := filepath.Abs(epubRootTempDir)
	println("临时文件夹地址: " + absEpubRootTempDir)

	err := os.MkdirAll(epubRootTempDir, 0777)
	if err != nil {
		return errors.New(fmt.Sprint("创建临时文件夹失败:%s", err))
	}
	// FIXME :: 记得删除临时文件
	//defer func() {
	//	if err := os.RemoveAll(epubRootTempDir); err != nil {
	//		panic(fmt.Sprintf("删除临时文件夹失败: %s", err))
	//	}
	//}()

	ep.writeMimetype(epubRootTempDir)
	ep.createEpubFolders(epubRootTempDir)

	// 必须在以下方法之后调用
	// createEpubFolders()
	ep.writeContainerFile(epubRootTempDir)

	// 必须在以下方法之后调用
	// createEpubFolders()
	err = ep.writeCSSFiles(epubRootTempDir)
	if err != nil {
		return err
	}

	// 必须在以下方法之后调用
	// createEpubFolders()
	err = ep.writeFonts(epubRootTempDir)
	if err != nil {
		log.Printf("err: %s", err.Error())
		return err
	}

	// 必须在以下方法之后调用
	// createEpubFolders()
	err = ep.writeImages(epubRootTempDir)
	if err != nil {
		return err
	}

	// 必须在以下方法之后调用
	// createEpubFolders()
	err = ep.writeVideos(epubRootTempDir)
	if err != nil {
		return err
	}

	// 必须在以下方法之后调用
	// createEpubFolders()
	err = ep.writeAudios(epubRootTempDir)
	if err != nil {
		return err
	}

	// 必须在以下方法之后调用
	// createEpubFolders()
	ep.writeSections(epubRootTempDir)

	// 必须在以下方法之后调用
	// createEpubFolders()
	// writeSections()
	// epub3 没有 toc.ncx 了, 只有nav.html
	ep.writeNav(epubRootTempDir)

	// 必须在以下方法之后调用
	// createEpubFolders()
	// writeCSSFiles()
	// writeImages()
	// writeVideos()
	// writeAudios()
	// writeSections()
	// writeToc()
	ep.writePackageFile(epubRootTempDir)

	// Must be called last
	_, err = ep.Build(epubRootTempDir, outFile)

	return nil
}

// Build 将epub的所有资源添加到zip中, 并进行zip打包
func (ep *EpubInfo) Build(zipSourceDir string, outFile io.Writer) (int64, error) {
	counter := &writeCounter{}
	teeWriter := io.MultiWriter(counter, outFile)

	z := zip.NewWriter(teeWriter)

	skipMimetypeFile := false

	addFileToZip := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Get the path of the file relative to the folder we're zipping
		relativePath, err := filepath.Rel(zipSourceDir, path)

		if err != nil {
			// tempDir and path are both internal, so we shouldn't get here
			return err
		}
		relativePath = filepath.ToSlash(relativePath)

		// Only include regular files, not directories
		info, err := d.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}

		var w io.Writer
		if filepath.FromSlash(path) == filepath.Join(zipSourceDir, mimetypeFilename) {
			// Skip the mimetype file if it's already been written
			if skipMimetypeFile == true {
				return nil
			}
			// The mimetype file must be uncompressed according to the EPUB spec
			w, err = z.CreateHeader(&zip.FileHeader{
				Name:   relativePath,
				Method: zip.Store,
			})
		} else {
			w, err = z.Create(relativePath)
		}
		if err != nil {
			return fmt.Errorf("error creating zip writer: %w", err)
		}

		r, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("error opening file %v being added to EPUB: %w", path, err)
		}
		defer func() {
			if err := r.Close(); err != nil {
				panic(err)
			}
		}()

		_, err = io.Copy(w, r)
		if err != nil {
			return fmt.Errorf("error copying contents of file being added EPUB: %w", err)
		}
		return nil
	}

	//开始添加文件
	minetypeFilepath := filepath.Join(zipSourceDir, mimetypeFilename)
	//minetypeInfo, err := fs.Stat(os.DirFS(zipSourceDir), minetypeFilepath)
	minetypeInfo, err := os.Stat(minetypeFilepath)
	if err != nil {
		if err := z.Close(); err != nil {
			panic(err)
		}
		return counter.Total, fmt.Errorf("unable to get FileInfo for mimetype file: %w", err)
	}

	err = addFileToZip(minetypeFilepath, fileInfoToDirEntry(minetypeInfo), nil)
	if err != nil {
		if err := z.Close(); err != nil {
			panic(err)
		}
		return counter.Total, fmt.Errorf("unable to get FileInfo for mimetype file: %w", err)
	}

	skipMimetypeFile = true

	// err = fs.WalkDir(os.DirFS(zipSourceDir), ".", addFileToZip)
	err = filepath.WalkDir(zipSourceDir, addFileToZip)

	if err != nil {
		if err := z.Close(); err != nil {
			panic(err)
		}
		return counter.Total, fmt.Errorf("unable to add file to EPUB: %w", err)
	}

	println("所有文件打包完毕")

	err = z.Close()
	return counter.Total, err
}
