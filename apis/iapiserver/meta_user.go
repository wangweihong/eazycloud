package iapiserver

import "github.com/wangweihong/eazycloud/apis/imachinery"

// +k8s:deepcopy-gen=true
type User struct {
	imachinery.ObjectMeta
	Password string `json:"password,omitempty"`
}

func (u *User) Transfer() *User {
	u.Password = ""
	return nil
}
