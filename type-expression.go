package appfinger

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Expression struct {
	//表达式数组
	paramSlice []*Param
	//表达式原文
	value string
	//表达式逻辑字符串
	//value = (body="test" || header="tt") && response="aaaa"
	//expr  = (${1} || ${2}) && ${3}
	expr string
}

func NewExpression(expr string) (*Expression, error) {
	e := &Expression{}
	e.value = expr
	//修饰expr
	expr = strings.ReplaceAll(expr, `\"`, `[quota]`)
	//合法性校验
	if err := exprVerification(expr); err != nil {
		return nil, err
	}
	//提取param数组
	var paramSlice []*Param
	paramRawSlice := paramRegx.FindAllStringSubmatch(expr, -1)
	//对param进行解析
	for index, value := range paramRawSlice {
		expr = strings.Replace(expr, value[0], "${"+strconv.Itoa(index+1)+"}", 1)
		param, err := NewParam(value[0])
		if err != nil {
			return nil, err
		}

		paramSlice = append(paramSlice, param)
	}

	e.expr = expr
	e.paramSlice = paramSlice
	return e, nil
}

func (e *Expression) Match(banner *Banner) bool {
	expr := e.MakeBoolExpression(banner)
	b, _ := ParseBoolFromString(expr)
	return b
}

func (e *Expression) MakeBoolExpression(banner *Banner) string {
	var expr = e.expr
	for index, param := range e.paramSlice {
		b := param.Match(banner)
		expr = strings.Replace(expr, "${"+strconv.Itoa(index+1)+"}", strconv.FormatBool(b), 1)
	}
	return expr
}

func exprVerification(expr string) error {
	//把所有param替换为空
	str := paramRegx.ReplaceAllString(expr, "")
	//把所有逻辑字符替换为空
	str = regexp.MustCompile(`[&| ()]`).ReplaceAllString(str, "")
	//检测是否存在其他字符
	if str != "" {
		str = strings.ReplaceAll(str, `[quota]`, `\"`)
		return errors.New(strconv.Quote(str) + " is unknown")
	}
	return nil
}

type Param struct {
	keyword  string
	value    string
	operator int
}

const (
	Unequal    = iota // !=
	Equal             // =
	RegxEqual         // ~=
	SuperEqual        // ==
)

var keywordSlice = []string{
	"Title",
	"Header",
	"Body",
	"Response",
	"Protocol",
	"Cert",
	"Port",
	"Hash",
	"Icon",
}

var paramRegx = regexp.MustCompile(`([a-zA-Z0-9]+) *(!=|=|~=|==) *"([^"\n]+)"`)
var keywordRegx = regexp.MustCompile("^" + strings.Join(keywordSlice, "|") + "$")

func NewParam(expr string) (*Param, error) {
	p := paramRegx.FindStringSubmatch(expr)

	keyword := p[1]
	value := p[3]

	keyword = strings.ToUpper(keyword[:1]) + keyword[1:]

	if keywordRegx.MatchString(keyword) == false {
		return nil, errors.New(keyword + " keyword is unknown")
	}

	operator := ConvOperator(p[2])
	if operator == RegxEqual {
		_, err := regexp.Compile(value)
		if err != nil {
			return nil, err
		}
	}

	value = strings.ReplaceAll(value, `[quota]`, `"`)

	return &Param{
		keyword:  keyword,
		value:    value,
		operator: operator,
	}, nil
}

func ConvOperator(expr string) int {
	switch expr {
	case "!=":
		return Unequal
	case "=":
		return Equal
	case "~=":
		return RegxEqual
	case "==":
		return SuperEqual
	default:
		return 0
	}
}

func (p *Param) Match(banner *Banner) bool {
	subStr := p.value
	keyword := p.keyword

	v := reflect.ValueOf(*banner)
	str := v.FieldByName(keyword).String()

	switch p.operator {
	case Unequal:
		return !strings.Contains(str, subStr)
	case Equal:
		return strings.Contains(str, subStr)
	case RegxEqual:
		return regexp.MustCompile(subStr).MatchString(str)
	case SuperEqual:
		return str == subStr
	default:
		return false
	}
}

const (
	And = iota // &&
	Or         // ||
)

func ParseBoolFromString(expr string) (bool, error) {
	//去除空格
	expr = strings.ReplaceAll(expr, " ", "")

	//如果存在其他异常字符，则报错
	s := regexp.MustCompile(`true|false|&|\||\(|\)`).ReplaceAllString(expr, "")
	if s != "" {
		return false, errors.New(s + "is known")
	}
	return stringParse(expr)
}

func stringParse(expr string) (bool, error) {
	first := true
	operator := And

	for i := 0; i < len(expr); i++ {
		char := expr[i : i+1]
		if char == "t" {
			first = parseCoupleBool(first, true, operator)
			i += 3
		}
		if char == "f" {
			first = parseCoupleBool(first, false, operator)
			i += 4
		}
		if char == "&" {
			operator = And
			i += 1

		}
		if char == "|" {
			operator = Or
			i += 1
		}
		if char == "(" {
			length, err := findCoupleBracketIndex(expr[i:])
			if err != nil {
				return false, err
			}
			next, err := stringParse(expr[i+1 : i+length])
			if err != nil {
				return false, err
			}
			first = parseCoupleBool(first, next, operator)

			i += length
		}

	}
	return first, nil
}

func parseCoupleBool(first bool, next bool, operator int) bool {
	if operator == Or {
		return first || next
	}
	if operator == And {
		return first && next
	}
	return false
}

func findCoupleBracketIndex(expr string) (int, error) {

	var leftIndex []int
	var rightIndex []int

	for index, value := range expr {
		if value == '(' {
			leftIndex = append(leftIndex, index)
		}
		if value == ')' {
			rightIndex = append(rightIndex, index)
		}
	}

	if len(leftIndex) != len(rightIndex) {
		return 0, errors.New("bracket is not couple")
	}
	for i, index := range rightIndex {
		countLeft := strings.Count(expr[:index], "(")
		if countLeft == i+1 {
			return index, nil
		}

	}
	return 0, errors.New("bracket is not couple")
}
