package unioss

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

const (
	TempFileSuffix = ".temp" // Temp file suffix
)

type Qiniu struct {
	config        *Config
	bucketManager *storage.BucketManager
}

func newQiniu(config Config) (*Qiniu, error) {
	qiniu := Qiniu{}
	err := qiniu.init(config)
	if err != nil {
		return nil, err
	}
	return &qiniu, nil
}

func (q *Qiniu) init(config Config) error {
	if q.config == nil {
		q.config = &config
	}
	if q.bucketManager == nil {
		mac := qbox.NewMac(q.config.KeyID, q.config.KeySecret)
		cfg := storage.Config{
			UseHTTPS: true,
		}
		q.bucketManager = storage.NewBucketManager((*auth.Credentials)(mac), &cfg)
	}
	return nil
}

// 获取文件
func (q *Qiniu) GetObjectToFile(objectKey, filePath string) error {
	url := q.getDownloadUrl(objectKey)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	tempFilePath := filePath + TempFileSuffix
	fd, err := os.OpenFile(tempFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0664))
	if err != nil {
		return err
	}

	_, err = io.Copy(fd, resp.Body)
	fd.Close()
	if err != nil {
		return err
	}

	return os.Rename(tempFilePath, filePath)
}

// 删除文件
func (q *Qiniu) DeleteObject(objectKey string) error {
	return q.bucketManager.Delete(q.config.Bucket, objectKey)
}

func (q *Qiniu) PutObject(objectKey string, reader io.Reader) error {
	upToken := q.getUpToken()
	formUploader := q.getUploader()
	ret := storage.PutRet{}

	buf := new(bytes.Buffer)
	size, err := io.Copy(buf, reader)
	if err != nil {
		return err
	}

	err = formUploader.Put(context.Background(), &ret, upToken, objectKey, buf, size, nil)
	if err != nil {
		return err
	}

	return err
}

func (q *Qiniu) PutObjectFromFile(objectKey, filePath string) error {
	upToken := q.getUpToken()
	formUploader := q.getUploader()
	ret := storage.PutRet{}

	return formUploader.PutFile(context.Background(), &ret, upToken, objectKey, filePath, nil)
}

func (q *Qiniu) IsExists(objectKey string) (bool, error) {
	_, err := q.bucketManager.Stat(q.config.Bucket, objectKey)
	if err != nil {
		return false, err
	}

	return true, nil
}

// 获取下载链接
func (q *Qiniu) getDownloadUrl(key string) string {
	var url string

	if q.config.Private {
		mac := qbox.NewMac(q.config.KeyID, q.config.KeySecret)
		deadline := time.Now().Add(time.Second * 3600).Unix() // 1小时有效期
		url = storage.MakePrivateURL((*auth.Credentials)(mac), q.config.Domain, key, deadline)
	} else {
		url = storage.MakePublicURL(q.config.Domain, key)
	}

	return url
}

// 获取upToken
func (q *Qiniu) getUpToken() string {
	putPolicy := storage.PutPolicy{
		Scope: q.config.Bucket,
	}
	mac := qbox.NewMac(q.config.KeyID, q.config.KeySecret)
	return putPolicy.UploadToken((*auth.Credentials)(mac))
}

// 获取Uploader
func (q *Qiniu) getUploader() *storage.FormUploader {
	cfg := storage.Config{}
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	return storage.NewFormUploader(&cfg)
}

// 获取访问URL
func (q *Qiniu) GetObjectUrl(objectKey string) string {
	return fmt.Sprintf("https://%s/%s", q.config.Domain, objectKey)
}
