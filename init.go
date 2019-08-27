package wordschecker

// 字典规则前缀
var RulePrefix = "#rule#:"

// 字典文件路径，字典应保证一行一个关键词
var dictFile string
// 全局关键词字典
var dictionary *Dict
// 全局规则字典
var ruleDictionary *RuleDict


// 初始化
func init() {
	// 全局规则字典
	dictionary = NewDict(0)
	// 初始化规则字典
	ruleDictionary = NewRuleDict()
}
