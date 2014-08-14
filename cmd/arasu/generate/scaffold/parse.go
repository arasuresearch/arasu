package scaffold

import (
	"encoding/json"
	"github.com/arasuresearch/arasu/lib/stringer"
	"log"
	"strings"
)

var bigdataAttrTypes = []string{"int", "bool", "String", "DateTime", "num"}

func parseBigDataAttrs(a []string) (map[string]string, interface{}, interface{}, string) {
	//[Contact Post:Fn Post:Ln Profile:Name:String Profile:Age:int Profile:Dob:DateTime]
	//fmt.Println(c.Attrs, c.ClientAttrs, c.ClientModelViewAttrs, c.ClientModelMetadata)
	var changed bool
	for i, e := range a {
		changed = false
		if strings.Contains(e, "{") {
			changed = true
			e = strings.Replace(e, "{", "", -1)
		}
		if strings.Contains(e, "}") {
			changed = true
			e = strings.Replace(e, "}", "", -1)
		}
		if changed {
			a[i] = e
		}
	}
	attrs := map[string]string{}
	var ls string
	for _, e := range a {
		if strings.Contains(e, ":") {
			la := strings.SplitN(e, ":", 2)
			ls = la[0]
		} else {
			ls = e
		}
		attrs[stringer.Camelize(ls)] = "ColumnFamily"
	}

	cattrs := []string{"Id"}
	for _, e := range a {
		if strings.Contains(e, ":") {
			la := strings.Split(e, ":")
			cattrs = append(cattrs, stringer.Camelize(la[0])+"."+stringer.Camelize(la[1]))
		}
	}

	cmva := make(map[string]interface{})
	cmva["Id"] = "String"
	for _, e := range a {
		a := strings.Split(e, ":")
		switch len(a) {
		case 1:

		case 2:
			cf, ci, ct := stringer.Camelize(a[0]), stringer.Camelize(a[1]), "String"
			if cols, ok := cmva[cf]; ok {
				cols.(map[string]string)[ci] = ct
			} else {
				cmva[cf] = map[string]string{ci: ct}
			}
		default:
			cf, ci, ct := stringer.Camelize(a[0]), stringer.Camelize(a[1]), a[2]

			if cols, ok := cmva[cf]; ok {
				cols.(map[string]string)[ci] = ct
			} else {
				cmva[cf] = map[string]string{ci: ct}
			}
		}
	}
	for _, v := range cmva {
		if cfa, ok := v.(map[string]string); ok {
			for _, v := range cfa {
				if !stringer.Contains(bigdataAttrTypes, v) {
					log.Fatalf("%s type is not supported", v)
				}
			}
		}
	}
	data, _ := json.Marshal(cmva)
	delete(cmva, "Id")
	return attrs, cattrs, cmva, string(data)
}

//var bigdataAttrTypes = []string{"int", "bool", "String", "DateTime", "num"}

var rdbmsAttrTypes = []string{"int", "bool", "String", "DateTime", "num",
	"integer", "text", "date", "boolean", "string", "timestamp"}

func rdbmsAttrsRealType(typ string) (rtyp string) {
	rtyp = typ
	switch typ {
	case "int":
		rtyp = "integer"
	case "bool":
		rtyp = "boolean"
	case "DateTime":
		rtyp = "timestamp"
	case "num":
		rtyp = "float"
	}
	return
}
func rdbmsClientAttrsRealType(typ string) (rtyp string) {
	rtyp = typ
	switch typ {
	case "string", "text":
		rtyp = "String"
	case "integer":
		rtyp = "int"
	case "boolean":
		rtyp = "bool"
	case "timestamp", "date":
		rtyp = "DateTime"
	case "float", "decimal":
		rtyp = "num"
	}
	return
}

func parseRdbmsAttrs(a []string) (map[string]string, interface{}, interface{}, string) {
	attrs := map[string]string{}
	cmva := map[string]string{"Id": "int"}
	cols := []string{"Id"}
	var col string

	for _, e := range a {
		if strings.Contains(e, ":") {
			la := strings.SplitN(e, ":", 2)
			if !stringer.Contains(rdbmsAttrTypes, la[1]) {
				log.Fatalf("%s type is not supported", la[1])
			}
			attrs[la[0]] = stringer.Camelize(rdbmsAttrsRealType(la[1]))
			cmva[stringer.Camelize(la[0])] = stringer.Camelize(rdbmsClientAttrsRealType(la[1]))
			col = stringer.Camelize(la[0])
		} else {
			attrs[e] = "String"
			cmva[stringer.Camelize(e)] = "String"
			col = stringer.Camelize(e)
		}
		cols = append(cols, col)
	}

	//adding extra fields
	cmva["CreatedAt"] = "DateTime"
	cmva["UpdatedAt"] = "DateTime"
	cols = append(cols, "CreatedAt")
	cols = append(cols, "UpdatedAt")

	data, _ := json.Marshal(cmva)

	delete(cmva, "Id")
	delete(cmva, "CreatedAt")
	delete(cmva, "UpdatedAt")
	return attrs, cols, cmva, string(data)
}

//arasubuild;arasu new d0 -d mysql -ds rdbms

//arasu generate scaffold User name:string pass age:integer dob:timestamp
