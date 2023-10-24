package keycloak

const (
	// keycloak token claim
	ClaimSub               = "sub"
	ClaimEmail             = "email"
	ClaimPreferredUsername = "preferred_username"
	ClaimName              = "name"

	// realm list
	MasterRealmId = "master"

	// token grant type list
	GrantTypePassword          = "password"
	GrantTypeClientCredentials = "client_credentials"

	// keycloak errors
	ErrPasswordMustNotEqualWithLastPassword = "400 Bad Request: invalidPasswordHistoryMessage: Invalid password: must not be equal to any of last 3 passwords."
	ErrUserIsDisabled                       = "400 Bad Request: User is disabled"
	ErrInvalidBearerToken                   = "401 Unauthorized: invalid_grant: Invalid bearer token"
	ErrInvalidUserCredentials               = "401 Unauthorized: invalid_grant: Invalid user credentials"
	ErrUnknown                              = "403 Forbidden: unknown_error"
	ErrAccessDeniedNotAuthorized            = "403 Forbidden: access_denied: not_authorized"
	ErrCouldNotObtainAccessToken            = "403 Forbidden: invalid_bearer_token: Could not obtain bearer access_token from request."
	ErrRealmNotFound                        = "404 Not Found: Realm not found."
	ErrRealmDoesntExist                     = "404 Not Found: Realm does not exist"
	ErrUserNotFound                         = "404 Not Found: User not found"
	ErrRoleNotFound                         = "404 Not Found: Role not found"
	ErrCouldNotFindRole                     = "404 Not Found: Could not find role"
	ErrCouldNotFindRoleWithId               = "404 Not Found: Could not find role with id"
	ErrCouldNotFindClient                   = "404 Not Found: Could not find client"
	ErrUsernameAlreadyInUsed                = "409 Conflict: User exists with same username"
	ErrEmailAlreadyInUsed                   = "409 Conflict: User exists with same email"

	// other errors
	ErrtokenExpired             = "\"exp\" not satisfied"
	ErrFailedToFindKeyWithKeyID = "failed to find key with key ID"

	OpenIdConfigPath = "/realms/%s/.well-known/openid-configuration"
)
