/*
   Copyright © 2020, 2022 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/
package cssfs

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

// `removeCSSwhitespace()` removes unneeded whitespace from `aCSS`.
//
//	`aCSS` The raw CSS data.
func removeCSSwhitespace(aCSS []byte) []byte {
	for _, re := range cssREs {
		aCSS = re.regEx.ReplaceAll(aCSS, []byte(re.replace))
	}

	return aCSS
} // removeCSSwhitespace()

func Test_gzName(t *testing.T) {
	dir, _ := filepath.Abs(`./`)
	s1 := `./stylesheet`
	w1, _ := filepath.Abs(s1)
	w1 += cssGZnameSuffix
	s2 := `./stylesheetcss`
	w2, _ := filepath.Abs(s2)
	w2 += cssGZnameSuffix
	s3 := `./stylesheet.css`
	w3, _ := filepath.Abs(s3)
	w3 += cssGZnameSuffix
	fs := newFS(dir, true)

	type args struct {
		aFilename string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{" 0", args{``}, `/dev/null`},
		{" 1", args{s1}, w1},
		{" 2", args{s2}, w2},
		{" 3", args{s3}, w3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fs.gzName(tt.args.aFilename); got != tt.want {
				t.Errorf("gzName() = %v,\nwant %v", got, tt.want)
			}
		})
	}
} // Test_gzName()

func Test_minName(t *testing.T) {
	dir, _ := filepath.Abs(`./`)
	s1 := `./stylesheet`
	w1, _ := filepath.Abs(s1)
	w1 += cssMinNameSuffix
	s2 := `./stylesheetcss`
	w2, _ := filepath.Abs(s2)
	w2 += cssMinNameSuffix
	s3 := `./stylesheet.css`
	w3, _ := filepath.Abs(s3)
	w3 = w3[:len(w3)-4] + cssMinNameSuffix
	fs := newFS(dir, false)

	type args struct {
		aFilename string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{" 0", args{``}, `/dev/null`},
		{" 1", args{s1}, w1},
		{" 2", args{s2}, w2},
		{" 3", args{s3}, w3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fs.minName(tt.args.aFilename); got != tt.want {
				t.Errorf("minName() = '%v',\nwant `%v`", got, tt.want)
			}
		})
	}
} // Test_minName()

func Test_removeCSSwhitespace(t *testing.T) {
	c0, w0 := []byte(``), []byte(``)
	c1 := []byte(`/*
this are my css rules
*/
`)
	w1 := []byte(``)
	c2 := []byte(`
body { /* default background colour */
	background : #f9f9f3;
}
`)
	w2 := []byte(`body{background:#f9f9f3}`)
	c3 := []byte(`
@media screen {
	body {
		background : #f9f9f3; /* default background colour */
		display    :block;
	}
}

@media print {
	body {
		background: #fff;
	}
}

`)
	w3 := []byte(`@media screen{body{background:#f9f9f3;display:block}}@media print{body{background:#fff}}`)

	c4 := []byte(`p :link { display: inline; } a: hover { display :none; }`)
	w4 := []byte(`p :link{display:inline}a:hover{display:none}`)
	c5 := []byte(`div.empty { height: 0em; width: 0ex; } #main { max-width: 99%; }`)
	w5 := []byte(`div.empty{height:0;width:0}#main{max-width:99%}`)

	type args struct {
		aCSS []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{" 0", args{c0}, w0},
		{" 1", args{c1}, w1},
		{" 2", args{c2}, w2},
		{" 3", args{c3}, w3},
		{" 4", args{c4}, w4},
		{" 5", args{c5}, w5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeCSSwhitespace(tt.args.aCSS); string(got) != string(tt.want) {
				t.Errorf("removeCSSwhitespace() = '%s',\nwant '%s'", got, tt.want)
			}
		})
	}
} // Test_removeCSSwhitespace()

