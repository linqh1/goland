package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//时间必须是2016-01-02 下午3点4分5秒才能被正确识别
const timeLayout = "Mon, 02 Jan 2006 15:04:05 "

var (
	confFile        string
	apiuser         string
	apikey          string
	keyUserSet      []EncryptSet
	DefaultFilePath = "casencrypt-default.config"
)

type EncryptSet struct {
	key  string
	user string
	doc  string
}

func init() {
	flag.StringVar(&apiuser, "u", "", "api user")
	flag.StringVar(&apikey, "k", "", "api key")
	flag.StringVar(&confFile, "f", "", "read key and user from given file. file \"casencrypt-default.config\" at same directory with script will be use if run without any flag.")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s [options]:\n", os.Args[0])
		flag.PrintDefaults()
	}
	keyUserSet = make([]EncryptSet, 0)
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	DefaultFilePath = filepath.Join(dir, DefaultFilePath)
}

func main() {
	flag.Parse()
	keyIsSet, userIsSet := apikey != "", apiuser != ""
	if bool2int(keyIsSet)^bool2int(userIsSet) > 0 {
		fmt.Fprintf(os.Stderr, "[warning] user and key should both be set\n")
	}
	if keyIsSet && userIsSet {
		keyUserSet = append(keyUserSet, EncryptSet{apikey, apiuser, ""})
	}

	if confFile == "" && len(keyUserSet) == 0 {
		fmt.Printf("no flag provided. use default config: %v\n", DefaultFilePath)
		confFile = DefaultFilePath
	}

	if confFile != "" {
		readFromFile(confFile)
	}

	date := time.Now().UTC().Format(timeLayout + "GMT")
	for _, s := range keyUserSet {
		fmt.Println()
		result := encrypt(s, date)
		fmt.Printf("key:%v, user:%v doc:%v\n", s.key, s.user, s.doc)
		fmt.Printf("Date:%v\n", date)
		fmt.Printf("Authorization:%v\n", result)
	}

}

func readFromFile(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not open file[%v].error:%v\n", filepath, err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(file)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		lineString := scanner.Text()
		if strings.TrimSpace(lineString) == "" { // 空行不计算
			continue
		}
		fields := strings.Fields(lineString)
		if len(fields) < 2 {
			fmt.Fprintf(os.Stderr, "[error] format error in line:%v \n", lineNo)
			continue
		}
		set := EncryptSet{key: fields[1], user: fields[0]}
		if len(fields) > 2 {
			set.doc = strings.Join(fields[2:], " ")
		}
		keyUserSet = append(keyUserSet, set)

	}
}

func bool2int(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func encrypt(set EncryptSet, date string) string {
	result := base64.StdEncoding.EncodeToString(SignHmacSha1([]byte(set.key), []byte(date)))
	result = "Basic " + base64.StdEncoding.EncodeToString([]byte(set.user+":"+result))
	return result
}

func SignHmacSha1(key []byte, value []byte) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write(value)
	return mac.Sum(nil)
}
