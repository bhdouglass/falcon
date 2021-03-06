package main

import (
    "encoding/json"
    "fmt"
    "github.com/gosexy/gettext"
    "io/ioutil"
    "launchpad.net/go-unityscopes/v2"
    "log"
    "os"
    "os/exec"
    "sort"
    "strings"
)

const searchCategoryTemplate = `{
    "schema-version": 1,
    "template": {
        "category-layout": "%s",
        "collapsed-rows": 0,
        "card-layout": "%s",
        "card-size": "%s"
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

func (falcon *Falcon) firstChar(str string) string {
    return string([]rune(str)[0])
}

func (falcon *Falcon) appPreview(result *scopes.Result, metadata *scopes.ActionMetadata, reply *scopes.PreviewReply) error {
    var settings Settings
    falcon.base.Settings(&settings)

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

    idWidget := scopes.NewPreviewWidget("id", "text")
    if settings.Ids {
        idWidget.AddAttributeValue("text", app.Id)
    }

    var buttons []ActionInfo
    buttons = append(buttons, ActionInfo{Id: "launch", Label: "Launch"})

    if falcon.isFavorite(app.Id) {
        buttons = append(buttons, ActionInfo{Id: "unfavorite", Label: "Unfavorite"})
    } else {
        buttons = append(buttons, ActionInfo{Id: "favorite", Label: "Favorite"})
    }

    actionsWidget := scopes.NewPreviewWidget("actions", "actions")
    actionsWidget.AddAttributeValue("actions", buttons)

    return reply.PushWidgets(headerWidget, iconWidget, commentWidget, idWidget, actionsWidget)
}

//TODO cache this data
func (falcon *Falcon) getLibertineApps(query string) Applications {
    var appList Applications

    //Note, we don't want to fail if we encounter an error in here (hence the long chain of if/else)
    containerOutput, err := exec.Command("libertine-container-manager", "list").Output()
    if err != nil {
        log.Println("Error while listing libertine containers:")
        log.Println(err)
    } else {
        log.Printf("libertine containers: %s", containerOutput)
        containerList := strings.Split(string(containerOutput), "\n")
        log.Printf("containers: %v", containerList)
        for index := range containerList {
            if len(containerList[index]) > 0 {
                appOutput, err := exec.Command("libertine-container-manager", "list-apps", "--json", "--id", containerList[index]).Output()
                if err != nil {
                    log.Printf("Error while listing apps in %s:", containerList[index])
                    log.Println(err)
                } else {
                    var libertineApps LibertineApps
                    if err = json.Unmarshal(appOutput, &libertineApps); err != nil {
                        log.Printf("Error while decoding apps in %s:", containerList[index])
                        log.Println(err)
                    } else {
                        for jindex := range libertineApps.AppLaunchers {
                            if !libertineApps.AppLaunchers[jindex].NoDisplay {
                                //log.Printf("libertine app: %s", libertineApps.AppLaunchers[jindex].Name)

                                id := libertineApps.AppLaunchers[jindex].DesktopFileName //TODO parse this
                                start := strings.LastIndex(id, "/")
                                end := strings.LastIndex(id, ".desktop")
                                id = id[(start + 1):end]

                                var libertineApp Application
                                libertineApp.Id = id
                                libertineApp.Title = libertineApps.AppLaunchers[jindex].Name
                                libertineApp.Sort = strings.ToLower(libertineApps.AppLaunchers[jindex].Name)
                                libertineApp.Comment = ""
                                libertineApp.Icon = falcon.getIcon(id, libertineApps.AppLaunchers[jindex].Icons[0])
                                libertineApp.Uri = fmt.Sprintf("appid://%s/%s/0.0", containerList[index], id)
                                libertineApp.IsApp = true
                                libertineApp.IsDesktop = true

                                if (query == "" || strings.Index(strings.ToLower(libertineApp.Title), strings.ToLower(query)) >= 0) {
                                    appList = append(appList, libertineApp)
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    return appList
}

func (falcon *Falcon) appSearch(query string, reply *scopes.SearchReply) error {
    var settings Settings
    falcon.base.Settings(&settings)

    paths := []string{
        "/usr/share/applications/",
        "/home/phablet/.local/share/applications/",
    }

    var uappexplorer Application
    var uappexplorerScope Application
    var clickstore Application

    var appList Applications
    for index := range paths {
        path := paths[index]
        files, err := ioutil.ReadDir(path)
        if err != nil {
            log.Println(err)
        } else {
            for _, f := range files {
                content, err := ioutil.ReadFile(path + f.Name())
                if err != nil {
                    log.Fatalln(err)
                } else {
                    lines := strings.Split(string(content), "\n")

                    var app = Application{}
                    app.Desktop = string(content)
                    app.Uri = "application:///" + f.Name()
                    app.IsApp = true
                    app.IsDesktop = false

                    skip := true
                    nodisplay := false
                    onlyShowIn := "unity"

                    desktopMap := map[string] string{}
                    for _, line := range lines {
                        if strings.Contains(line, "Desktop Action") {
                            //TODO refactor this to be smarter when parsing the different sections of the desktop file
                            break
                        }

                        split := strings.Split(line, "=")
                        if (len(split) >= 2) {
                            key := split[0]
                            lkey := strings.ToLower(key)
                            value := strings.Replace(line, key + "=", "", 1)

                            desktopMap[lkey] = value
                        }
                    }

                    if value, ok := desktopMap["name"]; ok {
                        app.Title = value

                        fullLang := os.Getenv("LANG")
                        if (fullLang != "") {
                            split := strings.Split(fullLang, ".")
                            lang := strings.ToLower(split[0])

                            split = strings.Split(lang, "_")
                            shortLang := strings.ToLower(split[0])

                            if title, ok := desktopMap[fmt.Sprintf("name[%s]", lang)]; ok {
                                app.Title = title
                            } else if title, ok := desktopMap[fmt.Sprintf("name[%s]", shortLang)]; ok {
                                app.Title = title
                            }
                        }

                        if domain, ok := desktopMap["x-ubuntu-gettext-domain"]; ok {
                            gettext.BindTextdomain(domain, ".")
                            gettext.Textdomain(domain)
                            gettext.SetLocale(gettext.LC_ALL, "")

                            translation := gettext.Gettext(value)
                            if (translation != "") {
                                app.Title = translation
                            }
                        }

                        app.Sort = strings.ToLower(app.Title)
                    }

                    if value, ok := desktopMap["icon"]; ok {
                        if (value == "media-memory-sd") { //Special exception for the "External Drives" app
                            app.Icon = "file:///usr/share/icons/Humanity/devices/48/media-memory-sd.svg"
                        } else if (value != "" && value[0:1] == "/") {
                            app.Icon = "file://" + value
                        } else {
                            app.Icon = "file:///usr/share/icons/suru/apps/128/placeholder-app-icon.png"
                        }
                    }

                    if value, ok := desktopMap["comment"]; ok {
                        app.Comment = value
                    }

                    if value, ok := desktopMap["x-ubuntu-application-id"]; ok {
                        app.Id = falcon.extractId(value)
                    } else {
                        app.Id = falcon.extractId(strings.Replace(f.Name(), ".desktop", "", 1))
                    }

                    if value, ok := desktopMap["x-ubuntu-touch"]; (ok && strings.ToLower(value) == "true") {
                        skip = false
                    }

                    if value, ok := desktopMap["nodisplay"]; (ok && strings.ToLower(value) == "true") {
                        nodisplay = true
                    }

                    if value, ok := desktopMap["onlyshowin"]; ok {
                        onlyShowIn = strings.ToLower(value)
                    }

                    //Currently the scopes have their data and icons stored under these path
                    if (strings.Contains(app.Icon, "/home/phablet/.local/share/unity-scopes/") || strings.Contains(app.Icon, "/usr/lib/arm-linux-gnueabihf/unity-scopes/") || strings.Contains(app.Icon, "/usr/share/unity/scopes/")) {
                        name := strings.Replace(f.Name(), ".desktop", "", 1)

                        //Don't show this scope
                        if (name != "falcon.bhdouglass_falcon" && name != "com.canonical.scopes.clickstore") {
                            app.Id = name
                            //Setting a scope uri seems to have the unfortunate side effect of preventing a preview so Falcon handles the activation directly
                            //app.Uri = fmt.Sprintf("scope://%s", name)
                            app.Uri = name

                            nodisplay = false
                            skip = false
                            app.IsApp = false
                        }
                    }

                    app.Icon = falcon.getIcon(app.Id, app.Icon)

                    if (!skip && !nodisplay && onlyShowIn == "unity") {
                        if (strings.Contains(app.Id, "uappexplorer.bhdouglass")) {
                            uappexplorer = app
                        } else if (strings.Contains(app.Id, "uappexplorer-scope.bhdouglass")) {
                            uappexplorerScope = app
                        } else if (strings.Contains(app.Id, "openstore.openstore-team")) {
                            clickstore = app
                        }

                        if (query == "" || strings.Index(strings.ToLower(app.Title), strings.ToLower(query)) >= 0) {
                            appList = append(appList, app)
                        }
                    }
                }
            }
        }
    }

    //Desktop/Libertine Apps
    appList = append(appList, falcon.getLibertineApps(query)...)

    sort.Sort(appList)

    categories := map[string] *scopes.Category{};

    categoryLayouts := []string{"grid", "carousel", "vertical-journal", "horizontal-list"}
    cardLayouts := []string{"vertical", "vertical", "horizontal", "vertical"}
    cardSizes := []string{"small", "medium", "large"}

    favoritesTemplate := fmt.Sprintf(searchCategoryTemplate, categoryLayouts[settings.FavoritesLayout], cardLayouts[settings.FavoritesLayout], cardSizes[settings.FavoritesSize])
    appScopeTemplate := fmt.Sprintf(searchCategoryTemplate, categoryLayouts[settings.AppScopeLayout], cardLayouts[settings.AppScopeLayout], cardSizes[settings.AppScopeSize])

    categories["favorite"] = reply.RegisterCategory("favorites", "Favorites", "", favoritesTemplate)

    if (settings.Layout == 0) { //Group by apps & scopes
        categories["apps"] = reply.RegisterCategory("apps", "Apps", "", appScopeTemplate)
        categories["desktop"] = reply.RegisterCategory("desktop", "Desktop Apps", "", appScopeTemplate)
        categories["scopes"] = reply.RegisterCategory("scopes", "Scopes", "", appScopeTemplate)
    } else { //Group by first letter
        //TODO ignore A/An/The
        //TODO group numbers

        charMap := map[string] string{}
        for index := range appList {
            char := strings.ToUpper(falcon.firstChar(appList[index].Title))
            charMap[char] = char
        }

        var charList []string
        for index := range charMap {
            charList = append(charList, index)
        }

        sort.Strings(charList)
        for index := range charList {
            char := charList[index]
            categories[char] = reply.RegisterCategory(char, char, "", appScopeTemplate)
        }
    }

    searchTitle := "Search for more apps"
    if (query != "") {
        searchTitle = fmt.Sprintf("Search for apps like \"%s\"", query)
    }
    storeCategory := reply.RegisterCategory("store", searchTitle, "", fmt.Sprintf(searchCategoryTemplate, "grid", "vertical", "small"))

    //Favorites
    for index := range appList {
        app := appList[index]

        if falcon.isFavorite(app.Id) {
            result := scopes.NewCategorisedResult(categories["favorite"])
            result.SetURI(app.Uri)
            result.SetTitle(app.Title)
            result.SetArt(app.Icon)
            result.Set("app", app)
            result.Set("type", "app")
            result.SetInterceptActivation()

            if err := reply.Push(result); err != nil {
                log.Fatalln(err)
            }
        }
    }

    //Apps first, or all if they are joined
    for index := range appList {
        app := appList[index]

        //See note at next for loop
        if (settings.Layout == 0 && !app.IsApp) || falcon.isFavorite(app.Id) || (settings.Layout == 0 && settings.SeparateDesktop && app.IsDesktop) {
            continue
        }

        var result *scopes.CategorisedResult
        if (settings.Layout == 0) {
            if (app.IsApp) {
                result = scopes.NewCategorisedResult(categories["apps"])
            } else if (settings.ShowScopes) {
                result = scopes.NewCategorisedResult(categories["scopes"])
            }
        } else {
            char := strings.ToUpper(falcon.firstChar(app.Title))
            result = scopes.NewCategorisedResult(categories[char])

            if (app.IsDesktop) {
                result.Set("subtitle", "Desktop App")
            } else if (app.IsApp) {
                result.Set("subtitle", "App")
            } else {
                result.Set("subtitle", "Scope")
            }
        }

        result.SetURI(app.Uri)
        result.SetTitle(app.Title)
        result.SetArt(app.Icon)
        result.Set("app", app)
        result.Set("type", "app")
        result.SetInterceptActivation()

        if err := reply.Push(result); err != nil {
            log.Fatalln(err)
        }
    }

    //Desktop apps (if separated)
    if settings.SeparateDesktop && settings.Layout == 0 {
        for index := range appList {
            app := appList[index]

            if (!app.IsDesktop || falcon.isFavorite(app.Id)) {
                continue
            }

            result := scopes.NewCategorisedResult(categories["desktop"])
            result.SetURI(app.Uri)
            result.SetTitle(app.Title)
            result.SetArt(app.Icon)
            result.Set("app", app)
            result.Set("type", "app")
            result.SetInterceptActivation()

            if err := reply.Push(result); err != nil {
                log.Fatalln(err)
            }
        }
    }

    //Scopes last
    //TODO This is a really hacky looking way to make sure the apps go before the scopes, figure out a better way to do this
    if (settings.Layout == 0 && settings.ShowScopes) {
        for index := range appList {
            app := appList[index]

            if (app.IsApp || falcon.isFavorite(app.Id)) {
                continue
            }

            result := scopes.NewCategorisedResult(categories["scopes"])
            result.SetURI(app.Uri)
            result.SetTitle(app.Title)
            result.SetArt(app.Icon)
            result.Set("app", app)
            result.Set("type", "app")
            result.SetInterceptActivation()

            if err := reply.Push(result); err != nil {
                log.Fatalln(err)
            }
        }
    }

    //TODO make a setting for this
    //TODO give an option to search online immediately instead of going to another scope/app
    var store Application
    if (uappexplorerScope.Id != "") {
        store = uappexplorerScope

        if (query != "") {
            store.Uri = fmt.Sprintf("scope://%s?q=%s", store.Id, query)
        }
    } else if (uappexplorer.Id != "") {
        store = uappexplorer

        if (query != "") {
            store.Uri = fmt.Sprintf("https://uappexplorer.com/apps?q=%s&sort=relevance", query)
        }
    } else if (clickstore.Id != "") {
        store = clickstore

        if (query != "") {
            store.Uri = fmt.Sprintf("scope://%s?q=%s", store.Id, query)
        }
    }

    if (store.Id != "") {
        result := scopes.NewCategorisedResult(storeCategory)
        result.SetURI(store.Uri)
        result.SetTitle(store.Title)
        result.SetArt(store.Icon)
        result.Set("app", store)
        result.Set("type", "app")
        result.SetInterceptActivation()

        if err := reply.Push(result); err != nil {
            log.Fatalln(err)
        }
    }

    //Icon pack result
    iconPackCategory := reply.RegisterCategory("icon-packs", "Icon Packs", "", fmt.Sprintf(searchCategoryTemplate, "grid", "vertical", "small"))

    result := scopes.NewCategorisedResult(iconPackCategory)
    result.SetURI("scope://falcon.bhdouglass_falcon?q=icon-packs")
    result.SetTitle("Find Icon Packs")
    result.SetArt(falcon.getIcon("find-icon-packs", falcon.base.ScopeDirectory() + "/icon-packs.svg"))
    result.Set("type", "icon-packs")
    result.SetInterceptActivation()

    if err := reply.Push(result); err != nil {
        log.Fatalln(err)
    }

    return nil
}

func (falcon *Falcon) appPerformAction(result *scopes.Result, metadata *scopes.ActionMetadata, widgetId, actionId string) *scopes.ActivationResponse {
    var resp *scopes.ActivationResponse

    var app Application
    if err := result.Get("app", &app); err != nil {
        log.Println(err)
    }

    if actionId == "favorite" {
        if app.Id != "" {
            falcon.favorite(app.Id)
        }

        //redirect to blank search
        query := scopes.NewCannedQuery("falcon.bhdouglass_falcon", "", "")
        resp = scopes.NewActivationResponseForQuery(query)
    } else if actionId == "unfavorite" {
        if app.Id != "" {
            falcon.unfavorite(app.Id)
        }

        //redirect to blank search
        query := scopes.NewCannedQuery("falcon.bhdouglass_falcon", "", "")
        resp = scopes.NewActivationResponseForQuery(query)
    } else { //action is launch
        if app.IsApp {
            resp = scopes.NewActivationResponse(scopes.ActivationNotHandled)
        } else {
            query := scopes.NewCannedQuery(app.Id, "", "")
            resp = scopes.NewActivationResponseForQuery(query)
        }
    }

    return resp
}

func (falcon *Falcon) appActivate(result *scopes.Result, metadata *scopes.ActionMetadata) *scopes.ActivationResponse {
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

    return resp
}
