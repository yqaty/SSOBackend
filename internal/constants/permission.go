package constants

const (
	GrantUserRoleAction    = "Grant"
	RevokeUserRoleAction   = "Revoke"
	ChangeUserRoleResource = "UniqueSSO::Internal::User::Role"

	AddRoleAction      = "ADD"
	DeleteRoleAction   = "DELETE"
	ChangeRoleResource = "UniqueSSO::Internal::Role"

	AddObjectAction      = "ADD"
	RemoveObjectAction   = "DELETE"
	ChangeObjectResource = "UniqueSSO::Internal::Object"

	AddObjectGroupAction      = "ADD"
	RemoveObjectGroupAction   = "DELETE"
	ChangeObjectGroupResource = "UniqueSSO::Internal::ObjectGroup"

	GrantRoleObjectGroupAction    = "Grant"
	RevokeRoleObjectGroupAction   = "Revoke"
	ChangeRoleObjectGroupResource = "UniqueSSO::Internal::Role::ObjectGroup"
)
