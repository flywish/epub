package reader

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

func OpenEpub(src string, isMakeFile bool) (*EpubCloser, error) {
	originFile, err := os.Open(src)
	if err != nil {
		return nil, err
	}

	originFileInfo, err := originFile.Stat()
	if err != nil {
		originFile.Close()
		return nil, err
	}

	zipSource, err := zip.NewReader(originFile, originFileInfo.Size())
	if err != nil {
		return nil, err
	}

	// 初始化 epubInfo
	epubInfo := new(EpubCloser)
	epubInfo.f = originFile

	err = epubInfo.analyze(zipSource, isMakeFile)
	if err != nil {
		return nil, err
	}

	return epubInfo, nil
}

// 分析zip, 如果需要创建本地文件，则解压缩到本地
func (ep *EpubInfo) analyze(zipSource *zip.Reader, isMakeFile bool) error {
	ep.files = make(map[string]*zip.File)

	for _, file := range zipSource.File {
		ep.files[file.Name] = file

		// 创建本地文件
		path := filepath.Join("./tmp/epub_source", file.Name)
		rc, err := file.Open()
		if err != nil {
			return err
		}
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, 0777)
		} else {
			os.MkdirAll(filepath.Dir(path), 0777)
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			_, err = io.Copy(f, rc)
			f.Close()
		}
		rc.Close()
	}

	err := ep.setContainer()
	if err != nil {
		return err
	}

	err = ep.setPackages()
	if err != nil {
		return err
	}

	err = ep.setItems()
	if err != nil {
		return err
	}

	return nil
}

// 解析rootfiles -> rootfile
func (ep *EpubInfo) setContainer() error {
	containerFile, err := ep.files[containerPath].Open()
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, containerFile)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(buf.Bytes(), &ep.Container)
	if err != nil {
		return err
	}

	if len(ep.Container.RootFiles) < 1 {
		return err
	}
	return nil
}

// 解析 package
func (ep *EpubInfo) setPackages() error {
	for _, rootFile := range ep.Container.RootFiles {
		if ep.files[rootFile.FullPath] == nil {
			return errors.New("未找到full_path对应的文件")
		}

		packageFile, err := ep.files[rootFile.FullPath].Open()
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		_, err = io.Copy(&buf, packageFile)
		if err != nil {
			return err
		}

		err = xml.Unmarshal(buf.Bytes(), &rootFile.Package)
		if err != nil {
			return err
		}
	}
	return nil
}

// 解析
func (ep *EpubInfo) setItems() error {
	itemrefCount := 0
	for _, rootFile := range ep.Container.RootFiles {
		itemMap := make(map[string]*Item)

		// 获取 manifest 内容
		for i := range rootFile.Manifest.Items {
			item := &rootFile.Manifest.Items[i]
			itemMap[item.Id] = item

			// 根据item的href获取资源的绝对路径(相对解压之后的路径而言)
			absPath := path.Join(path.Dir(rootFile.FullPath), item.Href)
			item.f = ep.files[absPath]
		}

		// 获取 spine 内容
		for i := range rootFile.Spine.Itemrefs {
			itemref := &rootFile.Spine.Itemrefs[i]
			itemref.Item = itemMap[itemref.Idref]
			if itemref.Item == nil {
				return errors.New("对应资源缺失")
			}
		}

		itemrefCount += len(rootFile.Spine.Itemrefs)
	}

	if itemrefCount < 1 {
		return errors.New("未找到合法节点")
	}
	return nil
}

// ItemOpen 打开item对应的文件
func (item *Item) ItemOpen() (io.ReadCloser, error) {
	if item.f == nil {
		return nil, errors.New("文件不存在")
	}

	return item.f.Open()
}

// ItemInfo 获取item信息
func (item *Item) ItemInfo() (fs.FileInfo, error) {
	if item.f == nil {
		return nil, errors.New("文件不存在")
	}
	return item.f.FileInfo(), nil
}

// epub close
func (ep *EpubCloser) CloseEpub() {
	ep.f.Close()
}
