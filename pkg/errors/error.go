// nolint:errorlint
package errors

import (
	"fmt"
	"strings"
)

var _ error = &withStack{}

type withStack struct {
	Coder
	stack       []string // 状态栈
	description string   // 描述
}

func (m *withStack) Stack() []string {
	if m != nil {
		return m.stack
	}
	return nil
}

func (m *withStack) Description() string {
	if m != nil {
		return m.description
	}
	return ""
}

func (m *withStack) Detail() string {
	siList := m.StackInfo()
	callChain := ""
	lastService := ""
	for i, si := range siList {
		service := fmt.Sprintf("([%s:%s]", si.Host, si.Module)
		if i == 0 {
			callChain = fmt.Sprintf("%s<%s:%s>)", si.FuncName, si.FileName, si.Line)
		} else if service != lastService && lastService != "" {
			callChain = si.FuncName + ")->" + lastService + callChain
		} else {
			callChain = si.FuncName + "->" + callChain
		}
		lastService = service
	}
	callChain = lastService + callChain
	if m.description != "" {
		return fmt.Sprintf("stack:%s,code:%d,message:%s,desc:%s", callChain, m.Code(), m.Message(), m.description)
	}
	return fmt.Sprintf("stack:%s,code:%d,message:%s", callChain, m.Code(), m.Message())
}

type StackInfo struct {
	Host     string `json:"host"`
	PID      string `json:"pid"`
	Module   string `json:"module"`
	Code     string `json:"code"`
	FileName string `json:"file_name"`
	FuncName string `json:"func_name"`
	Line     string `json:"line"`
}

func (m *withStack) StackInfo() []StackInfo {
	siList := make([]StackInfo, 0, len(m.stack))
	for _, str := range m.stack {
		si := StackInfo{}
		slist := strings.Split(str, ",")
		for _, s := range slist {
			if strings.HasPrefix(s, "host:") {
				si.Host = strings.TrimPrefix(s, "host:")
			}
			if strings.HasPrefix(s, "pid:") {
				si.PID = strings.TrimPrefix(s, "pid:")
			}
			if strings.HasPrefix(s, "module:") {
				si.Module = strings.TrimPrefix(s, "module:")
			}
			if strings.HasPrefix(s, "code:") {
				si.Code = strings.TrimPrefix(s, "code:")
			}
			if strings.HasPrefix(s, "file:") {
				si.FileName = strings.TrimPrefix(s, "file:")
			}
			if strings.HasPrefix(s, "func:") {
				si.FuncName = strings.TrimPrefix(s, "func:")
			}
			if strings.HasPrefix(s, "line:") {
				si.Line = strings.TrimPrefix(s, "line:")
			}
		}
		siList = append(siList, si)
	}
	return siList
}

func (m withStack) Error() string {
	return fmt.Sprintf("%v:%v", m.Message()[MessageLangENKey], m.description)
}

// Wrap generate a new withStack error with `code` and `desc`.
func Wrap(code int, desc string) *withStack {
	codeMux.RLock()
	defer codeMux.RUnlock()

	errStack := &withStack{
		stack:       []string{newStack(code, Caller())},
		description: desc,
	}

	coder, exist := codes[code]
	if !exist {
		errStack.Coder = unknown
		return errStack
	}
	errStack.Coder = coder
	return errStack
}

// WrapStack generate a new withStack error with `code` , `desc`,`stack`.
func WrapStack(code int, desc string, stack []string) *withStack {
	codeMux.RLock()
	defer codeMux.RUnlock()

	stack = append(stack, newStack(code, Caller()))
	errStack := &withStack{
		stack:       stack,
		description: desc,
	}

	coder, exist := codes[code]
	if !exist {
		errStack.Coder = unknown
		return errStack
	}
	errStack.Coder = coder
	return errStack
}

// WrapF generate a new withStack error with `code` and desc `format+arg...`.
func WrapF(code int, format string, args ...interface{}) *withStack {
	codeMux.RLock()
	defer codeMux.RUnlock()

	errStack := &withStack{
		stack:       []string{newStack(code, Caller())},
		description: fmt.Sprintf(format, args...),
	}

	coder, exist := codes[code]
	if !exist {
		errStack.Coder = unknown
		return errStack
	}
	errStack.Coder = coder
	return errStack
}

// WrapError generate a new withStack error with `code` and desc `format+arg...`
// if err not nil, inherit its stack and error message, replace origin code.
func WrapError(code int, err error) *withStack {
	codeMux.RLock()
	defer codeMux.RUnlock()

	errStack := &withStack{
		stack:       []string{newStack(code, Caller())},
		description: "",
	}

	coder, exist := codes[code]
	if !exist {
		errStack.Coder = unknown
		if err != nil {
			errStack.description = err.Error()
		}
		return errStack
	}

	if err != nil {
		if st, ok := err.(*withStack); ok {
			errStack.stack = append(st.Stack(), errStack.stack...)
			errStack.description = st.description
		} else {
			errStack.description = err.Error()
		}
	}
	errStack.Coder = coder
	return errStack
}

// UpdateStack add a new layer to err's caller stack.
func UpdateStack(err error) *withStack {
	errStack := FromError(err)
	if errStack != nil {
		errStack.stack = append(errStack.stack, newStack(errStack.Code(), Caller()))
	}
	return errStack
}

func (m withStack) ToBasicJson() map[string]interface{} {
	out := make(map[string]interface{})
	out["desc"] = m.description
	out["message"] = m.Message()
	out["code"] = m.Code()

	return out
}

func (m withStack) ToDetailJson() map[string]interface{} {
	out := make(map[string]interface{})
	out["desc"] = m.description
	out["stack"] = m.StackInfo()
	out["message"] = m.Message()
	out["code"] = m.Code()
	out["http"] = m.HTTPStatus()
	return out
}

// ParseError parse any error into *withStack.
// nil error will return nil direct.
// None withStack error will be parsed as ErrUnknown.
func FromError(err error) *withStack {
	if err == nil {
		return nil
	}

	if v, ok := err.(*withStack); ok {
		return v
	}

	return &withStack{
		Coder:       unknown,
		stack:       []string{newStack(unknown.code, Caller())},
		description: err.Error(),
	}
}

func newStack(code int, caller string) string {
	return currentModule.String() + fmt.Sprintf(",code:%d,", code) + caller
}
