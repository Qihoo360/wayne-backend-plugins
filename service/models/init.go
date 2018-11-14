package models

import (
	"github.com/astaxie/beego/orm"
)

var (
	ServiceModel    *serviceModel
	ServiceTplModel *serviceTplModel
)

func init() {
	orm.RegisterModel()

	ServiceModel = &serviceModel{}
	ServiceTplModel = &serviceTplModel{}
}
