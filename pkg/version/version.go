package version

import (
	"encoding/json"
	"runtime"
	"runtime/debug"

	"k0s.io"
)

var (
	GitCommitString  string
	GitStateString   string
	GitBranchString  string
	GitSummaryString string
	BuildDateString  string
	VersionString    string
)

var BuildInfo, HasBuildInfo = debug.ReadBuildInfo()
var ModVersion = BuildInfo.Main.Version
var Info = &Version{
	GitCommit:  GitCommitString,
	GitState:   GitStateString,
	GitBranch:  GitBranchString,
	GitSummary: GitSummaryString,
	BuildDate:  BuildDateString,
	Version:    VersionString,
	GoVersion:  runtime.Version(),
}

func init() {
	if HasBuildInfo && ModVersion != "" {
		if VersionString == "" {
			VersionString = ModVersion
		}
		if Info.Version == "" {
			Info.Version = ModVersion
		}
	}
}

func GetVersion() k0s.Version {
	return Info
}

type Version struct {
	GitCommit  string
	GitState   string
	GitBranch  string
	GitSummary string
	BuildDate  string
	Version    string
	GoVersion  string
}

func (v *Version) GetGitCommit() string  { return v.GitCommit }
func (v *Version) GetGitState() string   { return v.GitState }
func (v *Version) GetGitBranch() string  { return v.GitBranch }
func (v *Version) GetGitSummary() string { return v.GitSummary }
func (v *Version) GetBuildDate() string  { return v.BuildDate }
func (v *Version) GetVersion() string    { return v.Version }
func (v *Version) GetGoVersion() string  { return v.GoVersion }

func Decode(data []byte) (k0s.Version, error) {
	v := &Version{}
	err := json.Unmarshal(data, v)
	if err != nil {
		return nil, err
	}
	return v, err
}
