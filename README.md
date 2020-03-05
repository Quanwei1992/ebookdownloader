# ebookdownloader
网文下载器

 [![GitHub license](https://img.shields.io/github/license/sndnvaps/ebookdownloader)](https://github.com/sndnvaps/ebookdownloader/blob/master/LICENSE)

[![Build Status](https://travis-ci.org/sndnvaps/ebookdownloader.svg?branch=master)](https://travis-ci.org/sndnvaps/ebookdownloader)[![release_version](https://img.shields.io/github/release/sndnvaps/ebookdownloader.svg)](https://github.com/sndnvaps/ebookdownloader/releases)[![Download Count](https://img.shields.io/github/downloads/sndnvaps/ebookdownloader/total.svg)](https://github.com/sndnvaps/ebookdownloader/releases)



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


  ## 更新日志

  [CHANGELOG](./CHANGELOG "日志文件")

  ----------

  ## 支持的小说网站

  网站名 | 网址 | 是否支持 | 备注 |
  :-: | :-: | :-: | :-: |
  笔趣阁 | https://www.xsbiquge.com/ | √ |
  笔趣阁 | https://www.biduo.cc/ | √ |
  999小说 | https://www.999xs.com/ | √ |
  顶点小说网 | <s>https://www.23us.la</s> | × | 因为网站原因无法打开，暂定无法使用 |


  ## To Do List

     [√]  1. 添加生成封面功能
     [√]  2. 添加不同平台的接口实现
     [√]  3. 添加生成二级目录的方法(已经添加相应的实例)
     [√]  4. 添加界面版本gui
     [√]  5. 添加http-server,做为后端
     [√]  6. 添加linux arm,arm64平台支持
     [√]  7. 需要限制并发数量，因为vps性能有限 -> 目前限制的并发数量为(300+49)*2 = 698
     [ ]  8. 使用boltdb记录小说数据（小说下载网站，bookid,uuid->NewV5格式，cover.jpg,mobi,azw3,txt,epub等位置及md5验证信息）
     [√]  9. 添加https://www.biduo.cc/ 小说网站支持

