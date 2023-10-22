package main

import (
	"os"
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

	initFiles := []string{
		".env", 
        "env.example", 
        "README.md", 
        "Dockerfile", 
        "docker-compose.yml", 
        ".air.toml",
        "go.mod",
	}

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

func (app *App) checkFiles(filenames []string, replacements map[string]string) error {
	for _, filename := range filenames {
		if err := app.checkFile(filename, replacements); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) checkFile(filename string, replacements map[string]string) error {

	outpath := path.Join(app.Paths.Root, filename)
	if err := utils.UpsertFile(outpath); err != nil {
		return err
	}

	templatePath, err := getTemplatePath(filename)
	if err != nil {
		return err
	}

	return utils.CopyTemplateFile(templatePath, outpath, replacements)
}

func getTemplatePath(filename string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return path.Join(pwd, "cli", "templates", filename), nil
}
