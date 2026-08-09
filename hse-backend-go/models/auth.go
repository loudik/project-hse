package models

// User - internal representation of a user row from the database (used
// between handlers, NOT the shape sent to the frontend - see EcmeUser)
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	RoleID   int    `json:"roleId"`
	RoleName string `json:"roleName"`
	Status   string `json:"status"`
}

// EcmeUser - shape expected by Ecme's session store (useSessionUser)
type EcmeUser struct {
	UserName  string   `json:"userName"`
	Email     string   `json:"email"`
	Avatar    string   `json:"avatar"`
	Authority []string `json:"authority"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse - persis bentuk yang dipakai AuthProvider.jsx:
// const resp = await apiSignIn(values)
// handleSignIn({ accessToken: resp.token }, resp.user)
type LoginResponse struct {
	Token string   `json:"token"`
	User  EcmeUser `json:"user"`
}

// MenuNode - satu node menu, bisa punya children (tree)
type MenuNode struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	Path     *string    `json:"path"`
	Icon     *string    `json:"icon"`
	Children []MenuNode `json:"children,omitempty"`
}
