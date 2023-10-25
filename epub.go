package epub

import (
	"encoding/base64"
	"epub/reader"
	"io/ioutil"
	"log"
)

// 读取解析epub文件
func Read(src string) {
	epubInfo, err := reader.OpenEpub(src)
	if err != nil {
		panic(err)
	}

	defer epubInfo.CloseEpub()

	book := epubInfo.EpubInfo.RootFiles[0]
	for _, item := range book.Manifest.Items {
		// TODO :: 可以将image,css,video,audio等资源上传到oss
		if item.MediaType == "image/jpeg" {
			// 处理图片
			content, err := item.ItemOpen()
			if err != nil {
				println("文件打开失败")
			}
			b, _ := ioutil.ReadAll(content)
			baseimg := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(b)
			log.Println(baseimg)
			break
		}

		//item.ItemInfo()
		//break

	}
}

func Write() {

}
