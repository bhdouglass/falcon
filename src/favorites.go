package main

import (
    "io/ioutil"
    "log"
    "strings"
)

func (falcon *Falcon) favorite(appId string) {
    falcon.favorites = append(falcon.favorites, appId)

    falcon.saveFavorites()
}

func (falcon *Falcon) unfavorite(appId string) {
    var newFavorites []string

    for _, id := range falcon.favorites {
        if id != appId {
            newFavorites = append(newFavorites, id)
        }
    }

    falcon.favorites = newFavorites
    falcon.saveFavorites()
}

func (falcon *Falcon) saveFavorites() {
    data := []byte(strings.Join(falcon.favorites, "\n"))
    if err := ioutil.WriteFile(falcon.favFile, data, 0777); err != nil {
        log.Fatalln(err)
    }
}

func (falcon *Falcon) loadFavorites() {
    content, err := ioutil.ReadFile(falcon.favFile)
    if err != nil {
        log.Fatalln(err)
    } else {
        falcon.favorites = strings.Split(string(content), "\n")
    }
}

func (falcon *Falcon) isFavorite(appId string) bool {
    var isFav bool = false

    for _, id := range falcon.favorites {
        if id == appId {
            isFav = true
            break
        }
    }

    return isFav
}
