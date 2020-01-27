## 实现目的

为了方便下载网上的小说，特意写了此项目

## API文档内容


### POST
此功能，主要用于下载小说
```bash
/post
query用到的参数
   ebhost 定义小说下载网站，类型string; 默认值为 xsbiquge.com; 可用参数为 23us.la, 999xs.com
   bookid 定义对应小说在网站中的id,类型string
   istxt 定义是否生成txt文件，类型为bool,可接受的值 0,1,true,false; 默认值为false
   ismobi 定义是否生成mobi文件，类型为bool,可接受的值 0,1,true,false; 默认值为false
  http_method: GET
返回值
 {
  "isTxt": isTxtStr, //类型 string
  "isMobi": isMobiStr, //类型 string
  "metainfo": {
        "ebhost": "xsbiquge.com",
        "bookid": "91_918743",
        "bookname": "我是谁",
        "author": "sndnvaps",
        "cover_url": "http://192.168.13.118:8090/public/我是谁-sndnvaps/cover.jpg",
        "description": "我是小说的简介信息",
        "txt_url_path": "http://192.168.13.118:8090/public/我是谁-sndnvaps/我是谁-sndnvaps.txt",
        "mobi_url_path": "http://192.168.13.118:8090/public/我是谁-sndnvaps/我是谁-sndnvaps.mobi"
            }

 }
```

测试例子
```bash
$ curl -X GET -v  "http://localhost:8080/post?ebhost=23us.la&bookid=0_062&istxt=true&ismobi=false"
```

### List
列举下载目录里面的所有文件
```bash
  /get_list
  不接受任何参数
  返回值
  {
     "files": [
         {
         "metainfo":{
            "ebhost": "xsbiquge.com",
            "bookid": "91_918743",
            "bookname": "我是谁",
            "author": "sndnvaps",
            "cover_url": "http://192.168.13.118:8090/public/我是谁-sndnvaps/cover.jpg",
            "description": "我是小说的简介信息",
            "txt_url_path": "http://192.168.13.118:8090/public/我是谁-sndnvaps/我是谁-sndnvaps.txt",
            "mobi_url_path": "http://192.168.13.118:8090/public/我是谁-sndnvaps/我是谁-sndnvaps.mobi"
            }
         },
         {
         "metainfo":{
            "ebhost": "xsbiquge.com",
            "bookid": "91_918748",
            "bookname": "我是谁1",
            "author": "sndnvaps",
            "cover_url": "http://192.168.13.118:8090/public/我是谁1-sndnvaps/cover.jpg",
            "description": "我是小说的简介信息",
            "txt_url_path": "http://192.168.13.118:8090/public/我是谁1-sndnvaps/我是谁1-sndnvaps.txt",
            "mobi_url_path": "http://192.168.13.118:8090/public/我是谁1-sndnvaps/我是谁1-sndnvaps.mobi"
            }
         }
     ]
  }
```

测试例子
```bash
$ curl -X GET -v http://localhost:8080/get_list
```

### Del
删除在服务器上面下载好的小说
```bash
 /del
 接受参数
 bookname
 返回值
 成功后的返回值
 {
    "Status": "bookname has been del"
 }
失败的返回值
 {
    "error": "bookname is not exists"
 }
```
 测试例子
```bash
$ curl -X GET -v "http://localhost:8080/del/who-am-i.txt"
```
返回结果
```json
{
   "Status": "who-am-i.txt has been delete"
}
```

### 显示服务器版本信息
此功能主要用于显示服务器的版本信息
```bash
 /stat
 不接受任何参数
 返回值
 {
  "ebookdownloader_Version": Version,
  "HashCommit": Commit,
  "SystemBuildTime": BuildTime,
  "hostinfo": {
               "host": "localhost",
               "port": "8080",
               "url_base": "http://localhost:8080"
             }

 }
```

测试例子
```bash
$ curl -X GET -v http://localhost:8080/stat
```

### Upload, 此功能已经作废
此功能是上传文件到服务器上面
```bash
  /upload
  form 传入参数 file
  http_method: POST
  返回值 
    {
        "filepath": "http://localhost:port/file/filename"
    }
```
测试使用
```bash
$ curl -X POST --form "file=@./hello.txt" http://localhost:8080/upload
```
