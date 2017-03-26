package utils

import (
	"github.com/mozillazg/go-pinyin"
	"regexp"
	"strings"
)

func ConvertPY(str string) string {
	c := ""
	match, _ := regexp.Match("[a-zA-Z]", []byte{str[0]})
	if match {
		c = string(str[0])
	} else {
		args := pinyin.NewArgs()
		names := pinyin.LazyPinyin(str, args)
		c = string(names[0][0])
	}

	return strings.ToUpper(c)
}
