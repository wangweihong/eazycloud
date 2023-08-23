package errorutil

func ErrorMsg(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
