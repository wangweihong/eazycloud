package stringutil

func BothEmptyOrNone(str1, str2 string) bool {
	return (str1 == "" && str2 == "") || (str1 != "" && str2 != "")
}
