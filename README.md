# ebookdownloader
网文下载器

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
  .\ebookdownloader.exe --help #显示帮助信息
  ```

  ## 依赖程序 
    1. kindlegen.exe 支持windows平台
    2. kindlegenLinux 支持Linux 平台
    3. kindlegenMac 支持 Mac平台

  ## 懒人模式，直接下载编译好的程序
  
  到[这里](https://github.com/sndnvaps/ebookdownloader/releases)下载你需要的版本

  ## 更新日志

      2019.12.9 go版本添加 代理支持

      2019.12.8 go版本添加 Linux,Mac系统支持

      2019.12.6 go版本添加 进度条功能

      2019.12.5 添加go语言版本支持
      
      2019.8.22 python版本初始化  
