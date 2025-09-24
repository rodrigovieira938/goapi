package permission

type Permission struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UserPermission struct {
	UserID       int `json:"user_id"`
	PermissionID int `json:"permission_id"`
}
