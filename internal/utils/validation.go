package utils


func ValidateRegistration(password, confirm_password string) (map[string]string, bool) {
    result := make(map[string]string)
    valid := true
    if password != confirm_password {
        valid = false
        result["password"] = "Passwords do not match"
        result["confirm_password"] = "Passwords do not match"
    }
    if len(password) < 6 {
        valid = false
        result["password"] = "Password too short"
    }
    return result, valid
}
