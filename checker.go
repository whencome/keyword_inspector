package wordschecker

import (
	"strings"
	"sync"
)

// 评估对象
type Evaluator struct {
}

type evaluateScore struct {
	Score int
	Words []rune
	Info  string
}

type MatchWord struct {
	Keyword string
	Score   int
	Count   int
}

// EvaluateResult 评估结果
type EvaluateResult struct {
	Score    int
	Words    map[string]*MatchWord
	evalChan chan *evaluateScore
	stopChan chan struct{}
}

func newMatchWord(word string) *MatchWord {
	return &MatchWord{
		Keyword: word,
		Score:   0,
		Count:   0,
	}
}

func newEvaluateResult() *EvaluateResult {
	return &EvaluateResult{
		Score:    0,
		Words:    make(map[string]*MatchWord),
		evalChan: make(chan *evaluateScore),
		stopChan: make(chan struct{}),
	}
}

func newEvaluateScore() *evaluateScore {
	return &evaluateScore{
		Score: 0,
		Words: make([]rune, 0),
		Info:  "",
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

// NewEvaluator 创建一个新的Evaluator对象
func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// Evaluate 对输入语句进行评估检测
func (e *Evaluator) Evaluate(sentence string) *EvaluateResult {
	evalResult := newEvaluateResult()
	sentence = strings.TrimSpace(sentence)
	words := []rune(sentence)
	if len(words) == 0 {
		return evalResult
	}

	// 计算结果
	go func() {
		for {
			select {
			case es := <-evalResult.evalChan:
				if es.Score > 0 {
					evalResult.Score += es.Score
					keyword := string(es.Words)
					if _, ok := evalResult.Words[keyword]; !ok {
						evalResult.Words[keyword] = newMatchWord(keyword)
					}
					evalResult.Words[keyword].Score += es.Score
					evalResult.Words[keyword].Count++
				}
			case <-evalResult.stopChan:
				close(evalResult.evalChan)
				return
			}
		}
	}()

	// 开始检查
	var wg sync.WaitGroup
	wg.Add(2)
	// 检查关键词
	go func() {
		evalKeyword(evalResult.evalChan, words)
		wg.Done()
	}()
	// 规则检查
	go func() {
		evalRules(evalResult.evalChan, sentence)
		wg.Done()
	}()
	// 等待结束
	wg.Wait()
	evalResult.stopChan <- struct{}{}

	// 返回结果
	return evalResult
}

// 关键词匹配
func evalKeyword(scoreChan chan *evaluateScore, words []rune) {
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
		checkerPool <- struct{}{}
		go func(dict *Dict, word rune, chars []rune) {
			evalScore := newEvaluateScore()
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
func evalRules(scoreChan chan *evaluateScore, words string) {
	if ruleDictionary.Size() == 0 {
		return
	}
	for _, regRule := range ruleDictionary.Rules {
		// 获取匹配内容
		matches := regRule.Pattern.FindSubmatch([]byte(words))
		if matches == nil {
			continue
		}
		evalScore := newEvaluateScore()
		matchPart := string(matches[0])
		evalScore.SetWords([]rune(matchPart))
		evalScore.SetScore(regRule.Weight)
		scoreChan <- evalScore
	}
}
