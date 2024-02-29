package cli

import (
	"embed"
	"path"

    "github.com/nmusey/letsgo/internal/utils"
)

type App struct {
	Name  string
    Repo string
    Root string	
    Replacements map[string]string
    Filesystem embed.FS
}

//go:embed _templates/*
var filesystem embed.FS

func NewApp(name string, repo string, root string) error {
	app := App{
		Name:  name,
        Repo: repo,
        Root: path.Join(root, name),
        Replacements: map[string]string{
            "$appName": name,
            "$appRepo": repo,
            "$dbPort":  "5432",
	    },
        Filesystem: filesystem,
	}

    return app.initializeDirectory("")
}

func (app App) initializeDirectory(directory string) error {
    utils.UpsertFolder(path.Join(app.Root, directory))
    files := app.readDirectory(directory)
    if err := app.copyFiles(files, directory); err != nil {
        return err
    }
    
    return nil
}

func (app App) readDirectory(directory string) map[string]string {
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

func (app App) isEmbeddedDirectory(filepath string) bool {
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

func (app *App) copyFiles(files map[string]string, directory string) error {
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

func (app *App) copyFile(filepath string, contents string) error {
	outpath := path.Join(app.Root, filepath)
	return utils.CopyTemplateFile(contents, outpath, app.Replacements)
}

func (app *App) getTemplatePath(filepath string) string {
    // _templates is copied when running make build, so it can sit in the root directory but still run with the limitations of the embed package
    return path.Join("_templates", filepath)
}
