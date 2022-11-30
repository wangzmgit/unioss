package unioss

import (
	"fmt"
	"testing"
)

func TestOSS(t *testing.T) {
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
