package main

type Settings struct {
    Layout int64 `json:"layout"`
}

type ActionInfo struct {
	Id    string `json:"id"`
	Label string `json:"label"`
	Icon  string `json:"icon,omitempty"`
	Uri   string `json:"uri,omitempty"`
}

type Application struct {
    Id      string
    Title   string
    Comment string
    Icon    string
    Uri     string
    Desktop string
    IsApp   bool
    Sort    string
}

type RemoteScope struct {
    Id          string `json:id`
    Name        string `json:name`
    Icon        string `json:icon`
    Description string `json:description`
}

type Applications []Application

type IconPack struct {
    Title      string `json:title`
    Archive    string `json:archive`
    Author     string `json:"author,omitempty"`
    Maintainer string `json:"maintainer,omitempty"`
    Icon       string `json:icon`
    Preview    string `json:preview`
    Comment    string `json:"comment,omitempty"`
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
