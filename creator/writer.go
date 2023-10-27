package creator

import (
	"errors"
	"os"
	"time"
)

// 创建一个epub对象, 并设置默认的信息
func NewEpub(title string) *EpubInfo {
	epub := &EpubInfo{}

	epub.SetTitle(title)
	epub.SetLanguage("zh-cn")
	epub.SetDate(time.Now().Format("2006-01-02 15:04:05"))

	return epub
}

// 解析epub内设置的元素, 进行打包
func (ep *EpubInfo) Write(outFilePath string) error {
	// 创建一个空的zip资源包
	outFile, err := os.Create(outFilePath)
	if err != nil {
		return errors.New("创建目标文件失败")
	}
	defer outFile.Close()

	err = ep.Assemble(outFilePath)
	return err
}

// 将epub的所有资源添加到zip中, 并进行zip打包
func (ep *EpubInfo) Build(zipSourceDir, outFileName string) error {

	return nil
}
