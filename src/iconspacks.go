package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "launchpad.net/go-unityscopes/v2"
    "log"
    "os"
)

const iconPackCategoryTemplate = `{
    "schema-version": 1,
    "template": {
        "category-layout": "grid",
        "collapsed-rows": 0,
        "card-size": "small"
    },
    "components" : {
        "title": "title",
        "subtitle": "subtitle",
        "art": {
            "field": "art",
            "aspect-ratio": 1.13
        }
    }
}`

func (falcon *Falcon) iconPackUtilityPreview(result *scopes.Result, metadata *scopes.ActionMetadata, reply *scopes.PreviewReply) error {
    var subtype string
    if err := result.Get("sub-type", &subtype); err != nil {
        log.Println(err)
    }

    if subtype == "reset" {
        titleWidget := scopes.NewPreviewWidget("title", "header")
        titleWidget.AddAttributeValue("title", "Remove current icon pack")
        titleWidget.AddAttributeValue("subtitle", "Revert back to the default icons")

        var buttons []ActionInfo
        buttons = append(buttons, ActionInfo{Id: "icon-pack:reset", Label: "Revert"})

        actionsWidget := scopes.NewPreviewWidget("actions", "actions")
        actionsWidget.AddAttributeValue("actions", buttons)

        return reply.PushWidgets(titleWidget, actionsWidget)
    } else if subtype == "find" {
        titleWidget := scopes.NewPreviewWidget("title", "header")
        titleWidget.AddAttributeValue("title", "Find new icon packs")
        titleWidget.AddAttributeValue("subtitle", "Download icon pack apps from the app store")

        var buttons []ActionInfo
        buttons = append(buttons, ActionInfo{Id: "icon-pack:find", Uri: result.URI(), Label: "Search now"})

        actionsWidget := scopes.NewPreviewWidget("actions", "actions")
        actionsWidget.AddAttributeValue("actions", buttons)

        return reply.PushWidgets(titleWidget, actionsWidget)
    } else {
        titleWidget := scopes.NewPreviewWidget("title", "header")
        titleWidget.AddAttributeValue("title", "Submit an icon pack")
        titleWidget.AddAttributeValue("subtitle", "Contact the author to submit a new icon pack")

        var buttons []ActionInfo
        buttons = append(buttons, ActionInfo{Id: "icon-pack:contact", Uri: result.URI(), Label: "Contact the author"})

        actionsWidget := scopes.NewPreviewWidget("actions", "actions")
        actionsWidget.AddAttributeValue("actions", buttons)

        return reply.PushWidgets(titleWidget, actionsWidget)
    }

    return nil
}

func (falcon *Falcon) iconPackPreview(result *scopes.Result, metadata *scopes.ActionMetadata, reply *scopes.PreviewReply) error {
    var iconPack IconPack
    if err := result.Get("iconPack", &iconPack); err != nil {
        log.Println(err)
    }

    titleWidget := scopes.NewPreviewWidget("title", "header")
    titleWidget.AddAttributeValue("title", iconPack.Title)

    previewWidget := scopes.NewPreviewWidget("preview", "image")
    previewWidget.AddAttributeValue("source", iconPack.Preview)

    if err := reply.PushWidgets(titleWidget, previewWidget); err != nil {
        return err
    }

    if iconPack.Author != "" {
        authorWidget := scopes.NewPreviewWidget("author", "text")
        authorWidget.AddAttributeValue("text", fmt.Sprintf("<b>Author:</b> %s", iconPack.Author))

        if err := reply.PushWidgets(authorWidget); err != nil {
            return err
        }
    }

    if iconPack.Maintainer != "" {
        maintainerWidget := scopes.NewPreviewWidget("maintainer", "text")
        maintainerWidget.AddAttributeValue("text", fmt.Sprintf("<b>Maintainer:</b> %s", iconPack.Maintainer))

        if err := reply.PushWidgets(maintainerWidget); err != nil {
            return err
        }
    }

    if iconPack.Description != "" {
        descriptionWidget := scopes.NewPreviewWidget("description", "text")
        descriptionWidget.AddAttributeValue("text", iconPack.Description)

        if err := reply.PushWidgets(descriptionWidget); err != nil {
            return err
        }
    }


    if iconPack.Icons != falcon.iconPack {
        var buttons []ActionInfo
        buttons = append(buttons, ActionInfo{Id: "icon-pack:install", Label: "Activate"})

        actionsWidget := scopes.NewPreviewWidget("actions", "actions")
        actionsWidget.AddAttributeValue("actions", buttons)

        if err := reply.PushWidgets(actionsWidget); err != nil {
            return err
        }
    }

    return nil
}

