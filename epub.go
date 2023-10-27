package epub

import (
	"encoding/base64"
	"epub/creator"
	"epub/reader"
	"github.com/bmaupin/go-epub"
	"io/ioutil"
	"log"
)

// 读取解析epub文件
func Read(src string, isMakeFile bool) {
	epubInfo, err := reader.OpenEpub(src, isMakeFile)
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

// github 参考包的
func WriteDemo() {
	// 参考的包的效果
	ep := epub.NewEpub("参考测试")
	ep.SetAuthor("victor")

	ep.AddCSS("tmp/epub_source/OEBPS/Styles/sgc-toc.css", "sgc-toc.css")

	ep.AddImage("tmp/epub_source/OEBPS/Images/cover.jpg", "cover.jpg")

	ep.AddSection(`<h1 class="kindle-cn-copyright-title">版权信息</h1>
  <p class="kindle-cn-copyright-text">世界文明启示录/（德）伊瓦尔·里斯纳著；吴奕俊，鲍京秀译.——天津：天津人民出版社，2020.7</p>
  <p class="kindle-cn-copyright-text">ISBN 978-7-201-15984-3</p>
  <p class="kindle-cn-copyright-text">Ⅰ.①世…　Ⅱ.①伊……②吴……③鲍…　Ⅲ.①文化史-世界-通俗读物Ⅳ.①K103-49</p>
  <p class="kindle-cn-copyright-text">中国版本图书馆CIP数据核字（2020）第090438号</p>
  <p class="kindle-cn-copyright-text">世界文明启示录</p>
  <p class="kindle-cn-copyright-text">SHIJIE WENMING QISHILU</p>
  <p class="kindle-cn-copyright-text">〔德〕伊瓦尔·里斯纳　著　吴奕俊　鲍京秀　译</p>
  <p class="kindle-cn-copyright-text">出　　版　天津人民出版社</p>
  <p class="kindle-cn-copyright-text">出 版 人　刘庆</p>
  <p class="kindle-cn-copyright-text">地　　址　天津市和平区西康路35号康岳大厦</p>
  <p class="kindle-cn-copyright-text">邮政编码　300051</p>
  <p class="kindle-cn-copyright-text">邮购电话　（022）23332469</p>
  <p class="kindle-cn-copyright-text">网　　址　http://www.tjrmcbs.com</p>
  <p class="kindle-cn-copyright-text">电子信箱　tjrmcbs@126.com</p>
  <p class="kindle-cn-copyright-text">责任编辑　郭晓雪</p>
  <p class="kindle-cn-copyright-text">特约编辑　丁兴</p>
  <p class="kindle-cn-copyright-text">装帧设计　艺琳设计</p>
  <p class="kindle-cn-copyright-text">责任校对　余艳艳</p>
  <p class="kindle-cn-copyright-text">制版印刷　天津光之彩印刷有限公司</p>
  <p class="kindle-cn-copyright-text">经　　销　新华书店</p>
  <p class="kindle-cn-copyright-text">开　　本　710毫米×1000毫米1/16</p>
  <p class="kindle-cn-copyright-text">印　　张　20.5</p>
  <p class="kindle-cn-copyright-text">字　　数　250千字</p>
  <p class="kindle-cn-copyright-text">版次印次　2020年7月第1版　2020年7月第1次印刷</p>
  <p class="kindle-cn-copyright-text">定　　价　78.00元</p>
  <p class="kindle-cn-copyright-text">版权所有，侵权必究</p>
  <p class="kindle-cn-copyright-text">图书如出现印装质量问题，请致电联系调换（022-23332469）</p>`, "第一章 测试", "section_1.xhtml", "")

	ep.Write("参考的.epub")
}

func Write() {

	// 我自己的
	epub := creator.NewEpub("测试生成")
	epub.SetCreator("魏旭晖")
	epub.Write("tmp/final.epub")

}
