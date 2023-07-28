package log_test

import (
	"testing"

	"github.com/wangweihong/eazycloud/pkg/log"
)

func TestEvery(t *testing.T) {
	//{"slice": "[]string{\"a lll\"}"}
	log.Info("", log.Every("slice", []string{"a lll"}))
	//{"map": "map[string]interface {}{\"aaa\":123, \"name\":false, \"so\":\"bb\"}"}
	log.Info("", log.Every("map", map[string]interface{}{"aaa": 123, "so": "bb", "name": false}))
	type Person struct {
		Emails []string `mapstructure:"emails"`
	}
	var result Person
	result.Emails = []string{"a", "b"}

	var result2 Person
	result2.Emails = []string{"a", "b"}
	// {"object": "log_test.Person{Emails:[]string{\"a\", \"b\"}}"}
	log.Info("", log.Every("object", result))
	// {"objectP":"&log_test.Person{Emails:[]string{\"a\", \"b\"}}"}
	log.Info("", log.Every("objectP", &result))
	// {"objectP":"[]log_test.Person{log_test.Person{Emails:[]string{\"a\", \"b\"}},
	// log_test.Person{Emails:[]string{\"a\", \"b\"}}}"}
	log.Info("", log.Every("objectP", []Person{result, result2}))
}
