package dao

import (
	"database/sql"
	mydb "filestore/dao/mysql"
	"fmt"
)

//保存文件元信息到mysql
func OnFileUploadFinished(fileHash string, fileName string,
	fileSize int64, fileAddr string) bool {

	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_file(file_hash,file_name,file_size,file_addr,status) values(?,?,?,?,1)",
	)
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(fileHash, fileName, fileSize, fileAddr)
	if err != nil {
		fmt.Println("Failed to exec statement, err:" + err.Error())
		return false
	}

	if rf, err := ret.RowsAffected(); err == nil {
		if rf <= 0 {
			fmt.Printf("Failed with gash:%s has been uploaded before\n", fileHash)
			return false
		}
		return true
	}

	return false
}

type TableFileMeta struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

//查询文件元信息
func GetFileMeta(fileHash string) (*TableFileMeta, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_hash,file_name,file_size, file_addr " +
			"from tbl_file where file_hash=? and status=1 limit 1")

	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return nil, err
	}
	defer stmt.Close()

	tf := &TableFileMeta{}

	err = stmt.QueryRow(fileHash).Scan(&tf.FileHash, &tf.FileName, &tf.FileSize, &tf.FileAddr)
	if err != nil {
		fmt.Println("Failed to Select, err:" + err.Error())
		return nil, err
	}
	return tf, nil
}
