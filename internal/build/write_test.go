package build_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/weirdGuy/agentform/internal/build"
)

func write(t *testing.T, dir string, files ...build.File) {
	t.Helper()
	if err := build.Write(dir, files); err != nil {
		t.Fatalf("Write: unexpected error: %v", err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	return string(data)
}

func TestWriteCreatesOutputDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "gen", "dev")

	write(t, dir,
		build.File{Path: "main.py", Data: []byte("main\n")},
		build.File{Path: "pkg/util.py", Data: []byte("util\n")},
	)

	if got := readFile(t, filepath.Join(dir, "main.py")); got != "main\n" {
		t.Errorf("main.py = %q, want %q", got, "main\n")
	}
	if got := readFile(t, filepath.Join(dir, "pkg", "util.py")); got != "util\n" {
		t.Errorf("pkg/util.py = %q, want %q", got, "util\n")
	}
	if _, err := os.Stat(filepath.Join(dir, build.Marker)); err != nil {
		t.Errorf("marker file: %v", err)
	}
}

func TestWriteIntoEmptyExistingDir(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, build.File{Path: "main.py", Data: []byte("main\n")})
	if got := readFile(t, filepath.Join(dir, "main.py")); got != "main\n" {
		t.Errorf("main.py = %q, want %q", got, "main\n")
	}
}

func TestWriteRefusesForeignDir(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "precious.txt"), []byte("mine"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := build.Write(dir, []build.File{{Path: "main.py", Data: []byte("main\n")}})
	if err == nil {
		t.Fatal("Write: expected error for non-empty unmarked directory, got nil")
	}
	for _, want := range []string{dir, build.Marker} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("Write error = %q\nwant substring %q", err, want)
		}
	}
	if got := readFile(t, filepath.Join(dir, "precious.txt")); got != "mine" {
		t.Errorf("precious.txt = %q, want untouched", got)
	}
}

func TestWriteRemovesStaleFiles(t *testing.T) {
	dir := t.TempDir()
	write(t, dir,
		build.File{Path: "main.py", Data: []byte("v1\n")},
		build.File{Path: "pkg/util.py", Data: []byte("util\n")},
	)

	write(t, dir, build.File{Path: "main.py", Data: []byte("v2\n")})

	if got := readFile(t, filepath.Join(dir, "main.py")); got != "v2\n" {
		t.Errorf("main.py = %q, want %q", got, "v2\n")
	}
	if _, err := os.Stat(filepath.Join(dir, "pkg", "util.py")); !os.IsNotExist(err) {
		t.Errorf("pkg/util.py should be removed, stat err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "pkg")); !os.IsNotExist(err) {
		t.Errorf("emptied pkg/ directory should be removed, stat err = %v", err)
	}
}

func TestWriteLeavesHiddenEntriesAlone(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, build.File{Path: "main.py", Data: []byte("main\n")})

	if err := os.WriteFile(filepath.Join(dir, ".user"), []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".git", "config"), []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}

	write(t, dir, build.File{Path: "main.py", Data: []byte("main\n")})

	if got := readFile(t, filepath.Join(dir, ".user")); got != "keep" {
		t.Errorf(".user = %q, want untouched", got)
	}
	if got := readFile(t, filepath.Join(dir, ".git", "config")); got != "keep" {
		t.Errorf(".git/config = %q, want untouched", got)
	}
}

func TestWriteSkipsUnchangedFiles(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, build.File{Path: "main.py", Data: []byte("main\n")})

	past := time.Now().Add(-time.Hour)
	target := filepath.Join(dir, "main.py")
	if err := os.Chtimes(target, past, past); err != nil {
		t.Fatal(err)
	}

	write(t, dir, build.File{Path: "main.py", Data: []byte("main\n")})

	info, err := os.Stat(target)
	if err != nil {
		t.Fatal(err)
	}
	if !info.ModTime().Equal(past) {
		t.Errorf("unchanged main.py was rewritten: mtime = %v, want %v", info.ModTime(), past)
	}
}

func TestWriteRejectsBadPath(t *testing.T) {
	dir := t.TempDir()
	err := build.Write(dir, []build.File{{Path: "../escape.py", Data: []byte("x")}})
	if err == nil || !strings.Contains(err.Error(), "escapes") {
		t.Errorf("Write error = %v, want path-escape error", err)
	}
}
