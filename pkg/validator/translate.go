package validator

import "fmt"

type Translator interface {
	ZH() string
	EN() string
}

type NameTranslator struct{}

func (t NameTranslator) ZH() string {
	return fmt.Sprintf("必须匹配正则表达式 %s", nameRegixPattern)
}

func (t NameTranslator) EN() string {
	return fmt.Sprintf("name must match pattern:%s", nameRegixPattern)
}

type DescriptionTranslator struct{}

func (t DescriptionTranslator) ZH() string {
	return fmt.Sprintf("长度必须小于%d个字符", maxDescriptionLength)
}

func (t DescriptionTranslator) EN() string {
	return fmt.Sprintf("must be less than %d characters", maxDescriptionLength)
}

type IDTranslator struct{}

func (t IDTranslator) ZH() string {
	return fmt.Sprintf("未指定对象")
}

func (t IDTranslator) EN() string {
	return fmt.Sprintf("id missing")
}

type URLInvaliTranslator struct{}

func (t URLInvaliTranslator) ZH() string {
	return fmt.Sprintf("地址不正确")
}

func (t URLInvaliTranslator) EN() string {
	return fmt.Sprintf("wrong http url")
}
