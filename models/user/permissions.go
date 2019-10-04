package user

// CanViewUser returns true if THIS user can view USERNAME's details.
//
// Admins can view/edit all users. All others can only view/edit themselves.
func (u *User) CanViewUser(username string) bool {
	if u.IsAdmin() {
		return true
	}
	if u.Username == username {
		return true
	}
	return false
}

// CanEditUser returns true if THIS user can edit USERNAME's details.
//
// Admins can view/edit all users. All others can only view/edit themselves.
func (u *User) CanEditUser(username string) bool {
	if u.IsAdmin() {
		return true
	}
	if u.Username == username {
		return true
	}
	return false
}
