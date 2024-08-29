package generator

import (
	"os"
	"path"
	"strings"

	"golang.org/x/mod/modfile"
)

type GoProject struct {
	modFilePath string
}

func OpenProject(projectFolderPath string) (*GoProject, error) {
	modFilePath, found := findModFile(projectFolderPath)
	if !found {
		return nil, newProjectError("go.mod file not found", nil)
	}
	return &GoProject{
		modFilePath: modFilePath,
	}, nil
}

// AddDependency Adds the specified package to the current project.
func (p *GoProject) AddDependency(packageName string, version string) error {

	data, err := os.ReadFile(p.modFilePath)
	if err != nil {
		return newProjectError(errorCannotReadModFile, err)
	}

	modFile, err := modfile.Parse(p.modFilePath, data, nil)
	if err != nil {
		return newProjectError(errorCannotParseModFile, err)
	}

	for _, req := range modFile.Require {
		if strings.Compare(req.Mod.Path, packageName) == 0 {
			return nil
		}
	}

	if requireErr := modFile.AddRequire(packageName, version); requireErr != nil {
		return newProjectError(errorFailedToAddRequiredDependency, requireErr)
	}

	formattedModData, formatErr := modFile.Format()
	if formatErr != nil {
		return newProjectError(errorCannotFormatModFile, formatErr)
	}

	if err := os.WriteFile(p.modFilePath, formattedModData, os.ModePerm); err != nil {
		return newProjectError(errorCannotWriteModFile, err)
	}

	return nil
}

func findModFile(projectFolderPath string) (string, bool) {
	modFilePath := path.Join(projectFolderPath, "go.mod")
	if _, err := os.Stat(modFilePath); os.IsNotExist(err) {
		return "", false
	}
	return modFilePath, true
}
