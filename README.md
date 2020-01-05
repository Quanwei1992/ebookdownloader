# ebookdownloader
网文下载器

[![Build Status](https://travis-ci.org/sndnvaps/ebookdownloader.svg?branch=master)](https://travis-ci.org/sndnvaps/ebookdownloader)

# ebookdl 网文下载器，go语言版本

  ## 安装方法
  ```bash
  go get github.com/sndnvaps/ebookdownloader
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
  .\ebookdownloader.exe --help #显示帮助信息
  ```

  ## 依赖程序 
    1. kindlegen.exe 支持windows平台
    2. kindlegenLinux 支持Linux 平台
    3. kindlegenMac 支持 Mac平台

  ## 懒人模式，直接下载编译好的程序
  
  到[这里](https://github.com/sndnvaps/ebookdownloader/releases)下载你需要的版本

  ## 更新日志

      2020.01.05 go版本 更新
                 1. 实现二级目录直接写入 tpl_*.html文件当中
                 2. 添加tpl/tpl_volume.html 用于生成目录分卷
                 3. 实现mobi格式二级目录的生成（网站捉取二级目录部分，正在努力实现）

      2020.01.04 go版本 更新
                 1. 初始化 kindle二级目录支持代码(具体功能还在实现)
                 2. 更新ebookdl_test.go: 实例化 二级目录txt支持
                 
      2020.01.03 go版本 更新
                  1. 修改生成电子书的压缩比为-c2,使生成的文件更小
                  2. 添加生成awz3格式支持(注意，--mobi,--awz3只能使用一个，不能同时使用)
                  3. 修改封面的引用方法

      2019.12.29 go版本 完成实现 999xs.com平台的小说下载接口

      2019.12.27 go版本 实现不同小说平台的interface{}接口，方便加入新的小说网站

      2019.12.25 go版本 修改小说名字排版方式为坚排

      2019.12.22 go版本添加 
                 1. 简单代码测试
                 2. 使用图片格式的封面，方便后面使用 calibre更换封面

      2019.12.9 go版本添加 代理支持

      2019.12.8 go版本添加 Linux,Mac系统支持

      2019.12.6 go版本添加 进度条功能

      2019.12.5 添加go语言版本支持
      
      2019.8.22 python版本初始化

  ## To Do List

     [√]  1.添加生成封面功能
     [√]  2. 添加不同平台的接口实现
     [ ]  3. 添加生成二级目录的方法(已经添加相应的实例)
