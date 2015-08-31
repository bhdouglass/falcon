#include <unity/scopes/Category.h>
extern "C" {
#include "_cgo_export.h"
}
#include "scope.h"
#include "smartptr_helper.h"

using namespace unity::scopes;

ScopeAdapter::ScopeAdapter(GoInterface goscope) : goscope(goscope) {
}

void ScopeAdapter::start(std::string const &) {
    setScopeBase(goscope, reinterpret_cast<_ScopeBase*>(this));
}

void ScopeAdapter::stop() {
    setScopeBase(goscope, nullptr);
}

SearchQueryBase::UPtr ScopeAdapter::search(CannedQuery const &q,
                                     SearchMetadata const &metadata) {
    SearchQueryBase::UPtr query(new QueryAdapter(q, metadata, *this));
    return query;
}

PreviewQueryBase::UPtr ScopeAdapter::preview(Result const& result, ActionMetadata const& metadata) {
    PreviewQueryBase::UPtr query(new PreviewAdapter(result, metadata, *this));
    return query;
}

ActivationQueryBase::UPtr ScopeAdapter::activate(Result const& result, ActionMetadata const &metadata) {
    ActivationQueryBase::UPtr activation(new ActivationAdapter(result, metadata, *this));
    return activation;
}

ActivationQueryBase::UPtr ScopeAdapter::perform_action(Result const& result, ActionMetadata const &metadata, std::string const &widget_id, std::string const &action_id) {
    ActivationQueryBase::UPtr activation(new ActivationAdapter(result, metadata, widget_id, action_id, *this));
    return activation;
}

QueryAdapter::QueryAdapter(CannedQuery const &query,
                           SearchMetadata const &metadata,
                           ScopeAdapter &scope)
    : SearchQueryBase(query, metadata), scope(scope),
      cancel_channel(makeCancelChannel(), releaseCancelChannel) {
}

void QueryAdapter::cancelled() {
    sendCancelChannel(cancel_channel.get());
}

void QueryAdapter::run(SearchReplyProxy const &reply) {
    callScopeSearch(
        scope.goscope,
        reinterpret_cast<_CannedQuery*>(new CannedQuery(query())),
        reinterpret_cast<_SearchMetadata*>(new SearchMetadata(search_metadata())),
        const_cast<uintptr_t*>(reinterpret_cast<const uintptr_t*>(&reply)),
        cancel_channel.get());
}

PreviewAdapter::PreviewAdapter(Result const &result,
                               ActionMetadata const &metadata,
                               ScopeAdapter &scope)
    : PreviewQueryBase(result, metadata), scope(scope),
      cancel_channel(makeCancelChannel(), releaseCancelChannel) {
}

void PreviewAdapter::cancelled() {
    sendCancelChannel(cancel_channel.get());
}

void PreviewAdapter::run(PreviewReplyProxy const &reply) {
    callScopePreview(
        scope.goscope,
        reinterpret_cast<_Result*>(new Result(result())),
        reinterpret_cast<_ActionMetadata*>(new ActionMetadata(action_metadata())),
        const_cast<uintptr_t*>(reinterpret_cast<const uintptr_t*>(&reply)),
        cancel_channel.get());
}

ActivationAdapter::ActivationAdapter(Result const &result,
                                     ActionMetadata const &metadata,
                                     ScopeAdapter &scope)
    : ActivationQueryBase(result, metadata), scope(scope),
      is_action(false) {
}

ActivationAdapter::ActivationAdapter(Result const &result,
                                     ActionMetadata const &metadata,
                                     std::string const &widget_id,
                                     std::string const &action_id,
                                     ScopeAdapter &scope)
    : ActivationQueryBase(result, metadata, widget_id, action_id),
      scope(scope), is_action(true) {
}

ActivationResponse ActivationAdapter::activate() {
    ActivationResponse response(ActivationResponse::NotHandled);
    char *error = nullptr;
    if (is_action) {
        callScopePerformAction(
            scope.goscope,
            reinterpret_cast<_Result*>(new Result(result())),
            reinterpret_cast<_ActionMetadata*>(new ActionMetadata(action_metadata())),
            const_cast<char*>(widget_id().c_str()),
            const_cast<char*>(action_id().c_str()),
            reinterpret_cast<_ActivationResponse*>(&response),
            &error);
    } else {
        callScopeActivate(
            scope.goscope,
            reinterpret_cast<_Result*>(new Result(result())),
            reinterpret_cast<_ActionMetadata*>(new ActionMetadata(action_metadata())),
            reinterpret_cast<_ActivationResponse*>(&response),
            &error);
    }
    if (error != nullptr) {
        const std::string message(error);
        free(error);
        throw std::runtime_error(message);
    }
    return response;
}
