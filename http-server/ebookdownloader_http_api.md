## 实现目的

为了方便下载网上的小说，特意写了此项目


### api v2 添加权限管理
需要验证 username,password
也可以使用token

## API文档内容

### Login 登陆，生成 token
此功能主要用于进行登陆验证，生成token
```bash
/login
username 类型 string; 默认用户为admin; 传入方式 --form
password 类型 string;默认密码 admin; 传入方式 --form
http_method: POST
返回值
{
   "code": 200,
   "expire": "2020-01-28T19:23:09+08:00",
   "token":"eyJhbGciOiJIUzI1NiIsIn
R5cCI6IkpXVCJ9.eyJleHAiOjE1ODAyMTA1ODksImlkIjoiYWRtaW4iLCJvcmlnX2lhdCI6MTU4MDIwN
jk4OX0.XNPrk0LKcMJlJqf0Opx9JYh_kaKL_STT5p7J9_mkc0Y"
}
```

测试例子
```bash
$curl -X POST -form username=admin --form password=admin http://localhost:8080/login
```

### Logout 退出登陆，并删除 token
此功能用于退出当前用户的登陆，并删除相应的token
```bash
/logout
没有参数
http_method: GET
返回值
{
 "code": 200
}
```

测试例子
```bash
$curl -X GET http://localhost:8080/logout
```

### 更新 token
此功能，用于更新已经失效的 TOKEN

```bash
/auth/refresh_token
没有参数
http_method: GET
但需要定义传入的http header
传入header:
  Authorization:Bearer TOKEN
  Content-Type: application/json
返回值
{
   "code": 200,
   "expire": "2020-01-28T19:23:09+08:00",
   "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODAyMTA1ODksImlkIjoiYWRtaW4iLCJvcmlnX2lhdCI6MTU4MDIwNjk4OX0.XNPrk0LKcMJlJqf0Opx9JYh_kaKL_STT5p7J9_mkc0Y"
}
```

测试例子

```bash
$curl  -H "Authorization:Bearer token_string" -H "Content-Type: application/json" -X GET localhost:8000/auth/refresh_token 
```

### POST
此功能，主要用于下载小说
```bash
/auth/post
query用到的参数
   ebhost 定义小说下载网站，类型string; 默认值为 xsbiquge.com; 可用参数为 23us.la, 999xs.com
   bookid 定义对应小说在网站中的id,类型string
   istxt 定义是否生成txt文件，类型为bool,可接受的值 0,1,true,false; 默认值为false
   ismobi 定义是否生成mobi文件，类型为bool,可接受的值 0,1,true,false; 默认值为false
  http_method: GET
传入header:
  Authorization:Bearer TOKEN
  Content-Type: application/json
或者 cookie: 需要传入 token
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
$ curl -H "Authorization:Bearer token_string" -H "Content-Type: application/json" -X GET -v  "http://localhost:8080/auth/post?ebhost=23us.la&bookid=0_062&istxt=true&ismobi=false"
```

### List
列举下载目录里面的所有文件
```bash
  /auth/get_list
  不接受任何参数
   http_method: GET
   传入header:
  Authorization:Bearer TOKEN
  Content-Type: application/json
 或者 cookie: 需要传入 token
 ```

 ```json
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
$ curl  -H "Authorization:Bearer token_string" -H "Content-Type: application/json" -X GET -v http://localhost:8080/auth/get_list
```

### Del
删除在服务器上面下载好的小说

```bash
 /auth/del
 接受参数
 ebpath //定义小说的路径；类型 string; 格式： 小说名-作者
 bookname //定义小说名，支持格式有.txt,.mobi,azw3,.jpg,.json；类型string; 可接受特定命令del，用于删除小说对应的目录
 传入header:
  Authorization:Bearer TOKEN
  Content-Type: application/json
 或者 cookie: 需要传入 token
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
$ curl  -H "Authorization:Bearer xxxxxxxxx" -H "Content-Type: application/json" -X GET -v "http://localhost:8080/auth/del/我是谁-sndnvaps/我是谁-sndnvaps.txt"
$ curl  -H "Authorization:Bearer xxxxxxxxx" -H "Content-Type: application/json" -X GET -v "http://localhost:8080/auth/del/我是谁-sndnvaps/del"
```

返回结果
```json
{
   "Status": "我是谁-sndnvaps.txt has been delete"
}
{
      "Status": "我是谁-sndnvaps has been remove"
}
```

### 显示服务器版本信息
此功能主要用于显示服务器的版本信息

```bash
 /auth/stat
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
$ curl  -H "Authorization:Bearer xxxxxxxxx" -H "Content-Type: application/json" -X GET -v http://localhost:8080/auth/stat
```

## 任务调度功能

