package {{.Name}}s

import (
	. "github.com/arasuresearch/arasu/lib"
	"server/controllers/application"
	. "server/dstores/bigdata"
)

type {{.Cname}}sController struct {
	application.ApplicationController
}

func (c *{{.Cname}}sController) Index() ([]{{.Cname}}, error) {
	var all_{{.Name}}s []{{.Cname}}
	err := BigData.GetRows(&all_{{.Name}}s)
	return all_{{.Name}}s,err
}

func (c *{{.Cname}}sController) Show(id string) ({{.Name}} *{{.Cname}}, err error) {
	{{.Name}}.Id = id
	BigData.GetRow(&{{.Name}})
	return
}

func (c *{{.Cname}}sController) Create(attr map[string]interface{}) ({{.Name}} *{{.Cname}}, err error) {
	{{.Name}} = new({{.Cname}})
	if err = MapToBdObj(attr, {{.Name}}); err != nil {
		return
	}
	if err = BigData.Save({{.Name}});err != nil {
		return
	}
	return
}

func (c *{{.Cname}}sController) Update(id string, attr map[string]interface{}) ({{.Name}} *{{.Cname}}, err error) {
	{{.Name}} = new({{.Cname}})
	{{.Name}}.Id = id
	if err = MapToBdObj(attr, {{.Name}}); err != nil {
		return
	}
	if err = BigData.Save({{.Name}});err != nil {
		return
	}
	BigData.GetRow({{.Name}})
	return 
}

func (c *{{.Cname}}sController) Destroy(id string) (err error) {
	return BigData.DeleteRow(&{{.Cname}}{Id: id})
}
