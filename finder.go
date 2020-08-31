package gotype

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

type defaultSourceFinder struct {
	cache map[string][]string
}

func (s *defaultSourceFinder) GetPackageSourceFiles(packagePath string) ([]string, error) {
	if s.cache == nil {
		s.cache = make(map[string][]string)
	}

	if sourceFiles, ok := s.cache[packagePath]; ok {
		return sourceFiles, nil
	}

	packageDir, err := s.findPackageDir(packagePath)
	if err != nil {
		return nil, err
	}

	goSources, err := s.getGoSourcesInsideDir(packageDir)
	if err != nil {
		return nil, err
	}
	s.cache[packagePath] = goSources

	return goSources, nil
}

func (s *defaultSourceFinder) findPackageDir(packagePath string) (string, error) {
	moduleFile, goModFilePath, err := s.findModuleFile()
	if err != nil {
		return "", err
	}

	if packageDir, ok := s.findPackageDirByModulePath(
		packagePath,
		moduleFile.Module.Mod.Path,
		filepath.Dir(goModFilePath),
	); ok {
		return packageDir, nil
	}

	for _, req := range moduleFile.Require {
		modPath, err := s.findPackagePathByVersion(req.Mod)
		if err != nil {
			return "", err
		}

		if packageDir, ok := s.findPackageDirByModulePath(packagePath, req.Mod.Path, modPath); ok {
			return packageDir, nil
		}
	}

	modPath, err := s.findPackagePathByVersion(module.Version{})
	if err != nil {
		return "", err
	}

	if packageDir, ok := s.findPackageDirByModulePath(packagePath, "", modPath); ok {
		return packageDir, nil
	}

	return "", fmt.Errorf("package %s cannot be found in go.mod file", packagePath)
}

func (s *defaultSourceFinder) findModuleFile() (*modfile.File, string, error) {
	goModFilePath, err := s.findGoModFile()
	if err != nil {
		return nil, "", err
	}

	goModFile, err := os.Open(goModFilePath)
	if err != nil {
		return nil, "", fmt.Errorf("cannot open go.mod file: %w", err)
	}
	defer goModFile.Close()

	goModBytes, err := ioutil.ReadAll(goModFile)
	if err != nil {
		return nil, "", fmt.Errorf("cannot read go.mod file: %w", err)
	}

	moduleFile, err := modfile.Parse("go.mod", goModBytes, nil)
	if err != nil {
		return nil, "", fmt.Errorf("cannot parse go.mod file: %w", err)
	}

	return moduleFile, goModFilePath, nil
}

func (s *defaultSourceFinder) findGoModFile() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot get current working dir: %w", err)
	}

	maxDepth := 15
	for i := 0; i < maxDepth; i++ {
		goModPath := path.Join(wd, "go.mod")
		if _, err := os.Stat(goModPath); errors.Is(err, os.ErrNotExist) {
			if wd == "/" {
				break
			}
			wd = path.Join(wd, "..")
			continue
		}

		return goModPath, nil
	}

	return "", fmt.Errorf("no go.mod file found")
}

func (*defaultSourceFinder) findPackageDirByModulePath(
	targetPackagePath,
	modulePath,
	moduleDir string,
) (packageDir string, ok bool) {
	if !strings.HasPrefix(targetPackagePath, modulePath) {
		return "", false
	}
	return filepath.Join(moduleDir, targetPackagePath[len(modulePath):]), true
}

func (s *defaultSourceFinder) findPackagePathByVersion(modVer module.Version) (string, error) {
	lookupDir := make([]string, 0)

	if gohome, ok := os.LookupEnv("GOHOME"); ok {
		lookupDir = append(lookupDir, filepath.Join(gohome, "pkg", "mod", modVer.Path+"@"+modVer.Version))
	}

	// TODO (jauhararifin): this is for unix-based OS only.
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot get user home dir: %w", err)
	}
	lookupDir = append(lookupDir, filepath.Join(homedir, "go", "pkg", "mod", modVer.Path+"@"+modVer.Version))

	if goroot, ok := os.LookupEnv("GOROOT"); ok {
		lookupDir = append(lookupDir, filepath.Join(goroot, "src"))
	}

	if goInstallDir, err := s.findInstalledGoDir(); err == nil {
		lookupDir = append(lookupDir, goInstallDir)
	}
	lookupDir = append(lookupDir, filepath.Join("/", "usr", "local", "go", "src"))

	if gohome, ok := os.LookupEnv("GOHOME"); ok {
		lookupDir = append(lookupDir, filepath.Join(gohome, "src", "mod", modVer.Path+"@"+modVer.Version))
	}

	for _, modulePath := range lookupDir {
		if dir, ok := s.findPackagePathFromCandidatePath(modulePath); ok {
			return dir, nil
		}
	}

	return "", fmt.Errorf("cannot find module %s", modVer.String())
}

func (*defaultSourceFinder) findInstalledGoDir() (string, error) {
	cmd := exec.Command("whereis", "go")
	outputBuff := bytes.Buffer{}
	cmd.Stdout = &outputBuff
	cmd.Stderr = &outputBuff
	if err := cmd.Run(); err != nil {
		return "", err
	}
	line, _, err := bufio.NewReader(&outputBuff).ReadLine()
	if err != nil {
		return "", err
	}
	return filepath.Join(string(line), "..", ".."), nil
}

func (*defaultSourceFinder) findPackagePathFromCandidatePath(modulePath string) (string, bool) {
	stat, err := os.Stat(modulePath)
	if errors.Is(err, os.ErrNotExist) {
		return "", false
	}

	if !stat.IsDir() {
		return "", false
	}
	return modulePath, true
}

func (*defaultSourceFinder) getGoSourcesInsideDir(dir string) ([]string, error) {
	goSources := make([]string, 0)
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if path == dir {
			return nil
		}

		if err != nil {
			return fmt.Errorf("got error when listing file: %w", err)
		}

		if info.IsDir() {
			return filepath.SkipDir
		}

		if filepath.Ext(info.Name()) == ".go" {
			goSources = append(goSources, path)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("error while traversing the directories: %w", err)
	}
	return goSources, nil
}
