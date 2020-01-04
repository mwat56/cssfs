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
		fs http.FileSystem
	}

	// Internal list of regular expressions used by
	// the `removeCSSwhitespace()` function.
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

// `createMinFile()` generates a minified version of `aCSSName` in
// `aMinName` returning a possible I/O error.
//
//	`aCSSName` The filename of the original CSS file.
//	`aMinName` The filename of the minified CSS file.
func createMinFile(aCSSName, aMinName string) error {
	css, err := ioutil.ReadFile(aCSSName) // #nosec G304
	if err != nil {
		return err
	}
	for _, re := range cssREs {
		css = re.regEx.ReplaceAll(css, []byte(re.replace))
	}

	return ioutil.WriteFile(aMinName, css, 0640)
} // createMinFile()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// Open returns a `http.File` containing a minified CSS file.
//
//	`aName` is the name of the CSS file to open.
func (cf tCSSFilesFilesystem) Open(aName string) (http.File, error) {
	cName, _ := filepath.Abs(aName)
	mName := cName + cssNameSuffix

	mInfo, err := os.Stat(mName)
	if (nil != err) || (0 == mInfo.Size()) {
		if err = createMinFile(cName, mName); /* #nosec G104 */ nil != err {
			f, err := cf.fs.Open(cName)
			return tNoDirsFile{f}, err
		}

		f, err := cf.fs.Open(mName)
		return tNoDirsFile{f}, err
	}

	cInfo, err := os.Stat(cName)
	if nil != err {
		return nil, err
	}
	if mTime := mInfo.ModTime(); mTime.Before(cInfo.ModTime()) {
		if err = createMinFile(cName, mName); /* #nosec G104 */ nil != err {
			f, err := cf.fs.Open(cName)
			return tNoDirsFile{f}, err
		}
	}

	f, err := cf.fs.Open(mName)
	return tNoDirsFile{f}, nil
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
//	myHandler := http.FileServer(http.Dir("/tmp")))
//
// To use this implementation you'd use:
//
//	myHandler := css.FileServer(http.Dir("/tmp")))
//
//	`aRoot` The root of the filesystem to serve.
func FileServer(aRoot http.FileSystem) http.Handler {
	return http.FileServer(tCSSFilesFilesystem{aRoot})
} // FileServer()

/* _EoF_ */
