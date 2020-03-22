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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/donyori/gogo/encoding/hex"
	"github.com/donyori/gogo/encoding/hex/helper"
	"github.com/donyori/gogo/file"
)

func MainProcess() error {
	info, err := os.Lstat(Option.Input)
	if err != nil {
		return err
	}

	out := os.Stdout
	if Option.Output != "" {
		if Option.Append {
			out, err = os.OpenFile(Option.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		} else {
			out, err = os.Create(Option.Output)
		}
		if err != nil {
			return err
		}
		defer out.Close() // ignore error
	}

	if info.IsDir() {
		files, err := ioutil.ReadDir(Option.Input)
		if err != nil {
			return err
		}
		filenames := make([]string, len(files))
		for i := range files {
			filenames[i] = files[i].Name()
		}
		_, err = fmt.Fprintln(out, strings.Join(filenames, "\t"))
		return err
	} else if info.Mode()&os.ModeSymlink != 0 {
		dst, err := filepath.EvalSymlinks(Option.Input)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(out, dst)
		return err
	} else if !info.Mode().IsRegular() {
		return errors.New("input file is NOT a regular file, directory or symlink")
	}

	reader, err := file.ReadFile(Option.Input, &file.ReadOption{
		Offset:         Option.Offset,
		Limit:          Option.Limit,
		Raw:            Option.Raw,
		BufferSize:     0,
		BufferWhenOpen: true,
	})
	if err != nil {
		return err
	}
	defer reader.Close() // ignore error

	if Option.Binary {
		// Binary mode.
		_, err = io.CopyN(ioutil.Discard, reader, Option.Start)
		if errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return err
		}

		var w io.Writer = out

		if Option.Hex {
			dumper := hex.NewDumper(out, helper.ExampleDumpConfig(true, Option.HexLineWidth))
			defer dumper.Close() // ignore error
			w = dumper
		}

		_, err = io.CopyN(w, reader, Option.Length)
		if errors.Is(err, io.EOF) {
			return nil
		}
	} else {
		// Text mode.
		for count := int64(0); count < Option.Start; count++ {
			more := true
			for more {
				_, more, err = reader.ReadLine()
				if errors.Is(err, io.EOF) {
					return nil
				} else if err != nil {
					return err
				}
			}
		}

		var line []byte
		for count := int64(0); count < Option.Length; count++ {
			more := true
			for more {
				line, more, err = reader.ReadLine()
				if errors.Is(err, io.EOF) {
					return nil
				} else if err != nil {
					return err
				}
				_, err = out.Write(line)
				if err != nil {
					return err
				}
			}
			_, err = fmt.Fprintln(out)
			if err != nil {
				return err
			}
		}
	}
	return err
}
