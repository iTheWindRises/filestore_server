package meta

import (
	mydb "filestore/dao"
	"sort"
)

//文件元信息结构
type FileMeta struct {
	FileHash string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMeta map[string]FileMeta

func init() {
	fileMeta = make(map[string]FileMeta)
}

//新增/更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMeta[fmeta.FileHash] = fmeta
}

//新增/更新文件元信息到mysql
func UpdateFileMetaDB(fm FileMeta) bool {
	return mydb.OnFileUploadFinished(
		fm.FileHash, fm.FileName, fm.FileSize, fm.Location)
}

//获取文件元信息
func GetFileMeta(hash string) FileMeta {
	return fileMeta[hash]
}

//获取文件元信息 < mysql
func GetFileMetaDB(hash string) *FileMeta {
	tf, err := mydb.GetFileMeta(hash)
	if err != nil {
		return nil
	}

	fm := &FileMeta{
		FileHash: tf.FileHash,
		FileName: tf.FileName.String,
		FileSize: tf.FileSize.Int64,
		Location: tf.FileAddr.String,
	}
	return fm
}

//删除原信息
func RemoveFileMeta(fileHash string) {
	delete(fileMeta, fileHash)
}

func (f FileMeta) IsNil() bool {
	return f == (FileMeta{})
}

//获取指定长度最新的文件原信息
func GetLastFileMetas(count int) []FileMeta {
	if count > len(fileMeta) {
		count = len(fileMeta)
	}
	fMetaArray := make([]FileMeta, len(fileMeta))

	for _, v := range fileMeta {
		fMetaArray = append(fMetaArray, v)
	}
	sort.Sort(ByUploadTime(fMetaArray))
	return fMetaArray[0:count]
}