func Test_tCSSfilesFilesystem_createGZfile(t *testing.T) {
	dir, _ := filepath.Abs(`./`)
	c1 := `./css/stylesheet.css`
	c2 := `.././css/dark.css`
	c3 := `css/light.css`
	c4 := dir + `/css/fonts.css`
	fs := newFS(dir, true)
	defer func() {
		_ = os.Remove(fs.gzName(c1))
		_ = os.Remove(fs.gzName(c2))
		_ = os.Remove(fs.gzName(c3))
		_ = os.Remove(fs.gzName(c4))
	}()

	type args struct {
		aFilename string
	}
	tests := []struct {
		name    string
		fields  tCSSfilesFilesystem
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 0", fs, args{`does not exist`}, true},
		{" 1", fs, args{c1}, false},
		{" 2", fs, args{c2}, false},
		{" 3", fs, args{c3}, false},
		{" 4", fs, args{c4}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cf := tt.fields
			if !strings.HasPrefix(tt.args.aFilename, cf.root) {
				tt.args.aFilename = filepath.Join(cf.root,
					filepath.FromSlash(path.Clean(`/`+tt.args.aFilename)))
			}
			if err := cf.createGZfile(tt.args.aFilename); (err != nil) != tt.wantErr {
				t.Errorf("tCSSfilesFilesystem.createGZfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
} // Test_tCSSfilesFilesystem_createGZfile()

func Test_tCSSFilesFilesystem_createMinFile(t *testing.T) {
	dir, _ := filepath.Abs(`./`)
	c1 := `./css/stylesheet.css`
	c2 := `/css/dark.css`
	c3 := `css/light.css`
	c4 := dir + `/css/fonts.css`
	fs := newFS(dir, false)

	defer func() {
		_ = os.Remove(fs.minName(c1))
		_ = os.Remove(fs.minName(c2))
		_ = os.Remove(fs.minName(c3))
		_ = os.Remove(fs.minName(c4))
	}()

	type args struct {
		aFilename string
	}
	tests := []struct {
		name    string
		fields  tCSSfilesFilesystem
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 0", fs, args{`does not exist`}, true},
		{" 1", fs, args{c1}, false},
		{" 2", fs, args{c2}, false},
		{" 3", fs, args{c3}, false},
		{" 4", fs, args{c4}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cf := tt.fields
			if !strings.HasPrefix(tt.args.aFilename, cf.root) {
				tt.args.aFilename = filepath.Join(cf.root,
					filepath.FromSlash(path.Clean(`/`+tt.args.aFilename)))
			}
			if err := cf.createMinFile(tt.args.aFilename); (err != nil) != tt.wantErr {
				t.Errorf("tCSSFilesFilesystem.createMinFile() error = %v,\nwantErr %v", err, tt.wantErr)
			}
		})
	}
} // Test_tCSSFilesFilesystem_createMinFile()

func Test_tCSSFilesFilesystem_Open(t *testing.T) {
	dir, _ := filepath.Abs(`./`)
	c1 := `./css/stylesheet.css`
	c2 := `/css/dark.css`
	c3 := `css/light.css`
	c4 := `/css/fonts.css`
	fs1 := newFS(dir, false)
	fs2 := newFS(dir, true)

	defer func() {
		_ = os.Remove(fs1.minName(c1))
		_ = os.Remove(fs1.minName(c2))
		_ = os.Remove(fs1.minName(c3))
		_ = os.Remove(fs1.minName(c4))
		_ = os.Remove(fs2.gzName(c1))
		_ = os.Remove(fs2.gzName(c2))
		_ = os.Remove(fs2.gzName(c3))
		_ = os.Remove(fs2.gzName(c4))
	}()

	var hf http.File

	type args struct {
		aName string
	}
	tests := []struct {
		name    string
		cf      tCSSfilesFilesystem
		args    args
		want    http.File
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 0", fs1, args{`doesn't exist`}, hf, true},
		{" 1", fs1, args{c1}, hf, false},
		{" 2", fs1, args{c2}, hf, false},
		{" 3", fs1, args{c3}, hf, false},
		{" 4", fs1, args{c4}, hf, false},
		{"10", fs2, args{`doesn't exist`}, hf, true},
		{"11", fs2, args{c1}, hf, false},
		{"12", fs2, args{c2}, hf, false},
		{"13", fs2, args{c3}, hf, false},
		{"14", fs2, args{c4}, hf, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cf.Open(tt.args.aName)
			if (nil != err) != tt.wantErr {
				t.Errorf("tCSSFilesFilesystem.Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (nil == got) && (!tt.wantErr) {
				t.Errorf("tCSSFilesFilesystem.Open() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_tCSSFilesFilesystem_Open()
