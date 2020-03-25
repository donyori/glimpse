// glimpse. A tool for having a quick view of big files.
// Copyright (C) 2020 Yuan Gao
//
// This file is part of glimpse.
//
// glimpse is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/donyori/gogo/copyright/agpl3"
)

// User's input, a global value.
var Option struct {
	Input        string
	Start        int64
	Length       int64
	Output       string
	Append       bool
	Binary       bool
	Hex          bool
	HexLineWidth int
	Raw          bool
	Offset       int64
	Limit        int64
}

// Parse user's input.
func Parse() error {
	if flag.Parsed() {
		return nil
	}
	flag.StringVar(&Option.Output, "o", "",
		"Output filename. If unset, output to STDOUT.")
	flag.BoolVar(&Option.Append, "a", false,
		"Append the result to the output file. Only take effect when -o is set to a file.")
	flag.BoolVar(&Option.Binary, "b", false,
		"Use binary mode instead of text mode.")
	flag.BoolVar(&Option.Hex, "x", false,
		"Print the content in hex encoding, like `hexdump'. Only take effect when -b is set.")
	flag.IntVar(&Option.HexLineWidth, "w", 16,
		"The number of bytes to print per line. Non-positive values are regarded as using the default value (16). Only take effect when -b and -x are set.")
	flag.BoolVar(&Option.Raw, "r", false,
		"No further process such as decompression, print the file directly.")
	flag.Int64Var(&Option.Offset, "s", 0,
		"Offset of the file to read, in bytes, relative to the origin of the file for positive values, and relative to the end of the file for negative values. Only take effect when the file is a regular file.")
	flag.Int64Var(&Option.Limit, "m", 0,
		"Limit of the file to read, in bytes. Non-positive values are no limit. Only take effect when the file is a regular file.")
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		// Ignore errors in the following print functions.
		programName := filepath.Base(os.Args[0])
		if i := strings.LastIndex(programName, "."); i >= 0 {
			programName = programName[:i]
		}
		fmt.Fprintf(w, "Usage of %s: [options] filename [[start] length]\n", programName)
		fmt.Fprintln(w, "\nIf the file is a regular file, print a part of the file.")
		fmt.Fprintln(w, "If the file is a directory, list the content files.")
		fmt.Fprintln(w, "If the file is a symlink, print the destination filename of the symlink.")
		fmt.Fprintln(w, "start and length are the number of lines in text mode, or the number of bytes in binary mode.")
		fmt.Fprintln(w, "Setting start or length to non-positive values is regarded as using the default value.")
		fmt.Fprintln(w, "start and length only take effect when the file is a regular file.")
		fmt.Fprintln(w, "Default value - start: 0; length: 50 lines (in text mode) or 1024 bytes (in binary mode)")
		fmt.Fprintln(w, "\nOptions:")
		flag.PrintDefaults()
		fmt.Fprintln(w)
		agpl3.PrintNotice(w, programName, "2020", "Yuan Gao", "https://github.com/donyori/glimpse")
	}
	flag.Parse()
	if Option.HexLineWidth <= 0 { // non-positive values are regarded as using the default value
		Option.HexLineWidth = 16
	}
	n := flag.NArg()
	if n == 0 {
		return nil
	} else if n > 3 {
		return fmt.Errorf("too more arguments, got %d, want at most 3", n)
	}
	Option.Input = flag.Arg(0)
	var err error
	if n == 2 {
		Option.Length, err = strconv.ParseInt(flag.Arg(1), 0, 64)
		if err != nil {
			return fmt.Errorf(`invalid argument "length": %v`, err)
		}
	} else if n == 3 {
		Option.Start, err = strconv.ParseInt(flag.Arg(1), 0, 64)
		if err != nil {
			return fmt.Errorf(`invalid argument "start": %v`, err)
		}
		Option.Length, err = strconv.ParseInt(flag.Arg(2), 0, 64)
		if err != nil {
			return fmt.Errorf(`invalid argument "length": %v`, err)
		}
	}
	if Option.Start < 0 { // non-positive values are regarded as using the default value
		Option.Start = 0
	}
	if Option.Length <= 0 { // non-positive values are regarded as using the default value
		if Option.Binary {
			Option.Length = 1024
		} else {
			Option.Length = 50
		}
	}
	return nil
}
