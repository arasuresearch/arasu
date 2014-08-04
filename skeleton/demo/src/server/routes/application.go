package routes

import "boot"

var Routes = boot.App.Routes

func init() {
	Routes.Set("*", [][]string{
		{"*", "*", "*"},
	})
	// Routes.Set("", [][]string{
	// 	{"", "GET", "welcome.Index"},
	// })
}
