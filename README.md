# CSSfs

[![golang](https://img.shields.io/badge/Language-Go-green.svg)](https://golang.org)
[![GoDoc](https://godoc.org/github.com/mwat56/cssfs?status.svg)](https://godoc.org/github.com/mwat56/cssfs)
[![Go Report](https://goreportcard.com/badge/github.com/mwat56/cssfs)](https://goreportcard.com/report/github.com/mwat56/cssfs)
[![Issues](https://img.shields.io/github/issues/mwat56/cssfs.svg)](https://github.com/mwat56/cssfs/issues?q=is%3Aopen+is%3Aissue)
[![Size](https://img.shields.io/github/repo-size/mwat56/cssfs.svg)](https://github.com/mwat56/cssfs/)
[![Tag](https://img.shields.io/github/tag/mwat56/cssfs.svg)](https://github.com/mwat56/cssfs/tags)
[![View examples](https://img.shields.io/badge/learn%20by-examples-0077b3.svg)](https://github.com/mwat56/cssfs/blob/master/_demo/demo.go)
[![License](https://img.shields.io/github/mwat56/cssfs.svg)](https://github.com/mwat56/cssfs/blob/master/LICENSE)

----

- [CSSfs](#cssfs)
	- [Purpose](#purpose)
	- [Installation](#installation)
	- [Usage](#usage)
	- [Licence](#licence)

## Purpose

Cascading Style Sheets (CSS) needed a long time to become browser independently usable.
In this long time the capabilities seem to have grown almost beyond recognition.
One thing, however, that didn't change over the last 25 years is the fact that in the end a style sheet is just a plain old text file.
Sure, there is some sort of grammar and a certain vocabulary; but it all boils down to lines/paragraphs, preferably nicely formatted to make it easy to read and edit.

Easy, that is, for humans.
For the machines (browser) almost all of the tabs, spaces, and linebreaks used to format a style sheet are just white space eating memory, white noise to be removed – or at least actively ignored – during reading and interpreting the file.

Additionally, that white space – during transfer – uses bandwidth and hence time and energy.
So, getting rid of all those unneeded characters in the style sheet before transmitting it to the end-user (browser) saves both time and money for all parties involved.
And that's the whole purpose of this little package: To remove the unneeded characters from a style sheet.

## Installation

You can use `Go` to install this package for you:

	go get -u github.com/mwat56/cssfs

or you can just include it directly in your own server's code base

	import "github.com/mwat56/cssfs"

## Usage

This package exports a single function `FileServer()`.
So, while you're used to call

	myCSSDirectory := "./css"
	myCSSHandler := http.FileServer(http.Dir(myCSSDirectory)))

to create a fileserver for your static CSS files now, to use _this_ implementation, you'd just do:

	myCSSHandler := cssfs.FileServer(myCSSDirectory)

That's all.

Internally, whenever a CSS file is requested this package's fileserver checks whether there's already a minified version available and, if so, serves it.
Otherwise it creates the minified version from the original CSS file to be used for this and all following calls.

## Licence

        Copyright © 2020 M.Watermann, 10247 Berlin, Germany
                        All rights reserved
                    EMail : <support@mwat.de>

> This program is free software; you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation; either version 3 of the License, or (at your option) any later version.
>
> This software is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
>
> You should have received a copy of the GNU General Public License along with this program. If not, see the [GNU General Public License](http://www.gnu.org/licenses/gpl.html) for details.

----
