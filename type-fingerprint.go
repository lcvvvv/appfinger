package main

import (
	"errors"
	"httpfinger/iconhash"
	"httpfinger/product"
	"regexp"
	"strings"
)

type FingerPrint struct {
	//指纹适用范围
	Protocol string
	Path     string

	//指纹适配产品
	Product *product.Product

	//指纹识别规则
	Icon    string
	Hash    string
	Keyword *Expression
}

var regxKeywordFinger = regexp.MustCompile(`^([^ ]+) (.*)$`)

func NewFingerPrint(s string) (*FingerPrint, error) {

	if regxKeywordFinger.MatchString(s) == false {
		return nil, errors.New("syntax error")
	}
	i := regxKeywordFinger.FindStringSubmatch(s)
	cms := i[1]
	expr := i[2]
	expression, err := NewExpression(expr)
	if err != nil {
		return nil, err
	}

	p := product.New(cms, "", ",", ",", ",", ",", product.Application)

	return &FingerPrint{
		//指纹适用范围
		Protocol: "",
		Path:     "",
		//指纹适配产品
		Product: p,
		//指纹识别规则
		Icon:    "",
		Hash:    "",
		Keyword: expression,
	}, nil
}

func (f *FingerPrint) Match(banner *Banner) *product.Product {
	if f.Keyword != nil {
		if f.Keyword.Match(banner) {
			return f.Product
		}
	}

	if f.Hash != "" {
		if f.MatchHash(banner.Response) {
			return f.Product
		}
	}

	if f.Icon != "" {
		if f.MatchIconHash(banner.Response) {
			return f.Product
		}
	}

	return nil
}

func (f *FingerPrint) MatchIconHash(response string) bool {
	reader := strings.NewReader(response)
	hash, _ := iconhash.Get(reader)
	return f.Icon == hash
}

func (f *FingerPrint) MatchHash(response string) bool {
	return f.Hash == MD5Encode(response)
}