func (falcon *Falcon) iconPackSearch(query string, reply *scopes.SearchReply) error {
    var iconPacks []IconPack
    baseDir := "/opt/click.ubuntu.com/"

    files, err := ioutil.ReadDir(baseDir)
    if err != nil {
        log.Println(err)
    } else {
        for _, f := range files {
            path := baseDir + f.Name() + "/current/icon-pack-data.json";
            if _, err := os.Stat(path); err == nil {
                log.Printf("Found icon pack: %s", path)

                content, err := ioutil.ReadFile(path)
                if err != nil {
                    log.Printf("Error while reading icon pack: %s", path)
                    log.Println(err)
                } else {
                    var iconPack IconPack
                    if err = json.Unmarshal(content, &iconPack); err != nil {
                        log.Printf("Error while parsing icon pack: %s", path)
                        log.Println(err)
                    } else {
                        dir := baseDir + f.Name() + "/current/"
                        iconPack.Icons = dir + iconPack.Icons
                        iconPack.Icon = dir + iconPack.Icon
                        iconPack.Preview = dir + iconPack.Preview

                        iconPacks = append(iconPacks, iconPack)
                    }
                }
            }
        }
    }

    iconPackCategory := reply.RegisterCategory("icon-packs", "Installed Icon Packs", "", iconPackCategoryTemplate)

    for index := range iconPacks {
        iconPack := iconPacks[index]

        result := scopes.NewCategorisedResult(iconPackCategory)
        result.SetURI(iconPack.Icons)
        result.SetTitle(iconPack.Title)
        result.SetArt(iconPack.Icon)
        result.Set("type", "icon-pack")
        result.Set("iconPack", iconPack)

        if err := reply.Push(result); err != nil {
            log.Fatalln(err)
        }
    }

    utilitiesCategory := reply.RegisterCategory("icon-packs-utils", "Utilities", "", iconPackCategoryTemplate)

    findResult := scopes.NewCategorisedResult(utilitiesCategory)
    findResult.SetURI("https://open-store.io/?sort=relevance&search=icon-packs")
    findResult.SetTitle("Find new icon packs")
    findResult.SetArt(falcon.getIcon("find-new-icon-pack", falcon.base.ScopeDirectory() + "/find.svg"))
    findResult.Set("type", "icon-pack-utility")
    findResult.Set("sub-type", "find")
    findResult.SetInterceptActivation()

    if err := reply.Push(findResult); err != nil {
        log.Fatalln(err)
    }

    contactResult := scopes.NewCategorisedResult(utilitiesCategory)
    contactResult.SetURI("https://bhdouglass.com/contact.html")
    contactResult.SetTitle("Submit an icon pack")
    contactResult.SetArt(falcon.getIcon("submit-an-icon-pack", falcon.base.ScopeDirectory() + "/contact.svg"))
    contactResult.Set("type", "icon-pack-utility")
    contactResult.Set("sub-type", "contact")
    contactResult.SetInterceptActivation()

    if err := reply.Push(contactResult); err != nil {
        log.Fatalln(err)
    }

    if falcon.iconPack != "" {
        resetResult := scopes.NewCategorisedResult(utilitiesCategory)
        resetResult.SetURI("reset")
        resetResult.SetTitle("Remove current icon pack")
        resetResult.SetArt(falcon.getIcon("remove-current-icon-pack", falcon.base.ScopeDirectory() + "/reset.svg"))
        resetResult.Set("type", "icon-pack-utility")
        resetResult.Set("sub-type", "reset")
        resetResult.SetInterceptActivation()

        if err := reply.Push(resetResult); err != nil {
            log.Fatalln(err)
        }
    }

    return nil
}

func (falcon *Falcon) iconPackActivate(result *scopes.Result, metadata *scopes.ActionMetadata) *scopes.ActivationResponse {
    var resp *scopes.ActivationResponse

    if result.URI() == "reset" {
        falcon.saveIconPack("")

        //redirect to blank search
        query := scopes.NewCannedQuery("falcon.bhdouglass_falcon", "", "")
        resp = scopes.NewActivationResponseForQuery(query)
    } else {
        resp = scopes.NewActivationResponse(scopes.ActivationNotHandled)
    }

    return resp
}

func (falcon *Falcon) iconPackPerformAction(result *scopes.Result, metadata *scopes.ActionMetadata, widgetId, actionId string) *scopes.ActivationResponse {
    var resp *scopes.ActivationResponse

    if actionId == "icon-pack:install" {
        var iconPack IconPack
        if err := result.Get("iconPack", &iconPack); err != nil {
            log.Println(err)
        }

        falcon.saveIconPack(iconPack.Icons)

        //redirect to blank search
        query := scopes.NewCannedQuery("falcon.bhdouglass_falcon", "", "")
        resp = scopes.NewActivationResponseForQuery(query)
    } else if actionId == "icon-pack:reset" {
        falcon.saveIconPack("")

        //redirect to blank search
        query := scopes.NewCannedQuery("falcon.bhdouglass_falcon", "", "")
        resp = scopes.NewActivationResponseForQuery(query)
    } else {
        resp = scopes.NewActivationResponse(scopes.ActivationNotHandled)
    }

    return resp
}

func (falcon *Falcon) saveIconPack(dir string) {
    if falcon.iconPack != dir {
        falcon.iconPack = dir

        data := []byte(falcon.iconPack)
        if err := ioutil.WriteFile(falcon.iconPackFile, data, 0777); err != nil {
            log.Println(err)
        }

        falcon.refreshIconPack()
    }
}

func (falcon *Falcon) getIcon(id string, fallback string) string {
    iconFile := fallback

    if icon, ok := falcon.iconPackMap[id]; ok {
        checkFile := falcon.iconPack + "/" + icon.(string)
        if _, err := os.Stat(checkFile); err == nil {
            iconFile = checkFile
        }
    }

    return iconFile
}

func (falcon *Falcon) refreshIconPack() {
    content, err := ioutil.ReadFile(falcon.iconPack + "/icon-pack.json")
    if err == nil {
        var i interface{}
        if err := json.Unmarshal(content, &i); err != nil {
            log.Println(err)
        }

        falcon.iconPackMap = i.(map[string]interface{})
    } else {
        log.Println(err)
    }
}

func (falcon *Falcon) loadIconPack() {
    content, err := ioutil.ReadFile(falcon.iconPackFile)
    if err != nil {
        log.Println(err)
    } else {
        falcon.iconPack = string(content)

        falcon.refreshIconPack()
    }
}
