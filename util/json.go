package util

import (
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"io/ioutil"
)

func LoadJsonAndParse(f string, res interface{}) error {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return errors.Wrap(err, "加载插件配置文件失败")
	}
	err = json.Unmarshal(b, &res)
	if err != nil {
		return errors.Wrap(err, "解析插件配置文件失败")
	}
	return nil
}
