/*
Copyright © 2020, 2025 M.Watermann, 10247 Berlin, Germany

	   All rights reserved
	EMail : <support@mwat.de>
*/
package cssfs

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type (
	// Struct holding the ages of the three possible files.
	tCSSages struct {
		CSSAge time.Time // original CSS file
		MinAge time.Time // minified CSS file
		GzAge  time.Time // GZipped CSS file
	}

	// Simple struct embedding a `http.FileSystem` that
	// serves minified CSS files.
	tCSSfilesFilesystem struct {
		fs      http.FileSystem
		root    string
		useGZip bool // create compressed CSS files
	}

	// Internal list of regular expressions used by
	// the `createMinFile()` function.
	tCSSre struct {
		regEx   *regexp.Regexp
		replace string
	}

	// Simple struct embedding a `http.File` which ignores
	// directory reads.
	tNoDirsFile struct {
		http.File
	}
)

const (
	// Filename suffix for trimmed/minified CSS files.
	cssMinNameSuffix = `.min.css`
	// Note that we have to use the `.css` extension since stdlib
	// determines which data/type to send by file extension.

	// Filename suffix for compressed CSS files.
	cssGZnameSuffix = `.gz`
)

var (
	// Regular expressions to find/replace whitespace in a CSS file.
	cssREs = []tCSSre{
		{regexp.MustCompile(`(?s)\s*/\x2A.*?\x2A/\s*`), ` `}, // comment
		{regexp.MustCompile(`\s*([;\{,+!])\s*`), `$1`},       // punctuation
		{regexp.MustCompile(`\s*\}\s*\}\s*`), `}}`},          // dito
		// superfluous measurements units:
		{regexp.MustCompile(`(?i)([\s:])([+-]?0)(?:cm|em|ex|in|mm|pc|pt|px|rem|%)`), `0`},
		{regexp.MustCompile(`\s+(:\w)`), ` $1`},         // colon
		{regexp.MustCompile(`(\w:)\s+`), `$1`},          // dito
		{regexp.MustCompile(`\s+:\s+`), `:`},            // dito
		{regexp.MustCompile(`((\{.*?)\s+:\s*)`), `$2:`}, // dito
		{regexp.MustCompile(`\s*;?\}\s*`), `}`},         // final semicolon
		{regexp.MustCompile(`^\s+`), ``},                // leading whitespace
	}
)

// `createMinFile()` generates a minified version of file `aName`
// returning a possible I/O error.
//
//	`aFilename` The URLpath/filename of the original CSS file.
func (cf tCSSfilesFilesystem) createMinFile(aFilename string) error {
	cssData, err := os.ReadFile(aFilename) // #nosec G304
	if err != nil {
		return err
	}

	for _, re := range cssREs {
		cssData = re.regEx.ReplaceAll(cssData, []byte(re.replace))
	}

	return os.WriteFile(cf.minName(aFilename), cssData, 0644)
} // createMinFile()

// `createGZfile()` generates a minified version of file `aFilename`
// and compresses it, returning a possible I/O error.
//
//	`aName` The URLpath/filename of the original CSS file.
func (cf tCSSfilesFilesystem) createGZfile(aFilename string) error {
	mName := cf.minName(aFilename)
	mFile, err := os.OpenFile(mName, os.O_RDONLY, 0)
	if nil != err {
		// The minified CSS file couldn't be opened
		// hence try to create the minified version:
		if err = cf.createMinFile(aFilename); nil != err {
			// The CSS could not be minified
			// so we use the original CSS further on.
			mName = aFilename
		}

		if mFile, err = os.OpenFile(mName, os.O_RDONLY, 0); nil != err {
			// The original CSS file couldn't be opened:
			// nothing we can do about it.
			return err
		}
	}
	defer mFile.Close()

	zFile, err := os.OpenFile(cf.gzName(aFilename),
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644) // #nosec G302
	if nil != err {
		return err
	}
	defer zFile.Close()

	writer, _ := gzip.NewWriterLevel(zFile, gzip.BestCompression)
	defer writer.Close()

	_, err = io.Copy(writer, mFile)

	return err
} // createGZfile()

// `gzName()` returns the name to use for the compressed CSS file.
//
//	`aFilename` The name of the original CSS file.
func (cf tCSSfilesFilesystem) gzName(aFilename string) string {
	if 0 == len(aFilename) {
		return `/dev/null`
	}

	if !strings.HasPrefix(aFilename, cf.root) {
		aFilename = filepath.Join(cf.root,
			filepath.FromSlash(path.Clean(`/`+aFilename)))
	}

	if result, err := filepath.Abs(aFilename); nil == err {
		return result + cssGZnameSuffix
	}

	return aFilename + cssGZnameSuffix
} // gzName()

