package main

import (
    "fmt"
    "launchpad.net/go-unityscopes/v2"
    "log"
)

type Falcon struct {
    base *scopes.ScopeBase
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
    } else {
        log.Fatalln("unknown result type")
    }

    return err
}

func (falcon *Falcon) Search(query *scopes.CannedQuery, metadata *scopes.SearchMetadata, reply *scopes.SearchReply, cancelled <-chan bool) error {
    q := query.QueryString()
    log.Println(fmt.Sprintf("query: %s", q))

    if err := falcon.appSearch(q, reply); err != nil {
        log.Fatalln(err)
    }

    return nil
}

func (falcon *Falcon) PerformAction(result *scopes.Result, metadata *scopes.ActionMetadata, widgetId, actionId string) (*scopes.ActivationResponse, error) {
    log.Println(actionId)
    return falcon.appPerformAction(result, metadata, widgetId, actionId), nil
}

func (falcon *Falcon) Activate(result *scopes.Result, metadata *scopes.ActionMetadata) (*scopes.ActivationResponse, error) {
    var typ string
    if err := result.Get("type", &typ); err != nil {
        log.Println(err)
    }

    var resp *scopes.ActivationResponse
    if typ == "app" {
        resp = falcon.appActivate(result, metadata)
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
}

func main() {
    log.Println("falcon launching")

    scope := &Falcon{}
    if err := scopes.Run(scope); err != nil {
        log.Fatalln(err)
    }
}
