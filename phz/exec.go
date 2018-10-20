// Copyright (c) 2018 aerth. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of aerth nor the names of this project's
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package phz

import (
	logg "log"
	"os"
	"os/exec"
	"strings"
	"x/phzd/phz/tmplutils"

	"github.com/google/shlex"
)

var log = logg.New(os.Stderr, "phz: ", 0)

var DefaultFuncMap = AddMaps(tmplutils.All(), map[string]interface{}{
	"exec": execstring,
	"aqua": execaqua,

	// math
	/*	"add":      tmplutils.Add,
		"minus":    tmplutils.Minus,
		"div":      tmplutils.Div,
		"mod":      tmplutils.Mod,
		"mul":      tmplutils.Mul,
		"pow":      tmplutils.Pow,
		"sha256":   tmplutils.Sha256,
		"argon2id": tmplutils.Argon2id,
	*/
})

func AddMaps(m ...map[string]interface{}) map[string]interface{} {
	out := map[string]interface{}{}
	for _, mm := range m {
		for i, v := range mm {
			out[i] = v
		}
	}
	return out
}

func execstring(s string) string {
	cmdline, err := shlex.Split(os.ExpandEnv(s))
	if err != nil {
		return "error 5"
	}
	return execslice(cmdline)
}

func execslice(cmdline []string) string {
	cmd := exec.Command(cmdline[0])
	cmd.Args = cmdline[0:]
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("\t", err)
		return "eror 7"
	}
	return removenewline(b)
}

func removenewline(b []byte) string {
	return strings.TrimSuffix(strings.TrimSpace(string(b)), "\n")
}
