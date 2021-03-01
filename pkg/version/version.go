/*
Copyright 2021 Ciena Corporation

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

package version

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	version       string
	vcsURL        string
	vcsRef        string
	vcsCommitDate string
	vcsDirty      string
	goVersion     string
	os            string
	arch          string
	buildDate     string
)

type VersionSpec struct {
	Version       string `json:"version,omitempty"`
	VcsURL        string `json:"vcsurl,omitempty"`
	VcsRef        string `json:"vcsref,omitempty"`
	VcsCommitDate string `json:"vcscommitdate,omitempty"`
	VcsDirty      *bool  `json:"vcsdirty,omitempty"`
	GoVersion     string `json:"goversion,omitempty"`
	OS            string `json:"os,omitempty"`
	Arch          string `json:"arch,omitempty"`
	BuildDate     string `json:"buildtime,omitempty"`
}

func Version() VersionSpec {
	var dirty *bool = nil
	if vcsDirty != "" {
		if parsed, err := strconv.ParseBool(vcsDirty); err == nil {
			dirty = &parsed
		}
	}
	return VersionSpec{
		Version:       version,
		VcsURL:        vcsURL,
		VcsRef:        vcsRef,
		VcsCommitDate: vcsCommitDate,
		VcsDirty:      dirty,
		GoVersion:     goVersion,
		OS:            os,
		Arch:          arch,
		BuildDate:     buildDate,
	}
}

func (v VersionSpec) String() string {
	var ver strings.Builder
	if v.Version != "" {
		ver.WriteString(fmt.Sprintf("Version:       %s", v.Version))
	}
	if v.VcsURL != "" {
		if ver.Len() != 0 {
			ver.WriteString("\n")
		}
		ver.WriteString(fmt.Sprintf("VcsURL:        %s", v.VcsURL))
	}
	if v.VcsRef != "" {
		if ver.Len() != 0 {
			ver.WriteString("\n")
		}
		ver.WriteString(fmt.Sprintf("VcsRef:        %s", v.VcsRef))
	}
	if v.VcsCommitDate != "" {
		if ver.Len() != 0 {
			ver.WriteString("\n")
		}
		ver.WriteString(fmt.Sprintf("VcsCommitDate: %s", v.VcsCommitDate))
	}
	if v.VcsDirty != nil {
		if ver.Len() != 0 {
			ver.WriteString("\n")
		}
		ver.WriteString(fmt.Sprintf("VcsDirty:      %t", *v.VcsDirty))
	}
	if v.GoVersion != "" {
		if ver.Len() != 0 {
			ver.WriteString("\n")
		}
		ver.WriteString(fmt.Sprintf("GoVersion:     %s", v.GoVersion))
	}
	if v.OS != "" {
		if ver.Len() != 0 {
			ver.WriteString("\n")
		}
		ver.WriteString(fmt.Sprintf("OS:            %s", v.OS))
	}
	if v.Arch != "" {
		if ver.Len() != 0 {
			ver.WriteString("\n")
		}
		ver.WriteString(fmt.Sprintf("Arch:          %s", v.Arch))
	}
	if v.BuildDate != "" {
		if ver.Len() != 0 {
			ver.WriteString("\n")
		}
		ver.WriteString(fmt.Sprintf("BuildDate:     %s", v.BuildDate))
	}
	return ver.String()
}
