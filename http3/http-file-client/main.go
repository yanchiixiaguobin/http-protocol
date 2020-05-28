package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
)

func main() {
	// 创建表单文件
	// CreateFormFile 用来创建表单，第一个参数是字段名，第二个参数是文件名
	var fileName string
	for idx, args := range os.Args {
		if idx == 1 {
			fileName = args
		}
	}

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	formFile, err := writer.CreateFormFile("uploadfile", fileName)
	if err != nil {
		log.Fatalf("Create form file failed: %s\n", err)
	}

	// 从文件读取数据，写入表单
	srcFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("%Open source file failed: s\n", err)
	}
	defer srcFile.Close()
	_, err = io.Copy(formFile, srcFile)
	if err != nil {
		log.Fatalf("Write to form file falied: %s\n", err)
	}

	var keyLog io.Writer
	pool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}

	AddRootCA(pool)
	qconf := &quic.Config{}
	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: true,
			KeyLogWriter:       keyLog,
		},
		QuicConfig: qconf,
	}

	defer roundTripper.Close()

	hclient := &http.Client{
		Transport: roundTripper,
	}

	// 发送表单
	contentType := writer.FormDataContentType()
	writer.Close() // 发送之前必须调用Close()以写入结尾行
	rsp, err := hclient.Post("https://192.168.124.16:8088/upload", contentType, buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatalf("Got response for %s: %#v", "http://192.168.124.16:8088/upload", rsp)
}
