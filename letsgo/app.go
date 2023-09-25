package letsgo

import (
	"os"
	"path"

	"github.com/nmusey/letsgo/utils"
)

type App struct {
	Name  string
	Debug bool
	Paths AppPaths
}

type AppPaths struct {
	Root    string
	folders []string
}

func NewApp(name string, root string) error {
	paths := AppPaths{
		Root:    path.Join(root, name),
		folders: []string{"handlers", "views", "models", "data", "public", "middleware", "migrations"},
	}

	app := App{
		Name:  name,
		Debug: true, // TODO - make this ENV controlled
		Paths: paths,
	}

	if err := app.initPaths(paths); err != nil {
		return err
	}

	if err := app.checkFiles([]string{".env", "env.example"}); err != nil {
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

func (app *App) checkFiles(filenames []string) error {
	for _, filename := range filenames {
		if err := app.checkFile(filename); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) checkFile(filename string) error {
	replacements := map[string]string{
		"appName": app.Name,
		"dbPort":  "5432",
	}

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

	return path.Join(pwd, "templates", filename), nil
}
