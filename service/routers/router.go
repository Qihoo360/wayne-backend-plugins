package routers

import (
	"github.com/Qihoo360/wayne/src/backend/plugins/service/controller"
	"github.com/astaxie/beego"
)

func init() {
	nsWithApp := beego.NewNamespace("/api/v1",
		beego.NSNamespace("/apps/:appid([0-9]+)/services",
			beego.NSInclude(
				&controller.ServiceController{},
			)),
		beego.NSNamespace("/apps/:appid([0-9]+)/services/tpls",
			beego.NSInclude(
				&controller.ServiceTplController{},
			)),
	)

	beego.AddNamespace(nsWithApp)
}