// `minName()` returns the name to use for the trimmed CSS file.
//
//	`aFilename` The name of the original CSS file.
func (cf tCSSfilesFilesystem) minName(aFilename string) string {
	if 0 == len(aFilename) {
		return `/dev/null`
	}

	if !strings.HasPrefix(aFilename, cf.root) {
		aFilename = filepath.Join(cf.root,
			filepath.FromSlash(path.Clean(`/`+aFilename)))
	}
	if result, err := filepath.Abs(aFilename); nil == err {
		aFilename = result
	}

	if strings.HasSuffix(aFilename, `.css`) {
		return aFilename[:len(aFilename)-4] + cssMinNameSuffix
	}

	return aFilename + cssMinNameSuffix
} // minName()

// Open returns a `http.File` containing a minified CSS file.
//
//	`aName` The name of the CSS file to open.
func (cf tCSSfilesFilesystem) Open(aFilename string) (http.File, error) {
	var (
		ages                       tCSSages
		cFile, mFile, rFile, zFile http.File
		err, err1                  error
		fInfo                      os.FileInfo
		nullTime                   time.Time
	)
	defer func() {
		if nil != cFile {
			_ = cFile.Close()
		}
		if nil != mFile {
			_ = mFile.Close()
		}
		if nil != zFile {
			_ = zFile.Close()
		}
	}()

	if !strings.HasPrefix(aFilename, cf.root) {
		aFilename = filepath.Join(cf.root,
			filepath.FromSlash(path.Clean(`/`+aFilename)))
	}

	if cFile, err = os.OpenFile(aFilename, os.O_RDONLY, 0); nil == err {
		if fInfo, err = cFile.Stat(); nil == err {
			ages.CSSAge = fInfo.ModTime()
		}
	} else {
		err1 = err // WTF: we can't open the original CSS?
	}

	mName := cf.minName(aFilename)
	if mFile, err = os.OpenFile(mName, os.O_RDONLY, 0); nil == err {
		if fInfo, err = mFile.Stat(); nil == err {
			ages.MinAge = fInfo.ModTime()
		}
	}

	if ages.CSSAge.After(ages.MinAge) {
		if nil != mFile {
			_ = mFile.Close()
			ages.MinAge = nullTime // clear the time
		}
		// Original CSS is younger than the minified file
		// hence we create new minified CSS:
		if err = cf.createMinFile(aFilename); nil == err {
			if mFile, err = os.OpenFile(mName, os.O_RDONLY, 0); nil == err {
				if fInfo, err = mFile.Stat(); nil == err {
					ages.MinAge = fInfo.ModTime()
				}
			}
		}
	}

	if cf.useGZip {
		zName := cf.gzName(aFilename)
		if zFile, err = os.OpenFile(zName, os.O_RDONLY, 0); nil == err {
			if fInfo, err = zFile.Stat(); nil == err {
				ages.GzAge = fInfo.ModTime()
			}
		}
		if ages.MinAge.After(ages.GzAge) {
			if nil != zFile {
				_ = zFile.Close()
				ages.GzAge = nullTime // clear the time
			}
			// Wait a moment to avoid a possible identical
			// file creation time with the minified file:
			time.Sleep(time.Millisecond * 10)

			// Minified CSS is younger than the compressed file
			// hence we create new compressed CSS:
			if err = cf.createGZfile(aFilename); nil == err {
				if zFile, err = os.OpenFile(zName, os.O_RDONLY, 0); nil == err {
					if fInfo, err = zFile.Stat(); nil == err {
						ages.GzAge = fInfo.ModTime()
					}
				}
			}
		}
	}

	if ages.GzAge.After(ages.MinAge) && ages.MinAge.After(ages.CSSAge) {
		rFile, zFile = zFile, nil
		return tNoDirsFile{rFile}, nil
	}

	if ages.MinAge.After(ages.CSSAge) {
		rFile, mFile = mFile, nil
		return tNoDirsFile{rFile}, nil
	}

	rFile, cFile = cFile, nil
	// Here `err1` might be caused by an unsuccessful
	// opening of the supposed original CSS file:
	return tNoDirsFile{rFile}, err1
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

// `newFS()` returns a new `tCSSFilesFilesystem` instance.
//
//	`aRootDir` The root of the filesystem to serve.
//	`aGZip` Flag determining whether to use GZio or not.
func newFS(aRootDir string, aGZip bool) tCSSfilesFilesystem {
	if dir, err := filepath.Abs(aRootDir); nil == err {
		aRootDir = dir
		// In case of errors we try `aRootDir` as given.
	}

	return tCSSfilesFilesystem{
		fs:      http.Dir(aRootDir),
		root:    aRootDir,
		useGZip: aGZip,
	}
} // newFS()

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

	//TODO implement additional func argument

	return http.FileServer(newFS(aRootDir, false))
} // FileServer()

/* _EoF_ */
