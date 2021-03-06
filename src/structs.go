package main

type Settings struct {
    Layout          int64 `json:"layout"`
    Ids             bool  `json:"ids"`
    SeparateDesktop bool  `json:"separate_desktop"`
    FavoritesLayout int64 `json:"favorites_layout"`
    FavoritesSize   int64 `json:"favorites_size"`
    AppScopeLayout  int64 `json:"app_scope_layout"`
    AppScopeSize    int64 `json:"app_scope_size"`
    ShowScopes      bool  `json:"show_scopes"`
}

type ActionInfo struct {
	Id    string `json:"id"`
	Label string `json:"label"`
	Icon  string `json:"icon,omitempty"`
	Uri   string `json:"uri,omitempty"`
}

type Application struct {
    Id        string
    Title     string
    Comment   string
    Icon      string
    Uri       string
    Desktop   string
    IsApp     bool
    IsDesktop bool
    Sort      string
}

type RemoteScope struct {
    Id          string `json:id`
    Name        string `json:name`
    Icon        string `json:icon`
    Description string `json:description`
}

type Applications []Application

type IconPack struct {
    Title       string `json:"title"`
    Icons       string `json:"icons"`
    Author      string `json:"author,omitempty"`
    Maintainer  string `json:"maintainer,omitempty"`
    Icon        string `json:"icon"`
    Preview     string `json:"preview"`
    Description string `json:"description,omitempty"`
}

type LibertineApp struct {
    DesktopFileName string   `json:"desktop_file_name"`
    Icons           []string `json:"icons"`
    Name            string   `json:"name"`
    NoDisplay       bool     `json:"no_display"`
}

type LibertineApps struct {
    AppLaunchers []LibertineApp `json:"app_launchers"`
    Name string `json:"name"`
}

func (slice Applications) Len() int {
    return len(slice)
}

func (slice Applications) Less(a, b int) bool {
    return slice[a].Sort < slice[b].Sort;
}

func (slice Applications) Swap(a, b int) {
    slice[a], slice[b] = slice[b], slice[a]
}
