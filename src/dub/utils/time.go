package utils

//字符串转化为时间，必须是2006-01-02 15:04:05.999999格式
func TimeStrtoTime(str string) (t time.Time, err error) {
	timeFormat := "2006-01-02 15:04:05.999999"

	switch len(str) {
	case 10, 19, 21, 22, 23, 24, 25, 26: //up to "YYYY-MM-DD HH:MM:SS.MMMMMM"
		t, err = time.ParseInLocation(timeFormat[:len(str)], str, time.Local)
	default:
		err = fmt.Errorf("invalid time string: %s", str)
		return
	}

	return
}
