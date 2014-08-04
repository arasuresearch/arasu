package debug

import "config"

func init() {
	config.Config["debug"] = map[string]string{
		"secret":                "random",
		"Secret.Key":            "bPlNFGdSC2wd8f2QnFhk5A84JJjKWZdKH9H2FHFuvUs9Jz8UvBHv3Vc5awx39ivu",
		"Ssl":                   "enable",
		"Ssl.Cert":              "",
		"Ssl.Key":               "",
		"Cookie":                "true",
		"Cookie.Prefix":         "",
		"rdbms":                 "mysql,root:@tcp(127.0.0.1:3306)/{{.Name}}_development?charset=utf8&parseTime=True",
		"bigdata":               "hbase,root:@tcp(127.0.0.1:9090)/{{.Name}}_development?charset=utf8&parseTime=True",
		"header":                "true",
		"header.AllowAuthToken": "true",
		"Port":                  "4000",
		"Mode":                  "debug",
		"Format.Default":        "json",
		"Format.Allowed":        "*",
	}
}
