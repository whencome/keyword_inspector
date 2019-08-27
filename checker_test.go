package wordschecker

import (
	"testing"
	"fmt"
)

func init() {
	InitDict("/mnt/e/GoLibs/src/github.com/whencome/wordschecker/conf/sensitive_words.conf")
}

func Test_Evaluate(t *testing.T) {
	str := "这个是一段测试文本，欢迎加QQ123456，welcome,gay片"
	/*
	str := `本月初，美国科罗拉多州发生了一起警察开枪打死黑人青年事件，引发当地民众抗议。在黑人青年家人和一些媒体要求下，当地警方15日公布了警察执法记录仪拍摄的事发时的部分视频。


　　科罗拉多州的科罗拉多斯普林斯市警方公布的视频显示，警方在接到持枪抢劫报告后，对包括19岁黑人青年德文·贝利在内的两名黑人男子进行询问，其中一名警察要求两人举起双手接受另一名警察的搜身，以确定他们是否持有枪支。就在另一名警察准备对两人进行搜查时，贝利突然逃跑，于是警察朝他开枪。

　　对此，贝利家的律师基尔默表示，当时贝利手上并没有拿枪，却遭到警察的多次枪击。警方在追赶贝利过程中不应该使用致命力量，警方显然执法过度。



　　贝利家的律师 达罗尔德·基尔默：法律要求警方在追他的过程中，尽量不使用致命力量。但是结果却是，一个19岁的，之前没有任何犯罪记录的青年被打死了。

　　此前，警察开枪打死贝利一事已在当地引发多场民众抗议。贝利家属则希望警方能公布警察执法记录仪和行车记录仪拍摄到的事件全程视频。`
	*/
	evaluator := NewEvaluator(false)
	eval := evaluator.Evaluate(str)
	fmt.Printf("%#v\n", eval)
}
