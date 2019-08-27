package wordschecker

import (
	"os"
	"log"
	"bufio"
	"io"
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
		/*
		// 以;开头的行表示注释，直接忽略
		if strings.HasPrefix(line, ";") {
			continue
		}
		// 以#rule#:开头的表示规则，添加到规则字典列表
		if strings.HasPrefix(line, RulePrefix) {
			ruleDictionary.Add(line[len(RulePrefix):])
			continue
		}
		// 添加关键词
		dictionary.AddKeyword(line)
		*/
		di := NewDictItem(line)
		AddDictItem(di)
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
