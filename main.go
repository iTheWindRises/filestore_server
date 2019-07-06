package main

import (
	"filestore/consts"
	"filestore/handler"
	"fmt"
	"net/http"
)

func main() {
	//
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	//文件handler
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/query", handler.FileQueryHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)

	//user handler
	http.HandleFunc("/user/signup", handler.SignUpHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)

	fmt.Printf("To start server binded port%s\n", consts.ServerPort)
	err := http.ListenAndServe(consts.ServerPort, nil)

	if err != nil {
		fmt.Printf("Failed to start server err:%s", err.Error())
	}
}

func init() {
	fmt.Printf("init Server....\n")
}
