package douyu

func number2bytes(num int64) (bs []byte) {
	var (
		tmpBarr [64]byte
		cnt     int
	)
	for num > 0 {
		tmpBarr[cnt] = byte(num%10 + '0')
		cnt++
		num /= 10
	}
	bs = make([]byte, cnt)
	for i := cnt; i > 0; i-- {
		(bs)[cnt-i] = tmpBarr[i-1]
	}
	return
}

func bytes2number(data []byte) (number int64) {
	for _, v := range data {
		number = number*10 + int64(v-'0')
	}
	return
}
