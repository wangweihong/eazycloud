package httpcli

func SuccessStatus(status int) bool {
	return status >= 200 && status <= 399
}
