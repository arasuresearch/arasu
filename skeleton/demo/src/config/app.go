package config

var Config = map[string]map[string]string{
	"": {
		"Name":                  "{{.Name}}",
		"Description":           "A {{.Name}} Application",
		"ReCompileAtNewRequest": "server,conf",
		"secret":                "random",
		"Secret.Key":            "bPlNFGdSC2wd8f2QnFhk5A84JJjKWZdKH9H2FHFuvUs9Jz8UvBHv3Vc5awx39ivu",
		"Ssl":                   "enable",
		"Ssl.Cert":              "",
		"Ssl.Key":               "",
		"Cookie":                "true",
		"Cookie.Prefix":         "",
		"header":                "true",
		"header.AllowAuthToken": "true",
		"Port":                  "4000",
		"Mode":                  "debug",
		"Dstore":                "{{.Dstore}}",
		"Database":              "{{.Database}}",
		"Format.Default":        "json",
		"Format.Allowed":        "*",
	},
}

// "rdbms":                "mysql,root:@tcp(127.0.0.1:3306)/{{.Name}}_development?charset=utf8&parseTime=True",
// "bigdata":              "hbase,root:@tcp(127.0.0.1:9090)/{{.Name}}_development?charset=utf8&parseTime=True",
// "rdbms/mysql":           "mysql,root:@tcp(127.0.0.1:3306)/demo_development?charset=utf8&parseTime=True",
// "bigdata/hbase":         "hbase,root:@tcp(127.0.0.1:9090)/demo_development?charset=utf8&parseTime=True",
