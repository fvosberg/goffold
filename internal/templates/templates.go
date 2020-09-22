package templates

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"

	_ "github.com/fvosberg/goffold/internal/templates/goservice"
)

//go:generate statik -src=../../templates/go-service -dest . -p goservice

func ParseTo(template, destination string) error {
	if template != "go-service" {
		return errors.New("unsupported template")
	}
	statikTemplates, err := fs.New()
	if err != nil {
		return fmt.Errorf("initialization of statik fs failed: %w", err)
	}
	_, err = statikTemplates.Open("Dockerfile")
	fmt.Printf("Opening err: %s\n", err)
	fs.Walk(statikTemplates, "/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		destinationPath := filepath.Join(destination, path)
		if info.IsDir() {
			return ensureDir(destinationPath)
		}
		dest, err := os.Create(destinationPath)
		if err != nil {
			return fmt.Errorf("creation of file %q failed: %w", destinationPath, err)
		}
		defer dest.Close()
		src, err := statikTemplates.Open(path)
		if err != nil {
			return fmt.Errorf("opening source file %s failed: %w", path, err)
		}
		defer src.Close()
		return parseFileFunc(src, dest)
	})
	if err != nil {
		return fmt.Errorf("walking files failed: %w", err)
	}
	return nil
}

var (
	parseFileFunc = parseFile
)

func parseFile(src io.Reader, destination io.Writer) error {
	return nil
}

type goServiceArgs struct {
	ServiceName      string
	CommandName      string
	ConfigPrefix     string
	DockerRepository string
}

func destinationPath(templateRoot, templatePath, destinationRoot string) (string, error) {
	rel, err := filepath.Rel(templateRoot, templatePath)
	if err != nil {
		return "", fmt.Errorf("determining relative path of template %q to template root %q failed: %w", templatePath, templateRoot, err)
	}
	return filepath.Join(destinationRoot, rel), nil
}

func ensureDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}
