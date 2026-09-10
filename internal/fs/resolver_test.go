package fs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/azin-lang/Azin/internal/fs"
)

func TestResolverRelative(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "helper.az")
	if err := os.WriteFile(src, nil, 0o600); err != nil {
		t.Fatal(err)
	}

	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	subSrc := filepath.Join(sub, "util.az")
	if err := os.WriteFile(subSrc, nil, 0o600); err != nil {
		t.Fatal(err)
	}

	r := fs.NewResolver(nil)

	got, err := r.Resolve("./helper", dir)
	if err != nil {
		t.Fatalf("Resolve ./helper: %v", err)
	}
	if got != src {
		t.Errorf("got %q, want %q", got, src)
	}

	got2, err := r.Resolve("./sub/util", dir)
	if err != nil {
		t.Fatalf("Resolve ./sub/util: %v", err)
	}
	if got2 != subSrc {
		t.Errorf("got %q, want %q", got2, subSrc)
	}
}

func TestResolverRelativeNotFound(t *testing.T) {
	dir := t.TempDir()
	r := fs.NewResolver(nil)

	_, err := r.Resolve("./nonexistent", dir)
	if err == nil {
		t.Fatal("expected error for nonexistent import")
	}
}

func TestResolverSearchPaths(t *testing.T) {
	dir := t.TempDir()

	lib1 := filepath.Join(dir, "lib1")
	if err := os.Mkdir(lib1, 0o755); err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(lib1, "mymath.az")
	if err := os.WriteFile(src, nil, 0o600); err != nil {
		t.Fatal(err)
	}

	r := fs.NewResolver([]string{lib1})

	got, err := r.Resolve("mymath", dir)
	if err != nil {
		t.Fatalf("Resolve mymath: %v", err)
	}
	if got != src {
		t.Errorf("got %q, want %q", got, src)
	}
}

func TestResolverSearchPathsOrder(t *testing.T) {
	dir := t.TempDir()

	lib1 := filepath.Join(dir, "lib1")
	lib2 := filepath.Join(dir, "lib2")
	if err := os.MkdirAll(lib1, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(lib2, 0o755); err != nil {
		t.Fatal(err)
	}

	src1 := filepath.Join(lib1, "foo.az")
	src2 := filepath.Join(lib2, "foo.az")
	if err := os.WriteFile(src1, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(src2, nil, 0o600); err != nil {
		t.Fatal(err)
	}

	r := fs.NewResolver([]string{lib1, lib2})

	got, err := r.Resolve("foo", dir)
	if err != nil {
		t.Fatalf("Resolve foo: %v", err)
	}
	if got != src1 {
		t.Errorf("expected first match %q, got %q", src1, got)
	}
}

func TestResolverSearchNotFound(t *testing.T) {
	dir := t.TempDir()
	r := fs.NewResolver([]string{dir})

	_, err := r.Resolve("noexist", dir)
	if err == nil {
		t.Fatal("expected error for nonexistent import")
	}
}

func TestResolverRelativeWithExtension(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "helper.az")
	if err := os.WriteFile(src, nil, 0o600); err != nil {
		t.Fatal(err)
	}

	r := fs.NewResolver(nil)

	got, err := r.Resolve("./helper.az", dir)
	if err != nil {
		t.Fatalf("Resolve ./helper.az: %v", err)
	}
	if got != src {
		t.Errorf("got %q, want %q", got, src)
	}
}
