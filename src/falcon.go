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
    var app Application
    if err := result.Get("app", &app); err != nil {
        log.Println(err)
    }

    headerWidget := scopes.NewPreviewWidget("header", "header")
    headerWidget.AddAttributeValue("title", app.Title)

    iconWidget := scopes.NewPreviewWidget("art", "image")
    iconWidget.AddAttributeValue("source", app.Icon)

    commentWidget := scopes.NewPreviewWidget("content", "text")
    commentWidget.AddAttributeValue("text", app.Comment)

    var buttons []ActionInfo
    buttons = append(buttons, ActionInfo{Id: "launch", Uri: app.Uri, Label: "Launch"})

    if falcon.isFavorite(app.Id) {
        buttons = append(buttons, ActionInfo{Id: "unfavorite", Label: "Unfavorite"})
    } else {
        buttons = append(buttons, ActionInfo{Id: "favorite", Label: "Favorite"})
    }

    actionsWidget := scopes.NewPreviewWidget("actions", "actions")
    actionsWidget.AddAttributeValue("actions", buttons)

    messageWidget := scopes.NewPreviewWidget("message", "text")
    if falcon.isFavorite(app.Id) {
        messageWidget.AddAttributeValue("text", "Refresh scope to see changes")
    }

    return reply.PushWidgets(headerWidget, iconWidget, commentWidget, actionsWidget, messageWidget)
}

func (falcon *Falcon) Search(query *scopes.CannedQuery, metadata *scopes.SearchMetadata, reply *scopes.SearchReply, cancelled <-chan bool) error {
    q := query.QueryString()
    log.Println(fmt.Sprintf("query: %s", q))

    if falcon.favFile == "" {
        falcon.favFile = fmt.Sprintf("%s/favorites.txt", falcon.base.CacheDirectory())
        falcon.loadFavorites()
    }

    if err := falcon.addApps(q, reply); err != nil {
        log.Fatalln(err)
    }

    return nil
}

func (falcon *Falcon) PerformAction(result *scopes.Result, metadata *scopes.ActionMetadata, widgetId, actionId string) (*scopes.ActivationResponse, error) {
    var resp *scopes.ActivationResponse

    if actionId == "favorite" {
        var app Application
        if err := result.Get("app", &app); err != nil {
            log.Println(err)
        }

        if app.Id != "" {
            falcon.favorite(app.Id)
        }

        resp = scopes.NewActivationResponse(scopes.ActivationShowPreview)
    } else if actionId == "unfavorite" {
        var app Application
        if err := result.Get("app", &app); err != nil {
            log.Println(err)
        }

        if app.Id != "" {
            falcon.unfavorite(app.Id)
        }

        resp = scopes.NewActivationResponse(scopes.ActivationShowPreview)
    } else {
        resp = scopes.NewActivationResponse(scopes.ActivationNotHandled)
    }

    return resp, nil
}

func (falcon *Falcon) Activate(result *scopes.Result, metadata *scopes.ActionMetadata) (*scopes.ActivationResponse, error) {
    var resp *scopes.ActivationResponse
    var app Application
    if err := result.Get("app", &app); err != nil {
        log.Println(err)
    }

    if app.IsApp {
        //Let the uri handler open the app
        resp = scopes.NewActivationResponse(scopes.ActivationNotHandled)
    } else {
        //Do a canned query so the scopes can be previewed
        query := scopes.NewCannedQuery(app.Id, "", "")
        resp = scopes.NewActivationResponseForQuery(query)
    }

    return resp, nil
}

func (falcon *Falcon) SetScopeBase(base *scopes.ScopeBase) {
    falcon.base = base
}

func main() {
    log.Println("starting up")

    scope := &Falcon{}
    if err := scopes.Run(scope); err != nil {
        log.Fatalln(err)
    }
}
