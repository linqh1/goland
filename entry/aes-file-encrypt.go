package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const PasswordLength = 32

var (
	prefixBytes  = []byte("my en/decrypt")
	encryptMode  bool
	decryptMode  bool
	password     string
	inputFiles   []string
	bufferSize   uint
	inputBuffer  []byte
	outputBuffer []byte
	aesBlocker   cipher.Block
)

func main() {
	checkAndInit()
	for _, file := range inputFiles {
		dir, filePattern := filepath.Split(file)
		if dir == "" {
			dir = "."
		}
		i := 0
		_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			i++
			if i == 1 || info == nil { // i=1不对本目录进行处理
				return nil
			}
			if info.IsDir() {
				return filepath.SkipDir // 不扫描子目录
			}
			fileName := info.Name()
			match, _ := filepath.Match(filePattern, fileName)
			if !match {
				return nil
			}
			fmt.Printf("%v: ", fileName)
			fullName := filepath.Join(dir, fileName)
			fileObj, err := os.Open(fullName)
			if err != nil {
				fmt.Printf("open failed:%v\n", err)
				return nil
			}
			if encryptMode {
				outFile, err := os.Create(fileName + ".encrypt")
				if err != nil {
					fmt.Printf("create output file failed:%v\n", err)
					return nil
				}
				_, _ = outFile.Write(prefixBytes)
				handlerFile(fileObj, outFile, aesBlocker.Encrypt)
			} else {
				prefix := make([]byte, len(prefixBytes))
				fileObj.Read(prefix)
				if !bytes.Equal(prefix, prefixBytes) {
					fmt.Printf("can not descrypt this file\n")
					return nil
				}
				outFile, err := os.Create(fileName + ".decrypt")
				if err != nil {
					fmt.Printf("create output file failed:%v\n", err)
					return nil
				}
				handlerFile(fileObj, outFile, aesBlocker.Decrypt)
			}
			return nil
		})
	}
}

func handlerFile(file *os.File, outFile *os.File, bytesHandler func(dst, src []byte)) {
	defer file.Close()
	defer outFile.Close()
	outputWriter := bufio.NewWriterSize(outFile, int(bufferSize))
	defer outputWriter.Flush()
	for {
		n, err := file.Read(inputBuffer)
		if err != nil && err != io.EOF {
			fmt.Printf("read failed:%v\n", err)
			return
		}
		outputWriter.Write(handleBytes(inputBuffer[:n], bytesHandler))
		if err == io.EOF {
			return
		}
	}
}

func handleBytes(bytes []byte, handler func(dst, src []byte)) []byte {
	length := len(bytes)
	res := make([]byte, length)
	times := length / 16
	for i := 0; i < times; i++ {
		handler(res[i*16:(i+1)*16], bytes[i*16:(i+1)*16])
	}
	if mod := length % 16; mod != 0 { // 不足16位的不进行加解密
		for mod > 0 {
			res[length-mod] = bytes[length-mod]
			mod--
		}
	}
	return res
}

func checkAndInit() {
	if !encryptMode && !decryptMode {
		fmt.Println("please choose one mode: encrypt or decrypt")
		flag.Usage()
		os.Exit(1)
	}
	if encryptMode && decryptMode {
		fmt.Println("can only choose one mode: encrypt or decrypt")
		flag.Usage()
		os.Exit(1)
	}
	if password == "" {
		fmt.Println("please specify a password")
		flag.Usage()
		os.Exit(1)
	}
	inputFiles = flag.Args()
	if len(inputFiles) <= 0 {
		fmt.Println("please specify file to deal with")
		flag.Usage()
		os.Exit(1)
	}
	bytes := []byte(password)
	if len(bytes) > PasswordLength {
		fmt.Println("password too long")
		os.Exit(1)
	}
	newPasswordBytes := make([]byte, PasswordLength)
	copy(newPasswordBytes, bytes)
	remain := PasswordLength - len(bytes)
	for remain > 0 { // 不足32位填充0值
		newPasswordBytes[PasswordLength-remain] = 0
		remain--
	}
	block, err := aes.NewCipher(newPasswordBytes)
	if err != nil {
		fmt.Printf("invalid password:%v\n", err)
		os.Exit(1)
	}
	aesBlocker = block
	if mod := bufferSize % 16; mod != 0 { // 取16的整数
		bufferSize += 16 - mod
	}
	inputBuffer = make([]byte, bufferSize)
	outputBuffer = make([]byte, bufferSize)
}

func init() {
	flag.BoolVar(&encryptMode, "encrypt", false, "encrypt file")
	flag.BoolVar(&decryptMode, "decrypt", false, "decrypt file")
	flag.StringVar(&password, "password", "", "password for encrypt/decrypt")
	flag.UintVar(&bufferSize, "buffer", 1024, "read buffer")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "A simple tool to encrypt/decrypt file\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s file1 file2 file3\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
}