## create job
此功能用于创建定时下载任务（下载任务会在当前时间10分钟后执行)
需要配置/conf/ebd_conf.ini文件来定义ebookdownloader_cli的位置，要于kala程序运行目录一致
```bash
  /api/v1/job
  传入query参数可以是如下
    ebhost string类型 定义下载源
    bookid string类型 定义小说的bookid
    istxt string类型 true,false
    ismobi string类型 true,false
    meta string类型 true,false；默认为true
 传入header:
    Authorization:Bearer TOKEN
    Content-Type: application/json
  http_method: POST
  返回值
    成功创建的结果
  {
    	"code":        200,
		"msg":         "创建下载任务完成",
		"JobID":       jobID,
		"JobName":     jobName,
		"JobCMD":      cmd,
		"JobSchedule": schedule,
  }
   失败返回结果
     {
        	"code":   201,
			"msg":    "创建下载下载失败",
			"errMsg": err.Error(),
     }
```

测试例子
```bash
$curl  -H "Authorization:Bearer xxxxxxxxx" -H "Content-Type: application/json" -X POST http://192.168.13.107:8090/api/v1/job?bookid=90_90583&istxt=true&ismobi=true
```

### 查询任务信息
查询任务信息

```bash
  /api/v1/job/*id
  可传入参数 id
 传入header:
    Authorization:Bearer TOKEN
    Content-Type: application/json
  http_method: GET
  返回结果
  执行 /api/vi/job/ 后的结果
```

```json
{
    "id": ""
}
{
    "jobinfos": [
        {
            "id": "5cd1fd60-e6a8-42a3-63e9-c4b4474f26b4",
            "job": {
                "name": "Downloader ebook 垂钓之神-会狼叫的猪",
                "id": "5cd1fd60-e6a8-42a3-63e9-c4b4474f26b4",
                "command": "F:/work/ebookdownloader/ebookdownloader_cli.exe --ebhost=xsbiquge.com --bookid=90_90583 --txt --mobi --meta",
                "owner": "",
                "disabled": false,
                "dependent_jobs": null,
                "parent_jobs": null,
                "on_failure_job": "",
                "schedule": "R0/2020-02-14T19:01:43+08:00/",
                "retries": 0,
                "epsilon": "",
                "next_run_at": "2020-02-14T19:01:43+08:00",
                "resume_at_next_scheduled_time": false,
                "metadata": {
                    "success_count": 1,
                    "last_success": "2020-02-14T19:02:45.4369279+08:00",
                    "error_count": 0,
                    "last_error": "0001-01-01T00:00:00Z",
                    "last_attempted_run": "2020-02-14T19:01:43.0003567+08:00",
                    "number_of_finished_runs": 1
                },
                "type": 0,
                "remote_properties": {
                    "url": "",
                    "method": "",
                    "body": "",
                    "headers": null,
                    "timeout": 0,
                    "expected_response_codes": null
                },
                "stats": [
                    {
                        "job_id": "5cd1fd60-e6a8-42a3-63e9-c4b4474f26b4",
                        "ran_at": "2020-02-14T19:01:43.0003567+08:00",
                        "number_of_retries": 0,
                        "success": true,
                        "execution_duration": 62436571200
                    }
                ],
                "is_done": true
            }
        }
    ]
}
```

执行 /api/v1/job/

```json
  {
    "id": "5cd1fd60-e6a8-42a3-63e9-c4b4474f26b4"
}
{
    "jobinfo": {
        "name": "Downloader ebook 垂钓之神-会狼叫的猪",
        "id": "5cd1fd60-e6a8-42a3-63e9-c4b4474f26b4",
        "command": "F:/work/ebookdownloader/ebookdownloader_cli.exe --ebhost=xsbiquge.com --bookid=90_90583 --txt --mobi --meta",
        "owner": "",
        "disabled": false,
        "dependent_jobs": null,
        "parent_jobs": null,
        "on_failure_job": "",
        "schedule": "R0/2020-02-14T19:01:43+08:00/",
        "retries": 0,
        "epsilon": "",
        "next_run_at": "2020-02-14T19:01:43+08:00",
        "resume_at_next_scheduled_time": false,
        "metadata": {
            "success_count": 1,
            "last_success": "2020-02-14T19:02:45.4369279+08:00",
            "error_count": 0,
            "last_error": "0001-01-01T00:00:00Z",
            "last_attempted_run": "2020-02-14T19:01:43.0003567+08:00",
            "number_of_finished_runs": 1
        },
        "type": 0,
        "remote_properties": {
            "url": "",
            "method": "",
            "body": "",
            "headers": null,
            "timeout": 0,
            "expected_response_codes": null
        },
        "stats": [
            {
                "job_id": "5cd1fd60-e6a8-42a3-63e9-c4b4474f26b4",
                "ran_at": "2020-02-14T19:01:43.0003567+08:00",
                "number_of_retries": 0,
                "success": true,
                "execution_duration": 62436571200
            }
        ],
        "is_done": true
    }
}
```
测试例子

```bash
$curl  -H "Authorization:Bearer xxxxxxxxx" -H "Content-Type: application/json" -X GET http://192.168.13.107:8090/api/v1/job/
$curl  -H "Authorization:Bearer xxxxxxxxx" -H "Content-Type: application/json" -X GET http://192.168.13.107:8090/api/v1/job/5cd1fd60-e6a8-42a3-63e9-c4b4474f26b4
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
