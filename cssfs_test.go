/*
   Copyright © 2020 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package cssfs

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"net/http"
	"os"
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

func Test_removeCSSwhitespace(t *testing.T) {
	c0 := []byte(``)
	w0 := []byte(``)
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
		background : #f9f9f3;
	}
}

@media print {
	body {
		background : #fff;
	}
}

`)
	w3 := []byte(`@media screen{body{background:#f9f9f3}}@media print{body{background:#fff}}`)

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeCSSwhitespace(tt.args.aCSS); string(got) != string(tt.want) {
				t.Errorf("removeCSSwhitespace() = '%s',\nwant '%s'", got, tt.want)
			}
		})
	}
} // Test_removeCSSwhitespace()

func Test_tCSSFilesFilesystem_createMinFile(t *testing.T) {
	c1 := `./css/stylesheet.css`
	c2 := `./css/dark.css`
	c3 := `./css/light.css`
	c4 := `./css/fonts.css`
	fs := newFS(`./`)
	defer func() {
		_ = os.Remove(c1 + cssNameSuffix)
		_ = os.Remove(c2 + cssNameSuffix)
		_ = os.Remove(c3 + cssNameSuffix)
		_ = os.Remove(c4 + cssNameSuffix)
	}()
	type args struct {
		aName string
	}
	tests := []struct {
		name    string
		fields  tCSSFilesFilesystem
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", fs, args{c1}, false},
		{" 2", fs, args{c2}, false},
		{" 3", fs, args{c3}, false},
		{" 4", fs, args{c4}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cf := tCSSFilesFilesystem{
				fs:   tt.fields.fs,
				root: tt.fields.root,
			}
			if err := cf.createMinFile(tt.args.aName); (err != nil) != tt.wantErr {
				t.Errorf("tCSSFilesFilesystem.createMinFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
} // Test_tCSSFilesFilesystem_createMinFile()

func Test_tCSSFilesFilesystem_Open(t *testing.T) {
	c1 := `./css/stylesheet.css`
	c2 := `./css/dark.css`
	c3 := `./css/light.css`
	c4 := `./css/fonts.css`
	fs := newFS(`./`)
	defer func() {
		_ = os.Remove(c1 + cssNameSuffix)
		_ = os.Remove(c2 + cssNameSuffix)
		_ = os.Remove(c3 + cssNameSuffix)
		_ = os.Remove(c4 + cssNameSuffix)
	}()
	var hf http.File

	type args struct {
		aName string
	}
	tests := []struct {
		name    string
		cf      tCSSFilesFilesystem
		args    args
		want    http.File
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 0", fs, args{`doesn't exist`}, hf, true},
		{" 1", fs, args{c1}, hf, false},
		{" 2", fs, args{c2}, hf, false},
		{" 3", fs, args{c3}, hf, false},
		{" 4", fs, args{c4}, hf, false},
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
