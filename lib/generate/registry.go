package generate

type Kind struct {
	Name         string
	TemplateFile string
	OutputPath string
	FlagLong   string
	FlagShort  string
}

var Kinds = []Kind{
	{Name: "model", TemplateFile: "model.go.template", OutputPath: "lib/domain/{{.Package}}/model.go", FlagLong: "model", FlagShort: "m"},
	{Name: "service", TemplateFile: "service.go.template", OutputPath: "lib/domain/{{.Package}}/service.go", FlagLong: "service", FlagShort: "s"},
	{Name: "app-service", TemplateFile: "appservice.go.template", OutputPath: "lib/application/{{.Package}}/{{.Package}}service.go", FlagLong: "app-service", FlagShort: "a"},
	{Name: "repository", TemplateFile: "repository.go.template", OutputPath: "lib/infrastructure/database/{{.Package}}/repository.go", FlagLong: "repository", FlagShort: "r"},
}

func KindByName(name string) (Kind, bool) {
	for _, kind := range Kinds {
		if kind.Name == name {
			return kind, true
		}
	}

	return Kind{}, false
}
