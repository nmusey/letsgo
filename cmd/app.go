package cmd

import (
	"embed"
	"path"
	"path/filepath"

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

//go:embed templates/data/*
var dataFolder embed.FS

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
    
    templateDirectories := map[string]embed.FS {
        "templates": templateFolder, 
        "templates/data": dataFolder,
   }

	replacements := map[string]string{
		"$appName": app.Name,
		"$appRepo": repo,
		"$dbPort":  "5432",
	}

    for directory, fs := range templateDirectories {
        files := getInitFiles(directory, fs)    
        if err := app.checkFiles(files, replacements); err != nil {
            return err
        }
    }

	return nil
}

func getInitFiles(directory string, fs embed.FS) map[string]string {
    templates := make(map[string]string)
    files, err := fs.ReadDir(directory)
    if err != nil {
        panic(err)
    }

    for _, file := range files {
        filename := file.Name()
        contents, _ := fs.ReadFile(path.Join(directory, filename))
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

    // We need to treat go.mod specially because it can't be embedded directly
    if filepath.Base(outpath) == "go_mod" {
        dir := filepath.Dir(outpath)
        outpath = filepath.Join(dir, "go.mod")
    }

	return utils.CopyTemplateFile(contents, outpath, replacements)
}

