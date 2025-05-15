// vo包提供视图对象定义
package vo

// 导入错误处理包
import "errors"

// ApiCommonRes API通用响应结构
type ApiCommonRes struct {
	// 响应码，0表示成功
	Code int `json:"code"`
	// 响应消息
	Msg string `json:"msg"`
	// 响应数据
	Data interface{} `json:"data"`
}

// Error 返回响应中的错误信息
// 如果响应码不为0，返回错误消息，否则返回nil
func (this *ApiCommonRes) Error() error {
	if this.Code != 0 {
		return errors.New(this.Msg)
	}
	return nil
}
