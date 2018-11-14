package models

import (
	. "github.com/Qihoo360/wayne/src/backend/models"
)

type serviceTplModel struct{}

func (*serviceTplModel) Add(m *ServiceTemplate) (id int64, err error) {
	m.Service = &Service{Id: m.ServiceId}
	id, err = Ormer().Insert(m)
	return
}

func (*serviceTplModel) UpdateById(m *ServiceTemplate) (err error) {
	v := ServiceTemplate{Id: m.Id}
	// ascertain id exists in the database
	if err = Ormer().Read(&v); err == nil {
		m.Service = &Service{Id: m.ServiceId}
		_, err = Ormer().Update(m)
		return err
	}
	return
}

func (*serviceTplModel) GetById(id int64) (v *ServiceTemplate, err error) {
	v = &ServiceTemplate{Id: id}

	if err = Ormer().Read(v); err == nil {
		_, err = Ormer().LoadRelated(v, "Service")
		if err == nil {
			v.ServiceId = v.Service.Id
			return v, nil
		}
	}
	return nil, err
}

func (*serviceTplModel) DeleteById(id int64, logical bool) (err error) {
	v := ServiceTemplate{Id: id}
	// ascertain id exists in the database
	if err = Ormer().Read(&v); err == nil {
		if logical {
			v.Deleted = true
			_, err = Ormer().Update(&v)
			return err
		}
		_, err = Ormer().Delete(&v)
		return err
	}
	return
}
