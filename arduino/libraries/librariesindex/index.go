/*
 * This file is part of arduino-cli.
 *
 * arduino-cli is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * As a special exception, you may use this file as part of a free software
 * library without restriction.  Specifically, if other files instantiate
 * templates or use macros or inline functions from this file, or you compile
 * this file and link it with other files to produce an executable, this
 * file does not by itself cause the resulting executable to be covered by
 * the GNU General Public License.  This exception does not however
 * invalidate any other reasons why the executable file might be covered by
 * the GNU General Public License.
 *
 * Copyright 2017 ARDUINO AG (http://www.arduino.cc/)
 */

package librariesindex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/bcmi-labs/arduino-cli/arduino/resources"

	"github.com/bcmi-labs/arduino-cli/configs"
)

// Index represents the list of libraries available for download
type Index struct {
	Libraries map[string]*Library
}

// Library is a library available for download
type Library struct {
	Name     string
	Releases map[string]*Release

	Index *Index `json:"-"`
}

// Release is a release of a library available for download
type Release struct {
	Author             string
	Version            string
	Maintainer         string
	Sentence           string
	Paragraph          string
	Website            string
	Category           string
	Architectures      []string
	Types              []string
	Resource           *resources.DownloadResource
	InstalledDirectory string

	Library *Library `json:"-"`
}

type indexJSON struct {
	Libraries []indexRelease `json:"libraries"`
}

type indexRelease struct {
	Name            string   `json:"name,required"`
	Version         string   `json:"version,required"`
	Author          string   `json:"author"`
	Maintainer      string   `json:"maintainer"`
	Sentence        string   `json:"sentence"`
	Paragraph       string   `json:"paragraph"`
	Website         string   `json:"website"`
	Category        string   `json:"category"`
	Architectures   []string `json:"architectures"`
	Types           []string `json:"types"`
	URL             string   `json:"url"`
	ArchiveFileName string   `json:"archiveFileName"`
	Size            int64    `json:"size"`
	Checksum        string   `json:"checksum"`
}

// IndexPath is the path of the index file for libraries.
var IndexPath = configs.IndexPath("library_index.json")

// LoadIndex reads a library_index.json and create the corresponding Index
func LoadIndex() (*Index, error) {
	libFile, err := IndexPath.Get()
	if err != nil {
		return nil, fmt.Errorf("getting library_index.json path: %s", err)
	}

	buff, err := ioutil.ReadFile(libFile)
	if err != nil {
		return nil, fmt.Errorf("reading library_index.json: %s", err)
	}

	var i indexJSON
	err = json.Unmarshal(buff, &i)
	if err != nil {
		return nil, fmt.Errorf("parsing library_index.json: %s", err)
	}

	return i.extractIndex()
}

func (i indexJSON) extractIndex() (*Index, error) {
	index := &Index{
		Libraries: map[string]*Library{},
	}
	for _, indexLib := range i.Libraries {
		indexLib.extractLibraryIn(index)
	}
	return index, nil
}

func (indexLib *indexRelease) extractLibraryIn(index *Index) {
	library, exist := index.Libraries[indexLib.Name]
	if !exist {
		library := &Library{
			Name:     indexLib.Name,
			Releases: map[string]*Release{},
		}
		index.Libraries[indexLib.Name] = library
	}
	indexLib.extractReleaseIn(library)
}

func (indexLib *indexRelease) extractReleaseIn(library *Library) {
	release := &Release{
		Version:       indexLib.Version,
		Author:        indexLib.Author,
		Maintainer:    indexLib.Maintainer,
		Sentence:      indexLib.Sentence,
		Paragraph:     indexLib.Paragraph,
		Website:       indexLib.Website,
		Category:      indexLib.Category,
		Architectures: indexLib.Architectures,
		Types:         indexLib.Types,
		Resource: &resources.DownloadResource{
			URL:             indexLib.URL,
			ArchiveFileName: indexLib.ArchiveFileName,
			Size:            indexLib.Size,
			Checksum:        indexLib.Checksum,
			CachePath:       "libraries",
		},
		Library: library,
	}
	library.Releases[indexLib.Version] = release
}
