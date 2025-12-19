/*

   Copyright 2016-present Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

// Package common This package provides basic constants used by forms packages.
package common

import (
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/gosimple/slug"
	"github.com/webx-top/com"
)

var TplFuncs = func() template.FuncMap {
	return template.FuncMap{
		`RandomString`:    RandomString,
		`Slugify`:         Slugify,
		`SlugifyMaxWidth`: SlugifyMaxWidth,
	}
}

// RandomString 生成一个指定长度的随机字母数字字符串
// 如果未指定长度，则默认生成8位长度的字符串
func RandomString(length ...uint) string {
	if len(length) > 0 && length[0] > 0 {
		return com.RandomAlphanumeric(length[0])
	}
	return com.RandomAlphanumeric(8)
}

// Slugify 将字符串转换为URL友好的slug格式
func Slugify(v string) string {
	return slug.Make(v)
}

// SlugifyMaxWidth 将字符串转换为slug格式并限制最大长度
// v: 需要转换的原始字符串
// maxWidth: 返回字符串的最大长度限制
// 返回: 转换后的slug字符串，长度不超过maxWidth
func SlugifyMaxWidth(v string, maxWidth int) string {
	return com.Substr(slug.Make(v), ``, maxWidth)
}

// ParseFiles 从给定的文件路径中解析模板文件并返回一个模板对象。
// 如果配置了文件系统(FileSystem)，则使用ParseFS进行解析，否则从本地文件系统读取。
// 第一个文件将作为主模板，其余文件将作为附加模板进行解析。
// 返回解析后的模板对象和可能发生的错误。
func ParseFiles(files ...string) (*template.Template, error) {
	if !FileSystem.IsEmpty() {
		return ParseFS(FileSystem, files...)
	}
	name := filepath.Base(files[0])
	b, err := os.ReadFile(files[0])
	if err != nil {
		return nil, err
	}
	tmpl := template.New(name)
	tmpl.Funcs(TplFuncs())
	tmpl = template.Must(tmpl.Parse(string(b)))
	if len(files) > 1 {
		tmpl, err = tmpl.ParseFiles(files[1:]...)
	}
	return tmpl, err
}

// ParseFS 从给定的文件系统(fs)中解析指定的模板文件，并返回一个模板对象。
// 第一个文件将作为主模板，后续文件将作为子模板解析。
// 参数：
//
//	fs: 文件系统接口
//	files: 要解析的模板文件路径列表
//
// 返回值：
//
//	*template.Template: 解析后的模板对象
//	error: 解析过程中遇到的错误
func ParseFS(fs fs.FS, files ...string) (*template.Template, error) {
	name := filepath.Base(files[0])
	tmpl := template.New(name)
	tmpl.Funcs(TplFuncs())
	fp, err := fs.Open(files[0])
	if err != nil {
		return tmpl, err
	}
	b, err := io.ReadAll(fp)
	fp.Close()
	if err != nil {
		return tmpl, err
	}
	tmpl = template.Must(tmpl.Parse(string(b)))
	if len(files) > 1 {
		tmpl, err = tmpl.ParseFS(fs, files[1:]...)
	}
	return tmpl, err
}
