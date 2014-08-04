package {{.Name}}s

import (
	. "github.com/arasuresearch/arasu/lib"
	"server/controllers/application"
	. "server/dstores/rdbms"
)

type {{.Cname}}sController struct {
	application.ApplicationController
}

func (c *{{.Cname}}sController) Index() ([]{{.Cname}}, error) {
	var all_{{.Name}}s []{{.Cname}}
	Mysql.Find(&all_{{.Name}}s)
	return all_{{.Name}}s, nil
}

func (c *{{.Cname}}sController) Show(id int64) ({{.Name}} *{{.Cname}}, err error) {
	{{.Name}} = new({{.Cname}})
	return {{.Name}}, Mysql.Where(id).Find({{.Name}}).Error
}

func (c *{{.Cname}}sController) Create(attr map[string]interface{}) ({{.Name}} *{{.Cname}}, err error) {
	{{.Name}} = new({{.Cname}})
	if err = MapToObj(attr, {{.Name}}, false); err != nil {
		return
	}
	err = Mysql.Save({{.Name}}).Error
	return
}

func (c *{{.Cname}}sController) Update(id int64, attr map[string]interface{}) ({{.Name}} *{{.Cname}}, err error) {
	{{.Name}} = new({{.Cname}})
	Mysql.Where(id).Find({{.Name}})
	if err = MapToObj(attr, {{.Name}}, false); err != nil {
		return
	}
	return {{.Name}}, Mysql.Save({{.Name}}).Error
}

func (c *{{.Cname}}sController) Destroy(id int64) (err error) {
	{{.Name}} := new({{.Cname}})
	return Mysql.Where(id).Delete({{.Name}}).Error
}
