package pkg

func ValidateGender(g string) bool {
	if g == "male" || g == "female" {
		return true
	}
	return false
}
