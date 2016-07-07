package main

import (
    "fmt"
    "launchpad.net/go-unityscopes/v2"
    "log"
    "strings"
)

type Falcon struct {
    base *scopes.ScopeBase

    iconPack string
    iconPackFile string
    iconPackMap map[string]interface{}

    favFile string
    favorites []string
}

func (falcon *Falcon) Preview(result *scopes.Result, metadata *scopes.ActionMetadata, reply *scopes.PreviewReply, cancelled <-chan bool) error {
    var typ string
    if err := result.Get("type", &typ); err != nil {
        log.Println(err)
    }

    var err error
    if typ == "app" {
        err = falcon.appPreview(result, metadata, reply)
    } else if typ == "icon-pack" {
        err = falcon.iconPackPreview(result, metadata, reply)
    } else if typ == "icon-pack-utility" {
        err = falcon.iconPackUtilityPreview(result, metadata, reply)
    } else {
        log.Fatalln("unknown result type")
    }

    return err
}

func (falcon *Falcon) Search(query *scopes.CannedQuery, metadata *scopes.SearchMetadata, reply *scopes.SearchReply, cancelled <-chan bool) error {
    q := query.QueryString()
    log.Println(fmt.Sprintf("query: %s", q))

    if q == "icon-packs" {
        //TODO support searching within the icon packs
        if err := falcon.iconPackSearch(q, reply); err != nil {
            log.Fatalln(err)
        }
    } else {
        if err := falcon.appSearch(q, reply); err != nil {
            log.Fatalln(err)
        }
    }

    return nil
}

func (falcon *Falcon) PerformAction(result *scopes.Result, metadata *scopes.ActionMetadata, widgetId, actionId string) (*scopes.ActivationResponse, error) {
    var resp *scopes.ActivationResponse
    if strings.Contains(actionId, "icon-pack:") {
        resp = falcon.iconPackPerformAction(result, metadata, widgetId, actionId)
    } else {
        resp = falcon.appPerformAction(result, metadata, widgetId, actionId)
    }

    return resp, nil
}

func (falcon *Falcon) Activate(result *scopes.Result, metadata *scopes.ActionMetadata) (*scopes.ActivationResponse, error) {
    var typ string
    if err := result.Get("type", &typ); err != nil {
        log.Println(err)
    }

    var resp *scopes.ActivationResponse
    if typ == "app" {
        resp = falcon.appActivate(result, metadata)
    } else if typ == "icon-pack-utility" {
        resp = falcon.iconPackActivate(result, metadata)
    } else {
        resp = scopes.NewActivationResponse(scopes.ActivationNotHandled)
    }

    return resp, nil
}

func (falcon *Falcon) SetScopeBase(base *scopes.ScopeBase) {
    falcon.base = base

    if falcon.favFile == "" {
        falcon.favFile = fmt.Sprintf("%s/favorites.txt", falcon.base.CacheDirectory())
        falcon.loadFavorites()
    }

    if falcon.iconPackFile == "" {
        falcon.iconPackFile = fmt.Sprintf("%s/iconPack.txt", falcon.base.CacheDirectory())
        falcon.loadIconPack()
    }
}

func main() {
    log.Println("launching falcon")

    scope := &Falcon{}
    if err := scopes.Run(scope); err != nil {
        log.Fatalln(err)
    }
}
