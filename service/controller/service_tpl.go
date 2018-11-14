package controller

import (
	"encoding/json"
	"fmt"

	"github.com/Qihoo360/wayne/src/backend/controllers/base"
	"github.com/Qihoo360/wayne/src/backend/models"
	svcmodel "github.com/Qihoo360/wayne/src/backend/plugins/service/models"
	"github.com/Qihoo360/wayne/src/backend/util/hack"
	"github.com/Qihoo360/wayne/src/backend/util/logs"
	"k8s.io/api/core/v1"
)

// 服务模版相关操作
type ServiceTplController struct {
	base.APIController
}

func (c *ServiceTplController) URLMapping() {
	c.Mapping("List", c.List)
	c.Mapping("Create", c.Create)
	c.Mapping("Get", c.Get)
	c.Mapping("Update", c.Update)
	c.Mapping("Delete", c.Delete)
}

func (c *ServiceTplController) Prepare() {
	// Check administration
	c.APIController.Prepare()
	// Check permission
	perAction := ""
	_, method := c.GetControllerAndAction()
	switch method {
	case "Get", "List":
		perAction = models.PermissionRead
	case "Create":
		perAction = models.PermissionCreate
	case "Update":
		perAction = models.PermissionUpdate
	case "Delete":
		perAction = models.PermissionDelete
	}
	if perAction != "" {
		c.CheckPermission(models.PermissionTypeService, perAction)
	}
}

// @Title GetAll
// @Description get all ServiceTpl
// @Param	pageNo		query 	int	false		"the page current no"
// @Param	pageSize		query 	int	false		"the page size"
// @Param	name		query 	string	false		"name filter"
// @Param	deleted		query 	bool	false		"is deleted, default list all"
// @Success 200 {object} []models.ServiceTemplate success
// @router / [get]
func (c *ServiceTplController) List() {
	param := c.BuildQueryParam()
	name := c.Input().Get("name")
	if name != "" {
		param.Query["name__contains"] = name
	}

	isOnline := c.GetIsOnlineFromQuery()

	serviceId := c.Input().Get("serviceId")
	if serviceId != "" {
		param.Query["service_id"] = serviceId
	}

	var serviceTpls []models.ServiceTemplate
	total, err := models.ListTemplate(&serviceTpls, param, models.TableNameServiceTemplate, models.PublishTypeService, isOnline)
	if err != nil {
		logs.Error("list by param (%v) error. %v", param, err)
		c.HandleError(err)
		return
	}
	for index, tpl := range serviceTpls {
		serviceTpls[index].ServiceId = tpl.Service.Id
	}

	c.Success(param.NewPage(total, serviceTpls))
}

// @Title Create
// @Description create ServiceTpl
// @Param	body		body 	models.ServiceTemplate	true		"The ServiceTpl content"
// @Success 200 return models.ServiceTemplate success
// @router / [post]
func (c *ServiceTplController) Create() {
	var serviceTpl models.ServiceTemplate
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &serviceTpl)
	if err != nil {
		logs.Error("get body error. %v", err)
		c.AbortBadRequestFormat("ServiceTemplate")
	}
	err = validServiceTemplate(serviceTpl.Template)
	if err != nil {
		logs.Error("valid template err %v", err)
		c.AbortBadRequestFormat("KubeService")
	}

	serviceTpl.User = c.User.Name

	_, err = svcmodel.ServiceTplModel.Add(&serviceTpl)
	if err != nil {
		logs.Error("create error.%v", err.Error())
		c.HandleError(err)
		return
	}
	c.Success(serviceTpl)
}

func validServiceTemplate(serviceTplStr string) error {
	service := v1.Service{}
	err := json.Unmarshal(hack.Slice(serviceTplStr), &service)
	if err != nil {
		return fmt.Errorf("service template format error.%v", err.Error())
	}
	return nil
}

// @Title Get
// @Description find Object by id
// @Param	id		path 	int	true		"the id you want to get"
// @Success 200 {object} models.ServiceTemplate success
// @router /:id([0-9]+) [get]
func (c *ServiceTplController) Get() {
	id := c.GetIDFromURL()

	serviceTpl, err := svcmodel.ServiceTplModel.GetById(id)
	if err != nil {
		logs.Error("get template error %v", err)
		c.HandleError(err)
		return
	}

	c.Success(serviceTpl)
}

// @Title Update
// @Description update the ServiceTpl
// @Param	id		path 	int	true		"The id you want to update"
// @Param	body		body 	models.ServiceTemplate	true		"The body"
// @Success 200 models.ServiceTemplate success
// @router /:id([0-9]+) [put]
func (c *ServiceTplController) Update() {
	id := c.GetIDFromURL()
	var serviceTpl models.ServiceTemplate
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &serviceTpl)
	if err != nil {
		logs.Error("Invalid param body.%v", err)
		c.AbortBadRequestFormat("ServiceTemplate")
	}
	if err = validServiceTemplate(serviceTpl.Template); err != nil {
		logs.Error("valid template err %v", err)
		c.AbortBadRequestFormat("KubeService")
	}

	serviceTpl.Id = int64(id)
	err = svcmodel.ServiceTplModel.UpdateById(&serviceTpl)
	if err != nil {
		logs.Error("update error.%v", err)
		c.HandleError(err)
		return
	}
	c.Success(serviceTpl)
}

// @Title Delete
// @Description delete the ServiceTpl
// @Param	id		path 	int	true		"The id you want to delete"
// @Param	logical		query 	bool	false		"is logical deletion,default true"
// @Success 200 {string} delete success!
// @router /:id([0-9]+) [delete]
func (c *ServiceTplController) Delete() {
	id := c.GetIDFromURL()
	logical := c.GetLogicalFromQuery()

	err := svcmodel.ServiceTplModel.DeleteById(int64(id), logical)
	if err != nil {
		logs.Error("delete %d error.%v", id, err)
		c.HandleError(err)
		return
	}
	c.Success(nil)
}
