package controller

import (
	"encoding/json"

	"github.com/Qihoo360/wayne/src/backend/controllers/base"
	"github.com/Qihoo360/wayne/src/backend/models"
	svcmodel "github.com/Qihoo360/wayne/src/backend/plugins/service/models"
	"github.com/Qihoo360/wayne/src/backend/util/logs"
)

type ServiceController struct {
	base.APIController
}

func (c *ServiceController) URLMapping() {
	c.Mapping("GetNames", c.GetNames)
	c.Mapping("List", c.List)
	c.Mapping("Create", c.Create)
	c.Mapping("Get", c.Get)
	c.Mapping("Update", c.Update)
	c.Mapping("Delete", c.Delete)
}

func (c *ServiceController) Prepare() {
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

// @Title List/
// @Description get all id and names
// @Param	appId		query 	int	false		"the app id"
// @Param	deleted		query 	bool	false		"is deleted,default false."
// @Success 200 {object} []models.Service success
// @router /names [get]
func (c *ServiceController) GetNames() {
	filters := make(map[string]interface{})
	deleted := c.GetDeleteFromQuery()

	filters["Deleted"] = deleted
	if c.AppId != 0 {
		filters["App__Id"] = c.AppId
	}

	services, err := svcmodel.ServiceModel.GetNames(filters)
	if err != nil {
		logs.Error("get names error. %v, delete-status %v", err, deleted)
		c.HandleError(err)
		return
	}

	c.Success(services)
}

// @Title GetAll
// @Description get all Service
// @Param	pageNo		query 	int	false		"the page current no"
// @Param	pageSize		query 	int	false		"the page size"
// @Param	name		query 	string	false		"name filter"
// @Param	deleted		query 	bool	false		"is deleted, default list all"
// @Success 200 {object} []models.Service success
// @router / [get]
func (c *ServiceController) List() {
	param := c.BuildQueryParam()
	name := c.Input().Get("name")
	if name != "" {
		param.Query["name__contains"] = name
	}

	service := []models.Service{}
	if c.AppId != 0 {
		param.Query["App__Id"] = c.AppId
	} else if !c.User.Admin {
		param.Query["App__AppUsers__User__Id__exact"] = c.User.Id
		perName := models.PermissionModel.MergeName(models.PermissionTypeService, models.PermissionRead)
		param.Query["App__AppUsers__Group__Permissions__Permission__Name__contains"] = perName
		param.Groupby = []string{"Id"}
	}

	total, err := models.GetTotal(new(models.Service), param)
	if err != nil {
		logs.Error("get total count by param (%s) error. %v", param, err)
		c.HandleError(err)
		return
	}

	err = models.GetAll(new(models.Service), &service, param)
	if err != nil {
		logs.Error("list by param (%s) error. %v", param, err)
		c.HandleError(err)
		return
	}
	for key, one := range service {
		service[key].AppId = one.App.Id
	}

	c.Success(param.NewPage(total, service))
}

// @Title Create
// @Description create Service
// @Param	body		body 	models.Service	true		"The Service content"
// @Success 200 return models.Service success
// @router / [post]
func (c *ServiceController) Create() {
	var service models.Service
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &service)
	if err != nil {
		logs.Error("get body error. %v", err)
		c.AbortBadRequestFormat("Service")
	}

	service.User = c.User.Name
	_, err = svcmodel.ServiceModel.Add(&service)

	if err != nil {
		logs.Error("create error.%v", err.Error())
		c.HandleError(err)
		return
	}
	c.Success(service)
}

// @Title Get
// @Description find Object by id
// @Param	id		path 	int	true		"the id you want to get"
// @Success 200 {object} models.Service success
// @router /:id([0-9]+) [get]
func (c *ServiceController) Get() {
	id := c.GetIDFromURL()

	service, err := svcmodel.ServiceModel.GetById(int64(id))
	if err != nil {
		logs.Error("get by id (%d) error.%v", id, err)
		c.HandleError(err)
		return
	}

	c.Success(service)
	return
}

// @Title Update
// @Description update the Service
// @Param	id		path 	int	true		"The id you want to update"
// @Param	body		body 	models.Service	true		"The body"
// @Success 200 models.Service success
// @router /:id([0-9]+) [put]
func (c *ServiceController) Update() {
	id := c.GetIDFromURL()
	var service models.Service
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &service)
	if err != nil {
		logs.Error("Invalid param body.%v", err)
		c.AbortBadRequestFormat("Service")
	}

	service.Id = int64(id)
	err = svcmodel.ServiceModel.UpdateById(&service)
	if err != nil {
		logs.Error("update error.%v", err)
		c.HandleError(err)
		return
	}
	c.Success(service)
}

// @Title UpdateOrders
// @Description batch update the Orders
// @Param	body		body 	[]models.Service	true		"The body"
// @Success 200 models.Deployment success
// @router /updateorders [put]
func (c *ServiceController) UpdateOrders() {
	var services []*models.Service
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &services)
	if err != nil {
		logs.Error("Invalid param body.%v", err)
		c.AbortBadRequestFormat("services")
	}

	err = svcmodel.ServiceModel.UpdateOrders(services)
	if err != nil {
		logs.Error("update orders (%v) error.%v", services, err)
		c.HandleError(err)
		return
	}
	c.Success("ok!")
}

// @Title Delete
// @Description delete the Service
// @Param	id		path 	int	true		"The id you want to delete"
// @Param	logical		query 	bool	false		"is logical deletion,default true"
// @Success 200 {string} delete success!
// @router /:id([0-9]+) [delete]
func (c *ServiceController) Delete() {
	id := c.GetIDFromURL()

	logical := c.GetLogicalFromQuery()

	err := svcmodel.ServiceModel.DeleteById(int64(id), logical)
	if err != nil {
		logs.Error("delete %d error.%v", id, err)
		c.HandleError(err)
		return
	}
	c.Success(nil)
}
