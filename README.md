# go 对象存储库

支持`阿里云`、`腾讯云`、`七牛云`

## 安装
```
go get github.com/wangzmgit/uni-oss
```

## 基本用法
```
初始化存储，传入 `云存储类型` 和 `配置信息`
云存储类型可选值 ALIYUN TENCENT QINIU
NewStorage(unioss.ALIYUN,unioss.Config{

})
```
配置信息内容
```

KeyID     string //必填
对应阿里云access_key_id、腾讯secret_id，七牛云accessKey

KeySecret string //必填
对应阿里云access_key_secret、讯云secret_key，七牛云secretKey

Bucket    string //必填 
存储桶名称

Endpoint string //阿里云必填
具体参考阿里云访问域名和数据中心https://help.aliyun.com/document_detail/31837.htm

AppID  string //腾讯云必填
存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket

Region string //腾讯云必填
存储桶region可以在COS控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 


Domain string //七牛云必填（可以是七牛空间的默认域名）其他如果设置自定义域名可填写

Private bool //是否私有(仅七牛云)
```

## 使用示例
```
func TestOSS() {
	err := NewStorage(ALIYUN, Config{
		KeyID:     "",
		KeySecret: "",
		Endpoint:  "",
		Bucket:    "",
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	s, err := GetStorage()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(s.IsExists("test.txt"))
}

```

## 支持操作
```
// 获取文件到本地
GetObjectToFile(objectKey, filePath string) error
// 删除文件
DeleteObject(objectKey string) error
// 上传本地文件,第一个参数是 objectKey，第二个参数是 `io.Reader`。
Put(objectKey string, reader io.Reader) error
// 上传本地文件,第一个参数是 objectKey，第二个参数是本地文件路径。
PutObjectFromFile(objectKey, filePath string) error
// 获取文件是否存在
IsExists(objectKey string) (bool, error)
//获取文件url(公开)
GetObjectUrl(objectKey string) string
```


## 参考文档

[阿里云对象存储](https://help.aliyun.com/document_detail/32143.html)

[腾讯云对象存储](https://cloud.tencent.com/document/product/436/31215)

[七牛云对象存储](https://developer.qiniu.com/kodo/1238/go)




