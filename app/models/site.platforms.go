package models

import (
	"errors"
	"fmt"
)

type SiteRootInfo struct {
	DrupalVersion       string
	DefaultTheme        string
	AdministrationTheme string
	PHPConfig           string
	PHPOs               string
	DrushVersion        string
	DrushConfiguration  string
	DrushAliasFiles     string
	DrupalRoot          string
}

type PlatformInfo struct {
	Name       string
	Registered bool
	PlatformId int64
	RootInfo   *SiteRootInfo
	Sites      []*SiteInfo
}

type PlatformInputRequest struct {
	Name string
}

const (
	PlatformListKey = "%v+%v"
)

type PlatformList struct {
	List map[string]*PlatformInfo
}

func (pl *PlatformList) Add(rootPath string, platform *PlatformInfo) {
	if pl.List == nil {
		pl.List = make(map[string]*PlatformInfo)
	}

	key := fmt.Sprintf(PlatformListKey, rootPath, platform.Name)
	pl.List[key] = platform
}

func (pl *PlatformList) Get(rootPath string, name string) (*PlatformInfo, error) {
	key := fmt.Sprintf(PlatformListKey, rootPath, name)

	if info, ok := pl.List[key]; ok {
		return info, nil
	}

	return &PlatformInfo{}, errors.New("Platform not found in list.")
}

func (pl *PlatformList) ToSliceList() []*PlatformInfo {
	var items []*PlatformInfo

	for _, platform := range pl.List {
		items = append(items, platform)
	}

	return items
}

func (pi *PlatformInfo) AddSite(site *SiteInfo) {
	pi.Sites = append(pi.Sites, site)
}
