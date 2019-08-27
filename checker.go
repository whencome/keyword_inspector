package wordschecker

import (
	"strings"
	"sync"
)

// 评估对象
type Evaluator struct {
	// 开启增强检查(将中文转换成英文进行检查)
	EnableExtraCheck		bool
}

type evaluateScore struct {
	Score 			int
	Words			[]rune
	IsExtraCheck	bool
	Info 			string
}

type MatchWord struct {
	Keyword 		string
	Score 			int
	Count 			int
	ExtraScore		int
	ExtraCount		int
}

// 评估结果
type EvaluateResult struct {
	Score 			int
	Words			map[string]*MatchWord
	evalChan		chan *evaluateScore
	stopChan		chan struct{}
}

func NewMatchWord(word string) *MatchWord{
	return &MatchWord{
		Keyword		: word,
		Score		: 0,
		Count		: 0,
		ExtraScore	: 0,
		ExtraCount	: 0,
	}
}

func NewEvaluateResult() *EvaluateResult {
	return &EvaluateResult{
		Score 		: 0,
		Words 		: make(map[string]*MatchWord),
		evalChan	: make(chan *evaluateScore),
		stopChan  	: make(chan struct{}),
	}
}

func newEvaluateScore(isExtraCheck bool) *evaluateScore {
	return &evaluateScore{
		Score		 : 0,
		Words		 : make([]rune, 0),
		IsExtraCheck : isExtraCheck,
		Info : "",
	}
}

func (es *evaluateScore) SetScore(score int) {
	es.Score = score
}

func (es *evaluateScore) SetWords(words []rune) {
	es.Words = make([]rune, len(words))
	copy(es.Words, words)
}

func (es *evaluateScore) AddWord(word rune) {
	es.Words = append(es.Words, word)
}


func NewEvaluator(extraCheck bool) *Evaluator {
	return &Evaluator{
		EnableExtraCheck : extraCheck,
	}
}

func (e *Evaluator) Evaluate(sentence string) *EvaluateResult {
	evalResult := NewEvaluateResult()
	sentence = strings.TrimSpace(sentence)
	words := []rune(sentence)
	if len(words) == 0 {
		return evalResult
	}

	// 计算结果
	go func(){
		for {
			select {
			case es := <-evalResult.evalChan:
				if es.Score > 0 {
					evalResult.Score += es.Score
					keyword := string(es.Words)
					if _, ok := evalResult.Words[keyword]; !ok {
						evalResult.Words[keyword] = NewMatchWord(keyword)
					}
					evalResult.Words[keyword].Score += es.Score
					evalResult.Words[keyword].Count += 1
					if es.IsExtraCheck {
						evalResult.Words[keyword].ExtraScore += es.Score
						evalResult.Words[keyword].ExtraCount += 1
					}
				}
			case <-evalResult.stopChan:
				return
			}
		}
	}()

	// 开始检查
	var wg sync.WaitGroup
	wg.Add(2)
	// 检查关键词
	go func(){
		evalKeyword(evalResult.evalChan, words, false)
		wg.Done()
	}()
	// 规则检查
	go func(){
		evalRules(evalResult.evalChan, sentence, false)
		wg.Done()
	}()
	// 等待结束
	wg.Wait()
	evalResult.stopChan <- struct{}{}

	// 返回结果
	return evalResult
}

// 关键词匹配
func evalKeyword(scoreChan chan *evaluateScore, words []rune, isExtraCheck bool) {
	var innerWg sync.WaitGroup
	var checkerPool = make(chan struct{}, 5)
	for pos, word := range words {
		innerWg.Add(1)
		dict := dictionary.Find(word)
		if dict == nil {
			innerWg.Done()
			continue
		}
		// 关键词检查
		checkerPool<-struct{}{}
		go func(dict *Dict, word rune, chars []rune) {
			evalScore := newEvaluateScore(isExtraCheck)
			evalScore.AddWord(word)
			for _, char := range chars {
				newDict := dict.Find(char)
				if newDict == nil {
					break
				}
				evalScore.AddWord(char)
				dict = newDict
			}
			if len(dict.Keyword) == 0 && len(evalScore.Words) > 1 {
				evalScore.SetScore(dict.Weight)
				scoreChan <- evalScore
			}
			<-checkerPool
			innerWg.Done()
		}(dict, word, words[pos+1:])
	}
	innerWg.Wait()
}

// 规则匹配
func evalRules(scoreChan chan *evaluateScore, words string, isExtraCheck bool) {
	if ruleDictionary.Size() == 0 {
		return
	}
	for _, regRule := range ruleDictionary.Rules {
		// 获取匹配内容
		matches := regRule.Pattern.FindSubmatch([]byte(words))
		if matches == nil {
			continue
		}
		evalScore := newEvaluateScore(isExtraCheck)
		// 第一个匹配的是匹配内容本身，所以从第二个（下标1）开始取值
		// 暂时只返回第一个内容
		matchPart := string(matches[0])
		evalScore.SetWords([]rune(matchPart))
		evalScore.SetScore(regRule.Weight)
		scoreChan <- evalScore
	}
}
