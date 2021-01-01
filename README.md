# ebookdownloader
网文下载器

 [![GitHub license](https://img.shields.io/github/license/sndnvaps/ebookdownloader)](https://github.com/sndnvaps/ebookdownloader/blob/master/LICENSE)

[![Build Status](https://travis-ci.org/sndnvaps/ebookdownloader.svg?branch=master)](https://travis-ci.org/sndnvaps/ebookdownloader)[![release_version](https://img.shields.io/github/release/sndnvaps/ebookdownloader.svg)](https://github.com/sndnvaps/ebookdownloader/releases)[![Download Count](https://img.shields.io/github/downloads/sndnvaps/ebookdownloader/total.svg)](https://github.com/sndnvaps/ebookdownloader/releases)

[![Donate on PayPal](https://img.shields.io/badge/support-PayPal-blue?style=flat-square&logo=PayPal)](https://paypal.me/sndnvaps)



[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/sndnvaps/ebookdownloader/)

# ebookdl 网文下载器，go语言版本

  ## 安装方法
  ```bash
  go get github.com/sndnvaps/ebookdownloader/cli
  go get github.com/sndnvaps/ebookdownloader/gui
  go get github.com/sndnvaps/ebookdownloader/http-server
  ```
  ## 使用方法
  ```bash
  .\ebookdownloader.exe --bookid=0_642 --txt #只生成txt文本
  .\ebookdownloader.exe --bookid=0_642 --mobi #只生成mobi电子书
  .\ebookdownloader.exe --bookid=0_642 --txt --mobi #生成txt 和 mobi
  .\ebookdownloader.exe --bookid=0_642 --txt --awz3 #生成txt 和 awz3
  .\ebookdownloader.exe --proxy="http://proxyip:proxyport" --bookid=0_642 --mobi #生成mobi电子书，在下载章节的过程中使用 Proxy
  .\ebookdownloader.exe --ebhost=xsbiquge.com --bookid=0_642 --txt --mobi #使用xsbiquge.com做为下载源，生成txt 和 mobi
  .\ebookdownloader.exe --ebhost=999xs.com --bookid=0_642 --txt --mobi #使用999xs.com做为下载源，生成txt 和 mobi
   .\ebookdownloader.exe --ebhost=999xs.com --bookid=0_642 --txt --mobi --meta #使用999xs.com做为下载源，生成txt,mobi电子书，并生成meta.json文件于小说目录当中
  .\ebookdownloader.exe --ebhost=23us.la --bookid=127064 --pv #新功能，用于打印小说的分卷信息，此时不下载小说任何内容
  .\ebookdownloader.exe --bookid=0_0642 --json #生成json格式的小说数据
  .\ebookdownloader.exe conv --json=".\outputs\我是谁-sndnvaps\我是谁-sndnvaps.json" --txt --mobi #新功能，转换json格式到txt,mobi格式
  .\ebookdownloader.exe --help #显示帮助信息
  ```

  ## 依赖程序 
    1. kindlegen.exe 支持windows平台
    2. kindlegenLinux 支持Linux 平台
    3. kindlegenMac 支持 Mac平台
    4. cli/gui 两个项目，都需要在当前项目的根目录运行
    5. gui程序，需要依赖 https://github.com/akavel/rsrc ，项目来生成图标
    6. qemu-i386-static-armhf 支持在linux arm平台上运行 kindlegenLinux
    7. qemu-i386-static-arm64 支持在linux arm64平台上运行 kindlegenLinux
    8. http-server 项目依赖：
          github.com/ajvb/kala 项目，用于任务控制和管理
          kala需要与ebookdownloader_cli运行在同一个目录里面
    9. qtgui项目信赖
          goqt
          qt
          gcc
    10. chromedp 项目，用于捉取小说封面
           需要安装chrome浏览器，做支持

  ## 后端服务器 API接口
    主要目的是部署在vps上面，就可以方便随时下载小说了
   API接口文档
[ebookdownloader_http_api](http-server/ebookdownloader_http_api.md)

配置文件[ebdl_conf.ini](conf/ebdl_conf.ini)

  ## 懒人模式，直接下载编译好的程序
  
  墙里面使用gitee

  [gitee ebookdownloader release page](https://gitee.com/sndnvaps/ebookdownloader/releases "https://gitee.com/sndnvaps/ebookdownloader/releases")

墙外面使用github

  [github ebookdownloader release page](https://github.com/sndnvaps/ebookdownloader/releases "https://github.com/sndnvaps/ebookdownloader/releases")

---------------------



## ❤️ 赞助
如果你觉得EbookDownloader项目对你有帮助。你可以给我来一杯啤酒！

| PayPal打赏 | 微信打赏Pay | BTC打赏|
| :-: | :-: | :-: |
| <a href="https://paypal.me/sndnvaps"> paypal.me/sndnvaps </a>  |  <img src="qrcode/Donate_WeChatPay.jpg" width="250" height="250" alt="给sndnvaps来一杯啤酒"/>| bc1qs4zkl0fgh8csu7qn<br/>q547hymn0pxr68qrqw<br/>szst3e304grayanm4se<br/>8asf8 |


-------------

  ## 支持的小说网站

  网站名 | 网址 | 是否支持 | 备注 |
  :-: | :-: | :-: | :-: |
  笔趣阁 | https://www.xsbiquge.com/ | √ |
  笔趣阁 | https://www.biduo.cc/ | √ |
  笔趣阁 | https://www.biquwu.cc/ | √ | 网站的速度比较慢 |
  西西文学网 | http://www.xixiwx.com/ | √ |
  999小说 | https://www.999xs.com/ | √ |
  顶点小说网 | https://www.booktxt.net/ | √ |
  顶点小说网 | <s>https://www.23us.la</s> | × | 因为网站原因无法打开，暂定无法使用 |

------------

  ## 更新日志

  [CHANGELOG](./CHANGELOG "日志文件")

  -----------

  ## To Do List

     [√]  1. 添加生成封面功能
     [√]  2. 添加不同平台的接口实现
     [√]  3. 添加生成二级目录的方法(已经添加相应的实例)
     [√]  4. 添加界面版本gui
     [√]  5. 添加http-server,做为后端
     [√]  6. 添加linux arm,arm64平台支持
     [√]  7. 需要限制并发数量，因为vps性能有限 -> 目前限制的并发数量为(300+49)*2 = 698
     [√]  8. 使用boltdb记录小说数据（小说下载网站，bookid,uuid->NewV5格式，cover.jpg,mobi,azw3,txt,epub等位置及md5验证信息）
     [√]  9. 添加https://www.biduo.cc/ 小说网站支持
     [√]  10. 添加 https://www.booktxt.net/ 顶点小说网支持
     [√]  11. 添加 https://www.biquwu.cc/ 笔趣阁小说网支持
     [√]  12. 添加处理 作者，书名中包含的特殊字符串的方法
     [√]  13. 使用goqt对界面进行重写
     [ ]  14. 更新顶点小说 23us.la地址为23us.lol ->并重写其爬虫规则
