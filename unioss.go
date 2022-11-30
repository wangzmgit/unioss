package unioss

import (
	"errors"
	"io"
)

const (
	ALIYUN  = "aliyun"
	TENCENT = "tencent"
	QINIU   = "qiniu"
)

var currentStorage Storage

type Storage interface {
	GetObjectToFile(objectKey, downloadedFileName string) error
	DeleteObject(objectKey string) error
	PutObject(objectKey string, reader io.Reader) error
	PutObjectFromFile(objectKey, filePath string) error
	IsExists(objectKey string) (bool, error)
	GetObjectUrl(objectKey string) string
}

func NewStorage(ossName string, config Config) error {
	if config.KeyID == "" || config.KeySecret == "" || config.Bucket == "" {
		return errors.New("configuration not correct")
	}
	var err error
	switch ossName {
	case ALIYUN:
		currentStorage, err = newAliyun(config)
	case TENCENT:
		currentStorage, err = newTencent(config)
	case QINIU:
		currentStorage, err = newQiniu(config)
	default:
		return errors.New("driver not exists")
	}

	if err != nil {
		currentStorage = nil
		return err
	}

	return nil
}

func GetStorage() (Storage, error) {
	if currentStorage != nil {
		return currentStorage, nil
	}
	return nil, errors.New("driver not exists")
}
