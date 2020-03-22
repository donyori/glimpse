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
	"fmt"
	"os"

	"github.com/donyori/gogo/copyright/agpl3"
)

func main() {
	exitStatus := body()
	if exitStatus != 0 {
		os.Exit(exitStatus)
	}
}

func body() int {
	doResp, err := agpl3.RespShowWC(nil, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error occurs:", err)
		return 1
	}
	if doResp {
		return 0
	}
	err = Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error occurs:", err)
		return 1
	}
	err = MainProcess()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error occurs:", err)
		return 1
	}
	return 0
}
