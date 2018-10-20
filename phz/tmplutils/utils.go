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


package tmplutils

// /exec.go:19:14: undefined: tmplutils.Add
// ./exec.go:20:14: undefined: tmplutils.Minus
// ./exec.go:21:14: undefined: tmplutils.Div
// ./exec.go:22:14: undefined: tmplutils.Mod
// ./exec.go:23:14: undefined: tmplutils.Mul
// ./exec.go:24:14: undefined: tmplutils.Pow
// ./exec.go:25:14: undefined: tmplutils.Sha256
// ./exec.go:26:14: undefined: tmplutils.Argon2id
// FAIL    x/phzd/phz [build failed]

func Add(x, y int) int {
	return x + y
}

func Minus(x, y int) int {
	return x - y
}
func Div(x, y int) int         { return 0 }
func Mod(x, y int) int         { return 0 }
func Mul(x, y int) int         { return 0 }
func Pow(x, y int) int         { return 0 }
func Sha256(b []byte) []byte   { return []byte{5, 4, 3, 2, 1} }
func Argon2id(b []byte) []byte { return []byte{1, 2, 3, 4, 5} }
