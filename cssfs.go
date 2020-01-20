/*
   Copyright © 2020 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package cssfs

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

const (
	// Filename suffix for trimmed CSS files.
	cssNameSuffix = `.min`
)

type (
	// Simple struct embedding a `http.FileSystem` that
	// serves minified CSS file.
	tCSSFilesFilesystem struct {
		fs   http.FileSystem
		root string
	}

	// Internal list of regular expressions used by
	// the `createMinFile()` function.
	tCSSre struct {
		regEx   *regexp.Regexp
		replace string
	}

	// Simple struct embedding a `http.File` and ignoring
	// directory reads.
	tNoDirsFile struct {
		http.File
	}
)

var (
	// Regular expressions to find whitespace in a CSS file.
	cssREs = []tCSSre{
		{regexp.MustCompile(`(?s)\s*/\x2A.*?\x2A/\s*`), ` `}, /* comment */
		{regexp.MustCompile(`\s*([:;\{,+!])\s*`), `$1`},
		{regexp.MustCompile(`\s*\}\s*\}\s*`), `}}`},
		{regexp.MustCompile(`\s*;?\}\s*`), `}`},
		{regexp.MustCompile(`^\s+`), ``},
		{regexp.MustCompile(`\s+$`), ``},
	}
)

// `createMinFile()` generates a minified version of file `aName`
// returning a possible I/O error.
//
//	`aName` The filename of the original CSS file.
func (cf tCSSFilesFilesystem) createMinFile(aName string) error {
	if !path.IsAbs(aName) {
		aName = filepath.Join(cf.root,
			filepath.FromSlash(path.Clean(`/`+aName)))
	}

	cssData, err := ioutil.ReadFile(aName) // #nosec G304
	if err != nil {
		return err
	}

	for _, re := range cssREs {
		cssData = re.regEx.ReplaceAll(cssData, []byte(re.replace))
	}

	return ioutil.WriteFile(aName+cssNameSuffix, cssData, 0640) // #nosec G302
} // createMinFile()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// Open returns a `http.File` containing a minified CSS file.
//
//	`aName` The name of the CSS file to open.
func (cf tCSSFilesFilesystem) Open(aName string) (http.File, error) {
	mName := aName+cssNameSuffix

	mFile, err := cf.fs.Open(mName)
	if nil != err {
		if err = cf.createMinFile(aName); /* #nosec G104 */ nil != err {
			return nil, err
		}

		mFile, err = cf.fs.Open(mName)
		return tNoDirsFile{mFile}, err
	}

	// The minified file exists; now check whether it's
	// younger than the original CSS file.
	cFile, err := cf.fs.Open(aName)
	if nil != err {
		// The original CSS file got lost?
		return tNoDirsFile{mFile}, nil
	}
	cInfo, _ := cFile.Stat()

	mInfo, _ := mFile.Stat()
	mTime := mInfo.ModTime()
	if mTime.After(cInfo.ModTime()) {
		_ = cFile.Close()
		return tNoDirsFile{mFile}, nil
	}

	_ = mFile.Close()
	if err = cf.createMinFile(aName); /* #nosec G104 */ nil != err {
		return tNoDirsFile{cFile}, err
	}
	_ = cFile.Close()
	mFile, err = cf.fs.Open(mName)

	return tNoDirsFile{mFile}, err
} // Open()

// Readdir reads the contents of the directory associated with file `f`
// and returns a slice of up to `aCount` FileInfo values, as would
// be returned by `Lstat`, in directory order.
//
// NOTE: This implementation ignores `aCount` and returns nothing,
// i.e. both the `FileInfo` list and the `error` are `nil`.
func (f tNoDirsFile) Readdir(aCount int) ([]os.FileInfo, error) {
	return nil, nil
} // Readdir()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// FileServer returns a handler that serves HTTP requests
// with the contents of the file system rooted at `aRoot`.
//
// To use the operating system's file system implementation,
// use `http.Dir()`:
//
//	myHandler := http.FileServer(http.Dir("/cssdir"))
//
// To use this implementation you'd use:
//
//	myHandler := cssfs.FileServer("/cssdir")
//
//	`aRootDir` The root of the filesystem to serve.
func FileServer(aRootDir string) http.Handler {
	return http.FileServer(newFS(aRootDir))
} // FileServer()

// `newFS()` returns a new `tCSSFilesFilesystem` instance.
//
//	`aRootDir` The root of the filesystem to serve.
func newFS(aRootDir string) tCSSFilesFilesystem {
	dir, _ := filepath.Abs(aRootDir)

	return tCSSFilesFilesystem{
		fs:   http.Dir(dir),
		root: dir,
	}
} // newFS()

/* _EoF_ */
