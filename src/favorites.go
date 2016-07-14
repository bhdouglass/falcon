package main

import (
    "io/ioutil"
    "log"
    "regexp"
    "strings"
)

func (falcon *Falcon) extractId(id string) string {
    id = strings.ToLower(id)

    if id != "" {
        pos := strings.LastIndex(id, "_")

        if pos >= 0 {
            if r, err := regexp.Compile(`[\d\.]*`); err == nil {
                if r.MatchString(id[(pos + 1):len(id)]) {
                    id = id[0:pos]
                }
            } else {
                log.Println(err)
            }
        }
    }

    return id
}

func (falcon *Falcon) favorite(appId string) {
    falcon.favorites = append(falcon.favorites, falcon.extractId(appId))

    falcon.saveFavorites()
}

func (falcon *Falcon) unfavorite(appId string) {
    var newFavorites []string

    for _, id := range falcon.favorites {
        if falcon.extractId(id) != falcon.extractId(appId) {
            newFavorites = append(newFavorites, id)
        }
    }

    falcon.favorites = newFavorites
    falcon.saveFavorites()
}

func (falcon *Falcon) saveFavorites() {
    data := []byte(strings.Join(falcon.favorites, "\n"))
    if err := ioutil.WriteFile(falcon.favFile, data, 0777); err != nil {
        log.Println(err)
    }
}

func (falcon *Falcon) loadFavorites() {
    content, err := ioutil.ReadFile(falcon.favFile)
    if err != nil {
        log.Println(err)
    } else {
        falcon.favorites = strings.Split(string(content), "\n")

        for i := 0; i < len(falcon.favorites); i++ {
            falcon.favorites[i] = falcon.extractId(falcon.favorites[i])
        }
    }
}

func (falcon *Falcon) isFavorite(appId string) bool {
    var isFav bool = false
    appId = falcon.extractId(appId)

    for _, id := range falcon.favorites {
        if id == appId {
            isFav = true
            break
        }
    }

    return isFav
}
