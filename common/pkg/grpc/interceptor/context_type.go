package interceptor

type ContextKey string

const (
	ContextUserID ContextKey = "user_id"
	ContextRole   ContextKey = "role"

	ContextAdminID   ContextKey = "admin_id"
	ContextAdminRole ContextKey = "admin_role"

	ContextPrincipal ContextKey = "principal_type"
	ContextJWTToken  ContextKey = "jwt_token"
)

const (
	PrincipalUser  = "user"
	PrincipalAdmin = "admin"
)
