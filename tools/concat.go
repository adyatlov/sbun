package tools

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	taskLogDirName     = "task"
	executorLogDirName = "executor"
)

var stdoutNumberRegexp = regexp.MustCompile(`^stdout(\.[0-9]+)?(\.gz)?$`)
var stderrNumberRegexp = regexp.MustCompile(`^stderr(\.[0-9]*)?(\.gz)?$`)

func Concat(bundlePath string, compress bool) error {
	tasks, err := parseTasks(bundlePath)
	if err != nil {
		return fmt.Errorf("cannot parse tasks when concatenating: %v", err)
	}
	errs := make([]string, 0, 2)
	for _, task := range tasks {
		for _, dir := range []string{"", taskLogDirName, executorLogDirName} {
			dir = filepath.Join(task.DirNameAbsolute, dir)
			if !dirExists(dir) {
				continue
			}
			if err := concat(dir, compress); err != nil {
				errs = append(errs, err.Error())
				break
			}
		}
	}
	if len(errs) != 0 {
		return fmt.Errorf("erros occured when concatenating: %v", strings.Join(errs, "; "))
	}
	return nil
}

func concat(dir string, compress bool) error {
	errs := make([]string, 0, 2)
	if err := concatStdout(dir, compress); err != nil {
		errs = append(errs, err.Error())
	}
	if err := concatStderr(dir, compress); err != nil {
		errs = append(errs, err.Error())
	}
	if len(errs) != 0 {
		return fmt.Errorf("errors when concatenating stderr and stdout logs in dir %v: %v",
			dir, strings.Join(errs, ";"))
	}
	return nil
}

func concatStdout(dir string, compress bool) error {
	return concatInDirectory(dir, stdoutNumberRegexp, compress, filepath.Join(dir, "stdout_all"))
}

func concatStderr(dir string, compress bool) error {
	return concatInDirectory(dir, stderrNumberRegexp, compress, filepath.Join(dir, "stderr_all"))
}

func concatInDirectory(dir string, r *regexp.Regexp, compress bool, outName string) error {
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("cannot read dir %v while concatenating: %v", dir, err.Error())
	}
	paths := make([]string, 0, len(infos))
	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		paths = append(paths, filepath.Join(dir, info.Name()))
	}
	paths = filterPathsByFileName(paths, r)
	if len(paths) == 0 {
		return nil
	}
	sortPathsByFileName(paths, r)
	var out io.WriteCloser
	if compress {
		if filepath.Ext(outName) != ".gz" {
			outName += ".gz"
		}
		w, err := os.Create(outName)
		if err != nil {
			return fmt.Errorf("cannot create file %v while concatenating: %v", outName, err)
		}
		defer closeCloser(w)
		out = gzip.NewWriter(w)
	} else {
		out, err = os.Create(outName)
		if err != nil {
			return fmt.Errorf("cannot create file %v while concatenating: %v", outName, err)
		}
	}
	defer closeCloser(out)
	if err := concatFiles(paths, out); err != nil {
		return fmt.Errorf("cannot concatenate files in dir %v: %v", dir, err)
	}
	if err := removeFiles(paths); err != nil {
		return err
	}
	return nil
}

func concatFiles(paths []string, out io.Writer) error {
	for _, path := range paths {
		r, err := fileReader(path)
		if err != nil {
			return fmt.Errorf("cannot open input file while concatenating: %v", err.Error())
		}
		if _, err := io.Copy(out, r); err != nil {
			return fmt.Errorf("cannot copy bytes while concatenating: %v", err.Error())
		}
		closeCloser(r)
	}
	return nil
}

func removeFiles(paths []string) error {
	errs := make([]string, 0, len(paths))
	for _, path := range paths {
		if err := os.Remove(path); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) != 0 {
		return fmt.Errorf("cannot remove files: " + strings.Join(errs, "; "))
	}
	return nil
}

func filterPathsByFileName(paths []string, r *regexp.Regexp) []string {
	expectedPaths := make([]string, 0, len(paths))
	for _, path := range paths {
		if len(r.FindStringSubmatch(filepath.Base(path))) == 0 {
			continue
		}
		expectedPaths = append(expectedPaths, path)
	}
	return expectedPaths
}

func sortPathsByFileName(paths []string, r *regexp.Regexp) {
	sort.Slice(paths, func(i, j int) bool {
		return fileNumber(filepath.Base(paths[i]), r) > fileNumber(filepath.Base(paths[j]), r)
	})
}

// fileNumber("stdout.1.gz") returns 1, fileNumber("stdout.gz") returns 0
func fileNumber(fileName string, r *regexp.Regexp) int {
	groups := r.FindStringSubmatch(fileName)
	if len(groups) != 3 {
		panic("expected only matching files, got " + fileName)
	}
	if groups[1] == "" {
		return 0
	}
	n, err := strconv.Atoi(groups[1][1:])
	if err != nil {
		panic("unexpected conversion error: " + groups[1] + " is not a number")
	}
	return n
}

func fileReader(path string) (io.ReadCloser, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	if filepath.Ext(path) != ".gz" {
		return r, nil
	}
	gzr, err := gzip.NewReader(r)
	if err != nil {
		_ = r.Close()
		return nil, err
	}
	return newReadParentCloser(gzr, r), nil
}

func newReadParentCloser(rc io.ReadCloser, parent io.Closer) io.ReadCloser {
	return &struct {
		io.Reader
		io.Closer
	}{rc, newParentCloser(rc, parent)}
}

type parentCloser struct {
	io.Closer
	parent io.Closer
}

func newParentCloser(c io.Closer, parent io.Closer) *parentCloser {
	return &parentCloser{c, parent}
}

func (pc *parentCloser) Close() error {
	errs := make([]string, 0, 2)
	if err := pc.Closer.Close(); err != nil {
		errs = append(errs, err.Error())
	}
	if err := pc.parent.Close(); err != nil {
		errs = append(errs, "Cannot close parent: "+err.Error())
	}
	if len(errs) != 0 {
		return fmt.Errorf(strings.Join(errs, ";"))
	}
	return nil
}

func closeCloser(c io.Closer) {
	if err := c.Close(); err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "cannot close: %v", err.Error())
	}
}

func dirExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
