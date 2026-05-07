//go:build mage

package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var Default = Build

var releaseTargets = []struct {
	goos   string
	goarch string
}{
	{"linux", "arm64"},
	{"linux", "amd64"},
	{"windows", "arm64"},
	{"windows", "amd64"},
	{"freebsd", "arm64"},
	{"freebsd", "amd64"},
	{"openbsd", "arm64"},
	{"openbsd", "amd64"},
	{"darwin", "arm64"},
	{"darwin", "amd64"},
}

// Build generates code, builds the web frontend, cross-compiles server tools,
// and writes release archives under ./releases.
func Build() error {
	root, err := repoRoot()
	if err != nil {
		return err
	}

	if err := Protobufs(); err != nil {
		return err
	}
	if err := Frontend(); err != nil {
		return err
	}
	if err := Server(); err != nil {
		return err
	}
	if err := ReleaseBinaries(); err != nil {
		return err
	}
	if err := ReleaseArchives(); err != nil {
		return err
	}

	apk := filepath.Join(root, "frontend", "build", "app", "outputs", "flutter-apk", "app-release.apk")
	if exists(apk) {
		if err := copyFile(apk, filepath.Join(root, "releases", "forget-about-it.apk")); err != nil {
			return err
		}
	}
	return nil
}

// Protobufs regenerates Go and Dart protobuf sources from ../protobufs.
func Protobufs() error {
	root, err := repoRoot()
	if err != nil {
		return err
	}
	protobufs := filepath.Clean(filepath.Join(root, "..", "protobufs"))
	if err := run(protobufs, nil, "buf", "generate", "--template", "buf.gen.go.yaml"); err != nil {
		return err
	}
	return run(protobufs, nil, "buf", "generate", "--template", "buf.gen.dart.yaml")
}

// Frontend runs codegen, builds the web app, and copies it into server/web.
func Frontend() error {
	root, err := repoRoot()
	if err != nil {
		return err
	}
	frontend := filepath.Join(root, "frontend")
	serverWeb := filepath.Join(root, "server", "web")
	webBuild := filepath.Join(frontend, "build", "web")

	if err := run(frontend, nil, "dart", "run", "build_runner", "build", "-d"); err != nil {
		return err
	}
	if err := run(frontend, nil, "flutter", "build", "web", "--release", "--wasm"); err != nil {
		return err
	}

	if err := os.RemoveAll(serverWeb); err != nil {
		return err
	}
	return copyDir(webBuild, serverWeb)
}

// Server regenerates sqlc output for the Go server.
func Server() error {
	root, err := repoRoot()
	if err != nil {
		return err
	}
	return run(filepath.Join(root, "server"), nil, "sqlc", "generate")
}

// ReleaseBinaries cross-compiles server and provision binaries.
func ReleaseBinaries() error {
	root, err := repoRoot()
	if err != nil {
		return err
	}
	serverDir := filepath.Join(root, "server")
	tmpRelease := filepath.Join(root, "tmp-release")

	if err := os.RemoveAll(tmpRelease); err != nil {
		return err
	}

	for _, target := range releaseTargets {
		outDir := filepath.Join(tmpRelease, "forget-about-it-"+target.goos+"-"+target.goarch)
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			return err
		}
		ext := ""
		if target.goos == "windows" {
			ext = ".exe"
		}
		env := []string{"GOOS=" + target.goos, "GOARCH=" + target.goarch}
		if err := run(serverDir, env, "go", "build", "-ldflags", `-extldflags "-static"`, "-o", filepath.Join(outDir, "provision"+ext), "./cmd/provision"); err != nil {
			return err
		}
		if err := run(serverDir, env, "go", "build", "-ldflags", `-extldflags "-static"`, "-o", filepath.Join(outDir, "server"+ext), "./cmd/server"); err != nil {
			return err
		}
	}
	return nil
}

// ReleaseArchives packages tmp-release binaries into ./releases.
func ReleaseArchives() error {
	root, err := repoRoot()
	if err != nil {
		return err
	}
	tmpRelease := filepath.Join(root, "tmp-release")
	releases := filepath.Join(root, "releases")

	if err := os.RemoveAll(releases); err != nil {
		return err
	}
	if err := os.MkdirAll(releases, 0o755); err != nil {
		return err
	}

	for _, target := range releaseTargets {
		name := "forget-about-it-" + target.goos + "-" + target.goarch
		src := filepath.Join(tmpRelease, name)
		if target.goos == "windows" {
			if err := zipDir(src, filepath.Join(releases, name+".zip")); err != nil {
				return err
			}
			continue
		}
		if err := tarGzDir(src, filepath.Join(releases, name+".tar.gz")); err != nil {
			return err
		}
	}
	return nil
}

func repoRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if exists(filepath.Join(wd, "frontend", "pubspec.yaml")) && exists(filepath.Join(wd, "server", "go.mod")) {
			return wd, nil
		}
		next := filepath.Dir(wd)
		if next == wd {
			return "", fmt.Errorf("could not find repo root from %s", wd)
		}
		wd = next
	}
}

func run(dir string, env []string, name string, args ...string) error {
	fmt.Printf("running: %s %s\n", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = append(os.Environ(), env...)
	return cmd.Run()
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func zipDir(src, dst string) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	zw := zip.NewWriter(out)
	defer zw.Close()

	base := filepath.Dir(src)
	return filepath.WalkDir(src, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(base, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		info, err := entry.Info()
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = rel
		header.Method = zip.Deflate
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		_, err = io.Copy(writer, in)
		return err
	})
}

func tarGzDir(src, dst string) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	gw := gzip.NewWriter(out)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	base := filepath.Dir(src)
	return filepath.WalkDir(src, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(base, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = rel
		if runtime.GOOS == "windows" && !entry.IsDir() {
			header.Mode = 0o755
		}
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		_, err = io.Copy(tw, in)
		return err
	})
}
