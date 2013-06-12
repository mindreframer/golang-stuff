// Copyright 2013 Tumblr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package source

import (
	"os"
	"path"
)

type Jail struct {
	root string
}

func NewJail(root string) (*Jail, error) {
	j := &Jail{root}
	if err := j.mkdirs(); err != nil {
		return nil, err
	}
	return j, nil
}

func (j *Jail) mkdirs() error {
	if err := os.MkdirAll(path.Join(j.root, "src"), 0700); err != nil {
		return err
	}
	return nil
}

// AbsPkgPath returns the absolute local path of package pkg within the jail
func (j *Jail) AbsPkgPath(pkgPath string) string {
	return path.Join(j.root, "src", pkgPath)
}

func (j *Jail) MakePkgDir(pkgPaths ...string) error {
	for _, pkgPath := range pkgPaths {
		if err := os.MkdirAll(j.AbsPkgPath(pkgPath), 0700); err != nil {
			return err
		}
	}
	return nil
}

func (j *Jail) CreateSourceFile(pkgPath, fileName string) (*os.File, error) {
	absPath := j.AbsPkgPath(pkgPath)
	if err := os.MkdirAll(absPath, 0770); err != nil {
		return nil, err
	}
	f, err := os.Create(path.Join(absPath, fileName))
	if err != nil {
		return nil, err
	}
	return f, nil
}
