package creator

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

// 书名
func (ep *EpubInfo) SetTitle(title string) {
	ep.Lock()
	defer ep.Unlock()
	ep.MetaData.title = title
}

// Identifier
func (ep *EpubInfo) SetIdentifier(identifier string) {
	ep.Lock()
	defer ep.Unlock()
	ep.MetaData.identifier = identifier
}

// 设置创建日期
func (ep *EpubInfo) SetDate(date string) {
	ep.Lock()
	defer ep.Unlock()
	ep.MetaData.date = date
}

// 出版商
func (ep *EpubInfo) SetPublisher(publisher string) {
	ep.Lock()
	defer ep.Unlock()
	ep.MetaData.publisher = publisher
}

// 作者
func (ep *EpubInfo) SetCreator(creator string) {
	ep.Lock()
	defer ep.Unlock()
	ep.MetaData.creator = creator
}

// 格式
func (ep *EpubInfo) SetFormat(format string) {
	ep.Lock()
	defer ep.Unlock()
	ep.MetaData.format = format
}

// 来源
func (ep *EpubInfo) SetSource(source string) {
	ep.Lock()
	defer ep.Unlock()
	ep.MetaData.source = source
}

// 类型`
//func (ep *EpubInfo) SetType(t string) {
//	ep.Lock()
//	defer ep.Unlock()
//	ep.MetaData.Type = t
//}

// 设置描述
func (ep *EpubInfo) SetDescription(description string) {
	ep.Lock()
	defer ep.Unlock()
	ep.MetaData.description = description
}

// 权利
func (ep *EpubInfo) SetRights(rights string) {
	ep.Lock()
	defer ep.Unlock()
	ep.MetaData.rights = rights
}

// 涉及
func (ep *EpubInfo) SetRelation(relation string) {
	ep.Lock()
	defer ep.Unlock()
	ep.MetaData.relation = relation
}

// 投稿者
func (ep *EpubInfo) SetContributor(contributor string) {
	ep.Lock()
	defer ep.Unlock()
	ep.MetaData.contributor = contributor
}

// 语言
func (ep *EpubInfo) SetLanguage(language string) {
	ep.Lock()
	defer ep.Unlock()
	ep.MetaData.language = language
}

// 主题
func (ep *EpubInfo) SetSubject(subject string) {
	ep.Lock()
	defer ep.Unlock()
	ep.MetaData.subject = subject
}

// 根据设置的所有数据, 构建zip包内的资源
func (ep *EpubInfo) Assemble(outFilePath string) error {
	ep.Lock()
	defer ep.Unlock()

	// 创建一个随机的目录, 用于暂存epub资源
	//uuidDirName := uuid.Must(uuid.NewUUID()).String()
	uuidDirName := "70c635c0-7496-11ee-8f4d-1aacdefa13f1"
	epubRootTempDir := path.Join(tmpDir, uuidDirName)

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

	writeMimetype(epubRootTempDir)
	createEpubFolders(epubRootTempDir)

	// Must be called after:
	// createEpubFolders()
	writeContainerFile(epubRootTempDir)

	// Must be called after:
	// createEpubFolders()
	//err = ep.writeCSSFiles(epubRootTempDir)
	//if err != nil {
	//	return err
	//}

	// Must be called after:
	// createEpubFolders()
	//err = ep.writeFonts(epubRootTempDir)
	//if err != nil {
	//	return err
	//}

	// Must be called after:
	// createEpubFolders()
	//err = ep.writeImages(epubRootTempDir)
	//if err != nil {
	//	return err
	//}

	// Must be called after:
	// createEpubFolders()
	//err = ep.writeVideos(epubRootTempDir)
	//if err != nil {
	//	return err
	//}

	// Must be called after:
	// createEpubFolders()
	//err = ep.writeAudios(epubRootTempDir)
	//if err != nil {
	//	return err
	//}

	// Must be called after:
	// createEpubFolders()
	//ep.writeSections(epubRootTempDir)

	// Must be called after:
	// createEpubFolders()
	// writeSections()
	//ep.writeToc(epubRootTempDir)

	// Must be called after:
	// createEpubFolders()
	// writeCSSFiles()
	// writeImages()
	// writeVideos()
	// writeAudios()
	// writeSections()
	// writeToc()
	//ep.writePackageFile(epubRootTempDir)

	// Must be called last
	err = ep.Build(epubRootTempDir, outFilePath)

	return nil
}

// memetype 文件
func writeMimetype(rootEpubDir string) {
	mimetypeFilePath := filepath.Join(rootEpubDir, mimetypeFilename)
	if err := os.WriteFile(mimetypeFilePath, []byte(mediaTypeEpubTemplate), filePermissions); err != nil {
		panic(fmt.Sprintf("Error writing mimetype file: %s", err))
	}
}

// 创建Epub目录
func createEpubFolders(epubRootSourceDir string) {
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

func writeContainerFile(rootEpubDir string) {
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

func (ep *EpubInfo) writeCSSFiles() {

}
