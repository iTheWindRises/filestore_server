package handler

import (
	"encoding/json"
	"filestore/consts"
	"filestore/meta"
	util "filestore/utils"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		getUpload(w, r)
	}
	if r.Method == "POST" {
		postUpload(w, r)
	}
}

//上传成功
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Upload finished"))
}

//获得文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm()
	// fileHash := r.Form["fileHash"][0]
	fileHash := r.FormValue("fileHash")
	//fMeta := meta.GetFileMeta(fileHash)
	fMeta := meta.GetFileMetaDB(fileHash)

	if fMeta == nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

//查询批量的文件元信息
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm()
	// limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	limitCnt, _ := strconv.Atoi(r.FormValue("limit"))
	fileMetas := meta.GetLastFileMetas(limitCnt)

	json, err := json.Marshal(fileMetas)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

//文件下载
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	fileHash := r.FormValue("fileHash")
	fm := meta.GetFileMeta(fileHash)

	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", `attachment;filename=`+fm.FileName)
	w.Write(data)
}

//更新源信息文件名
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	opType := r.FormValue("op")
	fileHash := r.FormValue("fileHash")
	newFileName := r.FormValue("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	curFileMeta := meta.GetFileMeta(fileHash)
	if curFileMeta.IsNil() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//删除文件原信息
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	fileHash := r.FormValue("fileHash")

	fm := meta.GetFileMeta(fileHash)

	if !fm.IsNil() {
		os.Remove(fm.Location)
		meta.RemoveFileMeta(fileHash)
	}
	w.WriteHeader(http.StatusOK)

}

//GET
func getUpload(w http.ResponseWriter, r *http.Request) {
	// 返回上传的html页面
	file, err := ioutil.ReadFile("./static/view/index.html")
	if err != nil {
		io.WriteString(w, "internel server error")
		return
	}
	io.WriteString(w, string(file))
}

//POST
func postUpload(w http.ResponseWriter, r *http.Request) {
	//接收文件流并储存到本地
	file, header, err := r.FormFile("file")
	defer file.Close()

	if err != nil {
		fmt.Printf("Faiel to get data err:%s\n", err.Error())
		return
	}

	fileMeta := meta.FileMeta{
		FileName: header.Filename,
		Location: consts.FileSavePath + header.Filename,
		UploadAt: time.Now().Format(`2006-01-02 15:04:05`),
	}
	//创建本地文件
	newFile, err := os.Create(fileMeta.Location + fileMeta.FileName)
	defer newFile.Close()
	if err != nil {
		fmt.Printf("Faiel to create data err:%s\n", err.Error())
		return
	}

	buf := make([]byte, 1024, 1024)
	fileMeta.FileSize, err = io.CopyBuffer(newFile, file, buf)
	if err != nil {
		fmt.Printf("Faiel to save data err:%s\n", err.Error())
		return
	}

	newFile.Seek(0, 0)
	fileMeta.FileHash = util.FileSha1(newFile)
	//meta.UpdateFileMeta(fileMeta)
	ok := meta.UpdateFileMetaDB(fileMeta)
	fmt.Printf("metaFile save is :%v\n", ok)

	http.Redirect(w, r, "/file/upload/suc", http.StatusFound)

}
