package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
	"github.com/lucas-clemente/quic-go/quictrace"
)

var tracer quictrace.Tracer

func init() {
	tracer = quictrace.NewTracer()
}

type tracingHandler struct {
	handler http.Handler
}

var _ http.Handler = &tracingHandler{}

func (h *tracingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handler.ServeHTTP(w, r)
	if err := exportTraces(); err != nil {
		log.Fatal(err)
	}
}

func exportTraces() error {
	traces := tracer.GetAllTraces()
	if len(traces) != 1 {
		return errors.New("expected exactly one trace")
	}
	for _, trace := range traces {
		f, err := os.Create("trace.qtr")
		if err != nil {
			return err
		}
		if _, err := f.Write(trace); err != nil {
			return err
		}
		f.Close()
		fmt.Println("Wrote trace to", f.Name())
	}
	return nil
}

func main() {
	bind := "192.168.124.16:8088"
	handler := setupHandler()
	quicConf := &quic.Config{
		MaxIncomingStreams: 10000,
		MaxIncomingUniStreams: 10000,
		MaxReceiveStreamFlowControlWindow: 100 * (1 << 20), // 100M
		MaxReceiveConnectionFlowControlWindow: 100 * (1 << 20), // 100M
	}
	quicConf.QuicTracer = tracer

	server := http3.Server{
		Server:     &http.Server{Handler: handler, Addr: bind},
		QuicConfig: quicConf,
	}
	server.ListenAndServeTLS(GetCertificatePaths())
}

func setupHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", uploadHandle)
	return &tracingHandler{handler: mux}
}

func uploadHandle(w http.ResponseWriter, r *http.Request) {

	// 根据字段名获取表单文件
	formFile, header, err := r.FormFile("uploadfile")
	if err != nil {
		log.Printf("Get form file failed: %s\n", err)
		return
	}
	defer formFile.Close()

	// 获取当前目录
	directory, _ := os.Getwd()

	// 创建保存文件
	destFile, err := os.Create(directory + "/" + header.Filename)
	if err != nil {
		log.Printf("Create failed: %s\n", err)
		return
	}
	defer destFile.Close()

	// 读取表单文件，写入保存文件
	_, err = io.Copy(destFile, formFile)
	if err != nil {
		log.Printf("Write file failed: %s\n", err)
		return
	}
}
