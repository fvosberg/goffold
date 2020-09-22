package templates

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestParseToGoService(t *testing.T) {
	tmp, err := ioutil.TempDir("", "goffold-go-service-*")
	if err != nil {
		t.Fatalf("Determining temp dir failed: %s", err)
	}
	var walkedPaths []string
	parseFileFunc = func(source io.Reader, destination io.Writer) error {
		f, ok := source.(stater)
		if !ok {
			t.Fatalf("Unexpected reader (not a file): %#v", source)
		}
		stat, err := f.Stat()
		if err != nil {
			t.Fatalf("reading stats for file failed: %s", err)
		}
		sourcePath := stat.Sys().(*zip.FileHeader).Name
		walkedPaths = append(walkedPaths, sourcePath)
		v, ok := destination.(*os.File)
		if !ok {
			t.Errorf("Unexpected writer for path %s: %#v", sourcePath, destination)
		}
		relDestPath := strings.TrimPrefix(v.Name(), tmp)

		if !strings.HasSuffix(sourcePath, relDestPath) {
			t.Errorf("writer has not the right relative path to the destination %q (writer) %q (template)", v.Name(), stat.Name())
		}

		return nil
	}
	err = ParseTo("go-service", tmp)
	if err != nil {
		t.Fatalf("Parsing templates failed: %s", err)
	}

	expectedPaths := []string{
		".gitignore",
		"Dockerfile",
		"README.md",
		"bin/build.sh",
		"bin/env.sh",
		"bin/run.sh",
		"cmd/service/main.go",
		"test_infrastructure/Dockerfile",
		"test_infrastructure/docker-compose.yml",
	}

expectedLoop:
	for _, expected := range expectedPaths {
		for _, walked := range walkedPaths {
			if strings.HasSuffix(walked, expected) {
				continue expectedLoop
			}
		}
		t.Errorf("Missing expected path: %s", expected)
	}

walkedLoop:
	for _, walked := range walkedPaths {
		for _, expected := range expectedPaths {
			if strings.HasSuffix(walked, expected) {
				continue walkedLoop
			}
		}
		t.Errorf("Got unexpected path: %s", walked)
	}

	t.Logf("Walked paths: %#v", walkedPaths)
}

type stater interface {
	Stat() (os.FileInfo, error)
}
