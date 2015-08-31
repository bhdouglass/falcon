/*
Package scopes is used to write Unity scopes in Go.

Scopes are implemented through types that conform to the Scope interface.

    type MyScope struct {}

The scope has access to a ScopeBase instance, which can be used to access
various pieces of information about the scope, such as its settings:

    func (s *MyScope) SetScopeBase(base *scopes.ScopeBase) {
    }

If the scope needs access to this information, it should save the
provided instance for later use.  Otherwise, it the method body can be
left blank.

The shell may ask the scope for search results, which will cause the
Search method to be invoked:

    func (s *MyScope) Search(query *scopes.CannedQuery, metadata *scopes.SearchMetadata, reply *scopes.SearchReply, cancelled <-chan bool) error {
        category := reply.RegisterCategory("cat_id", "category", "", "")
        result := scopes.NewCategorisedResult(category)
        result.SetTitle("Result for " + query.QueryString())
        reply.Push(result)
        return nil
    }

In general, scopes will:

* Register result categories via reply.RegisterCategory()

* Create new results via NewCategorisedResult(), and push them with reply.Push(result)

* Check for cancellation requests via the provided channel.

The Search method will be invoked with an empty query when sufacing
results are wanted.

The shell may ask the scope to provide a preview of a result, which causes the Preview method to be invoked:

    func (s *MyScope) Preview(result *scopes.Result, metadata *scopes.ActionMetadata, reply *scopes.PreviewReply, cancelled <-chan bool) error {
        widget := scopes.NewPreviewWidget("foo", "text")
        widget.AddAttributeValue("text", "Hello")
        reply.PushWidgets(widget)
        return nil
    }

The scope should push one or more slices of PreviewWidgets using reply.PushWidgets.  PreviewWidgets can be created with NewPreviewWidget.

Additional data for the preview can be pushed with reply.PushAttr.

If any of the preview widgets perform actions that the scope should
respond to, the scope should implement the PerformAction method:

    func (s *MyScope) PerformAction(result *Result, metadata *ActionMetadata, widgetId, actionId string) (*ActivationResponse, error) {
        // handle the action and then tell the dash what to do next
        // through an ActivationResponse.
        resp := NewActivationResponse(ActivationHideDash) return resp, nil
    }

The PerformAction method is not part of the main Scope interface, so
the feature need only be implemented for scopes that use the feature.

Finally, the scope can be exported in the main function:

    func main() {
        if err := scopes.Run(&MyScope{}); err != nil {
            log.Fatalln(err)
        }
    }

The scope executable can be deployed to a scope directory named like:

    /usr/lib/${arch}/unity-scopes/${scope_name}

In addition to the scope executable, a scope configuration file named
${scope_name}.ini should be placed in the directory.  Its contents
should look something like:

    [ScopeConfig]
    DisplayName = Short name for the scope
    Description = Long description of scope
    Author =
    ScopeRunner = ${scope_executable} --runtime %R --scope %S
*/
package scopes
