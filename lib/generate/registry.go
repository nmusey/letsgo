package generate

type Kind struct {
	Name string
	TemplateFile string
	OutputFile string
	FlagLong  string
	FlagShort string
}

var Kinds = []Kind{
	{Name: "model", TemplateFile: "model.go.template", OutputFile: "model.go", FlagLong: "model", FlagShort: "m"},
	{Name: "repository", TemplateFile: "repository.go.template", OutputFile: "repository.go", FlagLong: "repository", FlagShort: "r"},
	{Name: "service", TemplateFile: "service.go.template", OutputFile: "service.go", FlagLong: "service", FlagShort: "s"},
}

func KindByName(name string) (Kind, bool) {
	for _, kind := range Kinds {
		if kind.Name == name {
			return kind, true
		}
	}

	return Kind{}, false
}
