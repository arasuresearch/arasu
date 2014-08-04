package routes

func init() {
	Routes.Set("{{.Name}}s", [][]string{
		{"", "GET", "{{.Name}}s.Index"},
		{"", "POST", "{{.Name}}s.Create"},
		{"new", "GET", "{{.Name}}s.New"},
		{":id/edit", "GET", "{{.Name}}s.Edit"},
		{":id", "GET", "{{.Name}}s.Show"},
		{":id", "PATCH", "{{.Name}}s.Update"},
		{":id", "PUT", "{{.Name}}s.Update"},
		{":id", "DELETE", "{{.Name}}s.Destroy"},
	})
}
