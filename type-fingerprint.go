package appfinger

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type FingerPrint struct {
	//指纹适用范围
	//Protocol string
	//Path     string

	//指纹适配产品
	ProductName string

	//指纹识别规则
	Keyword *Expression
}

var regxKeywordFinger = regexp.MustCompile(`^([^ \t]+)[ \t](.*)$`)

func NewFingerPrint(s string) (*FingerPrint, error) {
	//去除表达式尾部空格
	s = strings.TrimSpace(s)
	//对表达式合法性进行校验
	if regxKeywordFinger.MatchString(s) == false {
		return nil, errors.New(" syntax error :" + s)
	}
	//校验合法，序列化表达式
	i := regxKeywordFinger.FindStringSubmatch(s)
	productName := i[1]
	expressionString := i[2]
	expression, err := NewExpression(expressionString)
	if err != nil {
		return nil, err
	}

	return &FingerPrint{
		//指纹适用范围
		//Protocol: "",
		//Path:     "",
		//指纹适配产品
		ProductName: productName,
		//指纹识别规则
		//Icon:    "",
		//Hash:    "",
		Keyword: expression,
	}, nil
}

func (f *FingerPrint) Match(banner *Banner) string {
	if f.Keyword != nil {
		if f.Keyword.Match(banner) {
			fmt.Println(f.Keyword.MakeBoolExpression(banner))
			fmt.Println(f.Keyword.value)
			fmt.Println(f.Keyword.expr)
			fmt.Println(f.ProductName)
			return f.ProductName
		}
	}
	return ""
}
