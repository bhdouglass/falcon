#ifndef UNITYSCOPE_SCOPE_H
#define UNITYSCOPE_SCOPE_H

#include <memory>
#include <string>

#include <unity/scopes/ScopeBase.h>

class ScopeAdapter : public unity::scopes::ScopeBase
{
    friend class QueryAdapter;
    friend class PreviewAdapter;
    friend class ActivationAdapter;
public:
    ScopeAdapter(GoInterface goscope);
    virtual void start(std::string const&) override;
    virtual void stop() override;
    virtual unity::scopes::SearchQueryBase::UPtr search(unity::scopes::CannedQuery const &query, unity::scopes::SearchMetadata const &metadata) override;

    virtual unity::scopes::PreviewQueryBase::UPtr preview(unity::scopes::Result const& result, unity::scopes::ActionMetadata const& metadata) override;
    virtual unity::scopes::ActivationQueryBase::UPtr activate(unity::scopes::Result const& result, unity::scopes::ActionMetadata const &metadata) override;
    virtual unity::scopes::ActivationQueryBase::UPtr perform_action(unity::scopes::Result const& result, unity::scopes::ActionMetadata const &metadata, std::string const &widget_id, std::string const &action_id) override;

private:
    GoInterface goscope;
};

class QueryAdapter : public unity::scopes::SearchQueryBase
{
public:
    QueryAdapter(unity::scopes::CannedQuery const &query,
                 unity::scopes::SearchMetadata const &metadata,
                 ScopeAdapter &scope);
    virtual void cancelled() override;
    virtual void run(unity::scopes::SearchReplyProxy const &reply) override;
private:
    const ScopeAdapter &scope;
    std::unique_ptr<void, void(*)(GoChan)> cancel_channel;
};

class PreviewAdapter : public unity::scopes::PreviewQueryBase
{
public:
    PreviewAdapter(unity::scopes::Result const &result,
                   unity::scopes::ActionMetadata const &metadata,
                   ScopeAdapter &scope);
    virtual void cancelled() override;
    virtual void run(unity::scopes::PreviewReplyProxy const &reply) override;
private:
    const ScopeAdapter &scope;
    std::unique_ptr<void, void(*)(GoChan)> cancel_channel;
};

class ActivationAdapter : public unity::scopes::ActivationQueryBase
{
public:
    ActivationAdapter(unity::scopes::Result const &result,
                      unity::scopes::ActionMetadata const &metadata,
                      ScopeAdapter &scope);
    ActivationAdapter(unity::scopes::Result const &result,
                      unity::scopes::ActionMetadata const &metadata,
                      std::string const &widget_id,
                      std::string const &action_id,
                      ScopeAdapter &scope);
    virtual unity::scopes::ActivationResponse activate() override;
private:
    const ScopeAdapter &scope;
    bool is_action;
};

#endif
