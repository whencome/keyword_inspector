package wordschecker

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

// 初始化
func InitDict(dictFile string) {
	SetDictFile(dictFile)
	LoadDicts()
}

// 设置字典文件
func SetDictFile(path string) {
	dictFile = path
}

// 加载字典信息
func LoadDicts() {
	if dictionary != nil && dictionary.Size() > 0 {
		return
	}
	ReloadDicts()
}

// 重新加载字典信息
func ReloadDicts() {
	// 初始化字典
	dictionary = NewDict(0)
	ruleDictionary = NewRuleDict()
	// 检查字典文件是否存在
	// 如果字典文件不存在，则保留一个空字典，防止后续调用报错
	if dictFile == "" || !isFileExists(dictFile) {
		return
	}
	// 加载字典信息
	file, err := os.Open(dictFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// 读取文件内容
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			// 文件读取结束
			if err == io.EOF {
				break
			}
			// 暂时忽略其他错误
			log.Println(err)
			continue
		}
		line = strings.TrimSpace(line)
		di := NewDictItem(line)
		if di != nil {
			AddDictItem(di)
		}
	}
}

// 检查文件是否存在
func isFileExists(file string) bool {
	fInfo, err := os.Stat(file)
	if err == nil && !fInfo.IsDir() {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// 检查给定路径是否存在
func isPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
