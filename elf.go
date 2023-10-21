package initrd

import (
	"debug/elf"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

// ELFLibResolver resolves dynamically linked libraries of ELF file. It collects
// the libraries deduplicated for all files resolved with
// [ELFLibResolver.Resolve]. Once all files are resolved, call
// [ELFLibResolver.Libs] to get the complete list of libraries.
type ELFLibResolver struct {
	searchPaths []string
	libs        []string
}

// NewELFLibResolver creates a new [ELFLibResolver] set to the given
// searchPaths. If none are given, a default set is used that.
func NewELFLibResolver(searchPaths ...string) *ELFLibResolver {
	if len(searchPaths) == 0 {
		searchPaths = []string{
			absRootPath("lib"),
			absRootPath("lib64"),
			absRootPath("usr", "lib"),
			absRootPath("usr", "lib64"),
		}
	}
	return &ELFLibResolver{
		searchPaths: searchPaths,
		libs:        make([]string, 0),
	}
}

// Libs returns the list of resolved libraries so far.
func (r *ELFLibResolver) Libs() []string {
	out := make([]string, len(r.libs))
	copy(out, r.libs)
	return out
}

// Resolve analyzes the required linked libraries of the ELF file with the
// given path. The libraries are search for in the library search paths and
// are added with their absolute path to [ELFLibResolver]'s list of libs. Call
// [ELFLibResolver.Libs] once all files are resolved.
func (r *ELFLibResolver) Resolve(elfFilePath string) error {
	libs, err := LinkedLibs(elfFilePath)
	if err != nil {
		return fmt.Errorf("get linked libs: %v", err)
	}

	for _, lib := range libs {
		var found bool
		for _, searchPath := range r.searchPaths {
			path, err := filepath.EvalSymlinks(filepath.Join(searchPath, lib))
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					continue
				}
				return err
			}
			if !slices.Contains(r.libs, path) {
				r.libs = append(r.libs, path)
				if err := r.Resolve(path); err != nil {
					return err
				}
			}
			found = true
			break
		}
		if !found {
			return fmt.Errorf("lib could not be resolved: %s", lib)
		}
	}

	return nil
}

// LinkedLibs fetches the list of dynamically linked libraries from the ELF
// file.
func LinkedLibs(elfFilePath string) ([]string, error) {
	elfFile, err := elf.Open(elfFilePath)
	if err != nil {
		return nil, err
	}
	defer elfFile.Close()

	libs, err := elfFile.ImportedLibraries()
	if err != nil {
		return nil, fmt.Errorf("read libs: %v", err)
	}

	return libs, nil
}
