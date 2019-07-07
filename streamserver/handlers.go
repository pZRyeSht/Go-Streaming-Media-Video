package main

import (
	"io"
	"net/http"
	"os"
	"time"
	"html/template"
	"io/ioutil"
	"github.com/httprouter"
	"log"
)

func testPageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t, _:=template.ParseFiles("./videos/upload.html")
	t.Execute(w, nil)
}

func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid-id")
	vl := VIDEO_DIR + vid
	video, err := os.Open(vl)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error of open video: "+err.Error())
		return
	}
	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, "", time.Now(), video)

	defer video.Close()
}

func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//限制ioread最大能读取的文件大小或者能读到的最大的缓冲区大小
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	//首先校验文件大小
	if err:=r.ParseMultipartForm(MAX_UPLOAD_SIZE); err!=nil {
		sendErrorResponse(w, http.StatusBadRequest, "File is too large")
		return
	}

	//此处一定要这么写，否则读不出文件
	//从文件中读取
	file, _, err := r.FormFile("file")    //form name = file
	if err!=nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}

	//获取数据，写成二进制文件
	data, err:=ioutil.ReadAll(file)
	if err!=nil {
		log.Printf("Read file error: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}

	//写进将要保存的文件路径
	fn:=p.ByName("vid-id")
	err = ioutil.WriteFile(VIDEO_DIR + fn, data, 06666)
	if err!=nil {
		log.Printf("Write file error: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	//返回成功状态
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Uploaded successfully.")
}
