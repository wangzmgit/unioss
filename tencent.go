package unioss

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type Tencent struct {
	client *cos.Client
	config *Config
}

func newTencent(config Config) (*Tencent, error) {
	tencent := Tencent{}
	err := tencent.init(config)
	if err != nil {
		return nil, err
	}
	return &tencent, err
}

func (t *Tencent) init(config Config) error {
	if t.config == nil {
		t.config = &config
	}

	if config.AppID == "" || config.Region == "" {
		return errors.New("configuration not correct")
	}

	if t.client == nil {
		c := config
		buckerURL, _ := url.Parse(fmt.Sprintf("http://%s-%s.cos.%s.myqcloud.com", c.Bucket, c.AppID, c.Region))

		client := cos.NewClient(&cos.BaseURL{BucketURL: buckerURL}, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  c.KeyID,
				SecretKey: c.KeySecret,
			},
		})
		t.client = client
	}

	return nil
}

// 获取文件
func (t *Tencent) GetObjectToFile(objectKey, filePath string) error {
	response, err := t.client.Object.Download(context.Background(), objectKey, filePath, nil)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return httpError(response)
	}

	return nil
}

// 删除文件
func (t *Tencent) DeleteObject(objectKey string) error {
	response, err := t.client.Object.Delete(context.Background(), objectKey)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return httpError(response)
	}

	return nil
}

func (t *Tencent) PutObject(objectKey string, reader io.Reader) error {
	response, err := t.client.Object.Put(context.Background(), objectKey, reader, nil)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return httpError(response)
	}

	return nil
}

// 上传文件
func (t *Tencent) PutObjectFromFile(objectKey, filePath string) error {
	response, err := t.client.Object.PutFromFile(context.Background(), objectKey, filePath, nil)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return httpError(response)
	}

	return nil
}

func (t *Tencent) IsExists(objectKey string) (bool, error) {
	return t.client.Object.IsExist(context.Background(), objectKey)
}

// cos 请求错误
func httpError(response *cos.Response) error {
	bytes, err := io.ReadAll(response.Body)
	defer func() {
		err = response.Body.Close()
	}()
	if err != nil {
		return err
	}

	return errors.New(string(bytes))
}

// 获取访问URL
func (t *Tencent) GetObjectUrl(objectKey string) string {
	if t.config.Domain == "" {
		return fmt.Sprintf("https://%s-%s.cos.%s.myqcloud.com/%s",
			t.config.Bucket,
			t.config.AppID,
			t.config.Region,
			objectKey,
		)
	}
	return fmt.Sprintf("https://%s/%s", t.config.Domain, objectKey)
}
