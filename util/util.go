package util

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
)

type Sha1Stream struct {
	_sha1 hash.Hash
}

func (obj *Sha1Stream) Update(data []byte) {
	if obj._sha1 == nil {
		obj._sha1 = sha1.New()
	}
	obj._sha1.Write(data)
}

func (obj *Sha1Stream) Sum() string {
	return hex.EncodeToString(obj._sha1.Sum([]byte("")))
}

func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum(nil))
}

func Sha256(data []byte) string {
	_sha256 := sha256.New()
	_sha256.Write(data)
	return hex.EncodeToString(_sha256.Sum(nil))
}

//使用文件的二进制数据，进行SHA1的算法加密，得到的hash值用于作文件的ID
func FileSha1(file *os.File) string {
	_sha1 := sha1.New()  //生成一个hash.Hash接口

	//官方用法：io.WriteString(h, "His money is twice tainted:")，把一串字符串加密
	//这是对file对象进行SHA1加密
	io.Copy(_sha1, file)

	//Sum函数为_sha1对象的校验和，即最终的Hash函数，如果已经使用New()的方法，那么就不需要Sum的时候再传入参数了，所以参数为nil
	return hex.EncodeToString(_sha1.Sum(nil))
}


func MD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte(nil)))
}

func FileMD5(file *os.File) string {
	_md5 := md5.New()
	io.Copy(_md5, file)
	return hex.EncodeToString(_md5.Sum(nil))
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}
