package cmd

import (
    "embed"
	"path"

    "github.com/nmusey/letsgo/cli/utils"
)

type App struct {
	Name  string
    Repo string
	Paths AppPaths
}

type AppPaths struct {
	Root    string
	folders []string
}

//go:embed templates/*
var templateFolder embed.FS

func NewApp(name string, repo string, root string) error {
	paths := AppPaths{
		Root:    path.Join(root, name),
		folders: []string{},
	}

	app := App{
		Name:  name,
        Repo: repo,
		Paths: paths,
	}

	if err := app.initPaths(paths); err != nil {
		return err
	}

	initFiles := getInitFiles()
    
	replacements := map[string]string{
		"$appName": app.Name,
		"$appRepo": repo,
		"$dbPort":  "5432",
	}

	if err := app.checkFiles(initFiles, replacements); err != nil {
		return err
	}

	return nil
}

func getInitFiles() map[string]string {
    templates := make(map[string]string)
    files, err := templateFolder.ReadDir("templates")
    if err != nil {
        panic(err)
    }

    for _, file := range files {
        filename := file.Name()
        contents, _ := templateFolder.ReadFile(filename)
        templates[filename] = string(contents)
    }

    return templates
}

func (app *App) initPaths(paths AppPaths) error {
	utils.UpsertFolder(paths.Root)
	for _, dir := range paths.folders {
		dirPath := path.Join(paths.Root, dir)
		if err := utils.UpsertFolder(dirPath); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) checkFiles(files map[string]string, replacements map[string]string) error {
	for filename, contents := range files{
		if err := app.checkFile(filename, contents, replacements); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) checkFile(filename string, contents string, replacements map[string]string) error {
	outpath := path.Join(app.Paths.Root, filename)
	if err := utils.UpsertFile(outpath); err != nil {
		return err
	}

	return utils.CopyTemplateFile(filename, contents, outpath, replacements)
}

