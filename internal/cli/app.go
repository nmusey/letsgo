package cli

import (
	"embed"
	"os"
	"path"

	"github.com/nmusey/letsgo/internal/utils"
	"golang.org/x/mod/modfile"
)

type Template struct {
    Root string	
    TemplateName string
    Replacements map[string]string
    Filesystem embed.FS
}

//go:embed _templates/app/*
var appFilesystem embed.FS

func NewAppTemplate(name string, repo string, root string) error {
	template := Template{
        Root: path.Join(root, name),
        TemplateName: "app",
        Replacements: map[string]string{
            "$appName": name,
            "$appRepo": repo,
            "$dbPort":  "5432",
	    },
        Filesystem: appFilesystem,
	}

    return template.initializeDirectory("")
}

//go:embed _templates/package/*
var packageFilesystem embed.FS
func NewPackageTemplate(packageName string, root string) error {
    template := Template {
        Root: root,
        TemplateName: "package",
        Replacements: map[string]string{
            "$package": packageName,
        },
        Filesystem: packageFilesystem,
    }

    template.Replacements["$appRepo"] = template.getGoModuleName()
    return template.initializeDirectory("")
}

func (app Template) initializeDirectory(directory string) error {
    replacedDirectory := utils.ReplaceAllInString(directory, app.Replacements)
    utils.UpsertFolder(path.Join(app.Root, replacedDirectory))
    files := app.readDirectory(directory)
    if err := app.copyFiles(files); err != nil {
        return err
    }
    
    return nil
}

func (app Template) readDirectory(directory string) map[string]string {
    files, err := app.Filesystem.ReadDir(app.getTemplatePath(directory))
    if err != nil {
        panic(err)
    }

    templates := make(map[string]string)
    for _, file := range files {
        filepath := path.Join(directory, file.Name())

        if app.isEmbeddedDirectory(filepath) {
            app.initializeDirectory(filepath)
            continue
        }

        contents, _ := app.Filesystem.ReadFile(app.getTemplatePath(filepath))
        templates[filepath] = string(contents)
    }

    return templates
}

func (app Template) isEmbeddedDirectory(filepath string) bool {
    fsFile, err := app.Filesystem.Open(app.getTemplatePath(filepath))
    if err != nil {
        return false
    }

    fileStats, err := fsFile.Stat()
    if err != nil {
        return false
    }

    return fileStats.IsDir()
}

func (app *Template) copyFiles(files map[string]string) error {
	for filepath, contents := range files {
        // We need to treat go.mod specially because it can't be embedded directly
        if path.Base(filepath) == "go_mod" {
            filepath = path.Join(path.Dir(filepath), "go.mod")
        }

		if err := app.copyFile(filepath, contents); err != nil {
			return err
		}
	}

	return nil
}

func (app *Template) copyFile(filepath string, contents string) error {
    replacedPath := utils.ReplaceAllInString(filepath, app.Replacements)
	outpath := path.Join(app.Root, replacedPath)
	return utils.CopyTemplateFile(contents, outpath, app.Replacements)
}

func (app *Template) getTemplatePath(filepath string) string {
    // _templates is copied when running make build, so it can sit in the root directory but still run with the limitations of the embed package
    return path.Join("_templates", app.TemplateName, filepath)
}

func (app *Template) getGoModuleName() string {
    filepath := "./go.mod"
    contents, err := os.ReadFile(filepath)
    if err != nil {
        panic("Unable to read go.mod file: " + err.Error())
    }

    modfile, err := modfile.Parse("./go.mod", contents, nil) 
    return modfile.Module.Mod.Path
}
