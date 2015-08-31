package scopes

type PostLoginAction int

const (
	_ = PostLoginAction(iota)
	PostLoginDoNothing
	PostLoginInvalidateResults
	PostLoginContinueActivation
)

type accountDetails struct {
	ScopeID           string          `json:"scope_id"`
	ServiceName       string          `json:"service_name"`
	ServiceType       string          `json:"service_type"`
	ProviderName      string          `json:"provider_name"`
	LoginPassedAction PostLoginAction `json:"login_passed_action"`
	LoginFailedAction PostLoginAction `json:"login_failed_action"`
}

// RegisterAccountLoginResult configures a result such that the dash
// will attempt to log in to the account identified by (serviceName,
// serviceType, providerName).
//
// On success, the dash will perform the action specified by
// passedAction.  On failure, it will use failedAction.
func RegisterAccountLoginResult(result *CategorisedResult, query *CannedQuery, serviceName, serviceType, providerName string, passedAction, failedAction PostLoginAction) error {
	if result.URI() == "" {
		if err := result.SetURI(query.ToURI()); err != nil {
			return err
		}
	}
	return result.Set("online_account_details", accountDetails{
		ScopeID:           query.ScopeID(),
		ServiceName:       serviceName,
		ServiceType:       serviceType,
		ProviderName:      providerName,
		LoginPassedAction: passedAction,
		LoginFailedAction: failedAction,
	})
}
