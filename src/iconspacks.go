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

    //TODO hide if it's the active icon pack
    var buttons []ActionInfo
    buttons = append(buttons, ActionInfo{Id: "icon-pack:install", Label: "Install"})

    actionsWidget := scopes.NewPreviewWidget("actions", "actions")
    actionsWidget.AddAttributeValue("actions", buttons)

    return reply.PushWidgets(actionsWidget)
}

func (falcon *Falcon) iconPackSearch(query string, reply *scopes.SearchReply) error {
    //TODO see about finding a better home for this
    resp, err := http.Get("https://dl.dropboxusercontent.com/u/2138439/falcon/icon-packs.json")
    if err != nil {
        return err
    }

    defer resp.Body.Close()
    decoder := json.NewDecoder(resp.Body)

    var iconPacks []IconPack
    if err = decoder.Decode(&iconPacks); err != nil {
        return err
    }

    iconPackCategory := reply.RegisterCategory("icon-packs", "Icon Packs", "", searchCategoryTemplate)

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

    //TODO conact me button (for icon pack submissions)
    //TODO unset icon pack button

    return nil
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

//Based off code from https://socketloop.com/tutorials/golang-untar-or-extract-tar-ball-archive-example
func (falcon *Falcon) iconPackUntar(title string, filename string) string {
    dir := strings.ToLower(strings.Replace(title, " ", "-", -1))
    dir = fmt.Sprintf("%s/%s/", falcon.base.CacheDirectory(), dir)
    log.Println(dir)

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

    var iconPack IconPack
    if err := result.Get("iconPack", &iconPack); err != nil {
        log.Println(err)
    }

    if actionId == "icon-pack:install" {
        filename := falcon.iconPackDownload(iconPack.Title, iconPack.Archive)
        dir := falcon.iconPackUntar(iconPack.Title, filename)
        falcon.saveIconPack(dir)

        //redirect to blank search
        query := scopes.NewCannedQuery("falcon.bhdouglass_falcon", "", "")
        resp = scopes.NewActivationResponseForQuery(query)
    } else {
        resp = scopes.NewActivationResponse(scopes.ActivationNotHandled)
    }

    return resp
}

func (falcon *Falcon) saveIconPack(dir string) {
    falcon.iconPack = dir

    data := []byte(falcon.iconPack)
    if err := ioutil.WriteFile(falcon.iconPackFile, data, 0777); err != nil {
        log.Println(err)
    }
}

func (falcon *Falcon) loadIconPack() {
    content, err := ioutil.ReadFile(falcon.iconPackFile)
    if err != nil {
        log.Println(err)
    } else {
        falcon.iconPack = string(content)

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
}
