package utils

func DecideRole(requesterRole string, newUserRole string) bool {
	if requesterRole == "admin" {
		return true
	}else if requesterRole == "staff" && newUserRole == "customer" {
		return true
	}else{
		return false
	}
}