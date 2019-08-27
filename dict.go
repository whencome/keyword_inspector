package wordschecker

import (
	"strings"
	"regexp"
	"strconv"
)

// 普通关键词字典
type Dict struct {
	Depth		int
	Weight 		int 	// 权重
	Keyword		map[rune]*Dict
}

// 正则表达式规则词典
type RuleDict struct {
	Rules 		[]*RuleItem
}

type RuleItem struct {
	Weight 			int 	// 权重
	Rule 			string
	Pattern 		*regexp.Regexp
}

// 定义字典条目
type DictItem struct {
	OriginText 		string   	// 原始输入项目
	Keyword 		string 		// 字典内容
	Weight 			int 		// 权重信息
	IsNote			bool		// 是否是注释内容，注释内容不进入字典列表
	IsRule			bool 		// 是否是规则词条，规则词条进入RuleDict，使用正则匹配
}

// 创建普通词典
func NewDict(depth int) *Dict {
	return &Dict{
		Depth : depth,
		Weight : 1,
		Keyword : make(map[rune]*Dict),
	}
}

// 创建正则词典
func NewRuleDict() *RuleDict {
	return &RuleDict{
		Rules : make([]*RuleItem, 0),
	}
}

// 添加一个关键词
func (dict *Dict) Add(word []rune, weight int) {
	size := len(word)
	if size == 0 {
		return
	}
	// 添加关键词到链表中
	var newDict *Dict = dict
	for pos, char := range word {
		if _, ok := newDict.Keyword[char]; !ok {
			newDict.Keyword[char] = NewDict(newDict.Depth + 1)
		}
		newDict.Keyword[char].Add(word[pos+1:], weight)
		newDict = newDict.Keyword[char]
		if pos + 1 == size {
			newDict.Weight = weight
			break
		}
	}
}

// 添加关键词
func (dict *Dict) AddKeyword(word string, weight int) {
	dict.Add([]rune(strings.TrimSpace(word)), weight)
}

// 根据key查询对应的dict
func (dict *Dict) Find(word rune) *Dict {
	d, ok := dict.Keyword[word]
	if !ok {
		return nil
	}
	return d
}

// 获取词典大学
func (dict *Dict) Size() int {
	return len(dict.Keyword)
}

// 添加一条规则
func (rd *RuleDict) Add(rule string, weight int) {
	pattern, err := regexp.Compile(rule)
	if err != nil {
		return
	}
	rd.Rules = append(rd.Rules, NewRuleItem(rule, pattern, weight))
}

// 获取词典大学
func (rd *RuleDict) Size() int {
	return len(rd.Rules)
}

// 创建一个RuleItem
func NewRuleItem(rule string, pattern *regexp.Regexp, weight int) *RuleItem {
	return &RuleItem{
		Weight 	: weight,
		Rule 	: rule,
		Pattern	: pattern,
	}
}

// 创建一个DictItem
func NewDictItem(text string) *DictItem {
	di := &DictItem{
		OriginText : strings.TrimSpace(text),
	}
	// 解析内容
	if di.OriginText == "" {
		return di
	}
	// 以;开头的行表示注释，直接忽略
	if strings.HasPrefix(di.OriginText, ";") {
		di.IsNote = true
	}
	if strings.HasPrefix(di.OriginText, RulePrefix) {
		di.IsRule = true
		di.OriginText = di.OriginText[len(RulePrefix):]
	}
	// 获取权重信息
	weight := 1
	keyword := di.OriginText
	match, err := regexp.MatchString(`\=\>\s*\d+$`, di.OriginText)
	// 如果出错，则将此条目设置为注释，直接忽略
	if err != nil {
		di.IsNote = true
		return di
	}
	if match {
		parts := strings.Split(di.OriginText, "=>")
		if len(parts) >= 2 {
			keyword = strings.TrimSpace(parts[0])
			weightStr := parts[len(parts)-1]
			w, err := strconv.Atoi(weightStr)
			if err == nil && w > 0 {
				weight = w
			}
		}
	}
	di.Weight = weight
	di.Keyword = keyword
	// 返回结果
	return di
}

// 添加词条
func AddDictItem(di *DictItem) {
	// 如果是注释，则直接忽略
	if di.IsNote {
		return
	}
	// 根据类型添加
	if !di.IsRule {
		dictionary.AddKeyword(di.Keyword, di.Weight)
		return
	}
	// 添加规则
	ruleDictionary.Add(di.Keyword, di.Weight)
}
