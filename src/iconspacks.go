package main

import (
    "archive/tar"
    "compress/gzip"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "launchpad.net/go-unityscopes/v2"
    "log"
    "net/http"
    "os"
    "strings"
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
    if result.URI() == "reset" {
        titleWidget := scopes.NewPreviewWidget("title", "header")
        titleWidget.AddAttributeValue("title", "Remove current icon pack")
        titleWidget.AddAttributeValue("subtitle", "Revert back to the default icons")

        var buttons []ActionInfo
        buttons = append(buttons, ActionInfo{Id: "icon-pack:reset", Label: "Revert"})

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

    if iconPack.Comment != "" {
        commentWidget := scopes.NewPreviewWidget("comment", "text")
        commentWidget.AddAttributeValue("text", iconPack.Comment)

        if err := reply.PushWidgets(commentWidget); err != nil {
            return err
        }
    }

    var buttons []ActionInfo

    if falcon.iconPackDir(iconPack.Title) == falcon.iconPack {
        buttons = append(buttons, ActionInfo{Id: "icon-pack:install", Label: "Reinstall"})
    } else {
        buttons = append(buttons, ActionInfo{Id: "icon-pack:install", Label: "Install"})
    }

    actionsWidget := scopes.NewPreviewWidget("actions", "actions")
    actionsWidget.AddAttributeValue("actions", buttons)

    return reply.PushWidgets(actionsWidget)
}

func (falcon *Falcon) iconPackSearch(query string, reply *scopes.SearchReply) error {
    //TODO see about finding a better home for this
    resp, err := http.Get("https://falcon.bhdouglass.com/icon-packs.json")
    if err != nil {
        return err
    }

    defer resp.Body.Close()
    decoder := json.NewDecoder(resp.Body)

    var iconPacks []IconPack
    if err = decoder.Decode(&iconPacks); err != nil {
        return err
    }

    iconPackCategory := reply.RegisterCategory("icon-packs", "Icon Packs", "", iconPackCategoryTemplate)

    for index := range iconPacks {
        iconPack := iconPacks[index]

        result := scopes.NewCategorisedResult(iconPackCategory)
        result.SetURI(iconPack.Archive)
        result.SetTitle(iconPack.Title)
        result.SetArt(iconPack.Icon)
        result.Set("type", "icon-pack")
        result.Set("iconPack", iconPack)

        if err := reply.Push(result); err != nil {
            log.Fatalln(err)
        }
    }

    utilitiesCategory := reply.RegisterCategory("icon-packs-utils", "Utilities", "", iconPackCategoryTemplate)

    contactResult := scopes.NewCategorisedResult(utilitiesCategory)
    contactResult.SetURI("http://bhdouglass.com/contact.html")
    contactResult.SetTitle("Submit an icon pack")
    contactResult.SetArt(falcon.base.ScopeDirectory() + "/contact.svg")
    contactResult.Set("type", "icon-pack-utility")
    contactResult.SetInterceptActivation()

    if err := reply.Push(contactResult); err != nil {
        log.Fatalln(err)
    }

    if falcon.iconPack != "" {
        resetResult := scopes.NewCategorisedResult(utilitiesCategory)
        resetResult.SetURI("reset")
        resetResult.SetTitle("Remove current icon pack")
        resetResult.SetArt(falcon.base.ScopeDirectory() + "/reset.svg")
        resetResult.Set("type", "icon-pack-utility")
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

func (falcon *Falcon) iconPackDownload(title string, url string) string {
    filename := strings.ToLower(strings.Replace(title, " ", "-", -1))
    filename = fmt.Sprintf("%s/%s.tar.gz", falcon.base.TmpDirectory(), filename)

    output, err := os.Create(filename)
    if err != nil {
        log.Fatalln(err)
    }
    defer output.Close()

    response, err := http.Get(url)
    if err != nil {
        log.Fatalln(err)
    }
    defer response.Body.Close()

    if _, err := io.Copy(output, response.Body); err != nil {
        log.Fatalln(err)
    }

    return filename
}

func (falcon *Falcon) iconPackDir(title string) string {
    dir := strings.ToLower(strings.Replace(title, " ", "-", -1))
    dir = fmt.Sprintf("%s/%s/", falcon.base.CacheDirectory(), dir)

    return dir
}

//Based off code from https://socketloop.com/tutorials/golang-untar-or-extract-tar-ball-archive-example
func (falcon *Falcon) iconPackUntar(title string, filename string) string {
    dir := falcon.iconPackDir(title)

    if err := os.MkdirAll(dir, os.FileMode(0755)); err != nil {
        log.Fatalln(err)
    }

    file, err := os.Open(filename)
    if err != nil {
        log.Fatalln(err)
    }
    defer file.Close()

    fileReader, err := gzip.NewReader(file)
    if err != nil {
        log.Fatalln(err)
    }
    defer fileReader.Close()

    tarBallReader := tar.NewReader(fileReader)

    for {
        header, err := tarBallReader.Next()
        if err != nil {
            if err == io.EOF {
                break
            }

            log.Fatalln(err)
        }

        fname := header.Name
        pos := strings.Index(fname, "/")
        if pos >= 0 {
            //Remove the top level directory
            fname = fname[(pos + 1):len(fname)]
        }

        switch header.Typeflag {
            case tar.TypeDir: //Directory
                log.Println("Creating directory :", fname)
                if err = os.MkdirAll(dir + fname, os.FileMode(0755)); err != nil {
                    log.Fatalln(err)
                }

            case tar.TypeReg: //File
                log.Println("Untarring :", fname)
                writer, err := os.Create(dir + fname)
                if err != nil {
                    log.Fatalln(err)
                }

                io.Copy(writer, tarBallReader)
                writer.Close()

            default:
                log.Printf("Unable to untar type : %c in file %s", header.Typeflag, fname)
        }
    }

    return dir
}

func (falcon *Falcon) iconPackPerformAction(result *scopes.Result, metadata *scopes.ActionMetadata, widgetId, actionId string) *scopes.ActivationResponse {
    var resp *scopes.ActivationResponse

    if actionId == "icon-pack:install" {
        var iconPack IconPack
        if err := result.Get("iconPack", &iconPack); err != nil {
            log.Println(err)
        }

        filename := falcon.iconPackDownload(iconPack.Title, iconPack.Archive)
        dir := falcon.iconPackUntar(iconPack.Title, filename)
        falcon.saveIconPack(dir)

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

func (falcon *Falcon) refreshIconPack() {
    content, err := ioutil.ReadFile(falcon.iconPack + "icon-pack.json")
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
