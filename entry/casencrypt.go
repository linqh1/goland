package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"
)

const apiUser = "ms-portal"
const apiKey  = "20ED46FCD81E99E8EE6CEE97777401F2"
//时间必须是2016-01-02 下午3点4分5秒才能被正确识别
const timeLayout = "Mon, 02 Jan 2006 15:04:05 "

func main() {
	date := time.Now().UTC().Format(timeLayout + "GMT")
	result := base64.StdEncoding.EncodeToString(SignHmacSha1([]byte(apiKey),[]byte(date)))
	result = "Basic " + base64.StdEncoding.EncodeToString([]byte("ms-portal:"+result))
	fmt.Println(date)
	fmt.Println(result)
}

func GmtTimeFormat() string{
	//result := time.Now()
	//_, offset := result.Zone()
	//result = result.Add(-time.Second * time.Duration(offset)) //补偿当前时区和GMT时区的时间偏差
	//return result.Format(timeLayout + "GMT")
	return time.Now().UTC().Format(timeLayout + "GMT")
}

func SignHmacSha1(key []byte,value []byte) []byte {
	mac := hmac.New(sha1.New, []byte(apiKey))
	mac.Write(value)
	return mac.Sum(nil)
}
