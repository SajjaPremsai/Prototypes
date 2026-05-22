

func readSimpleString(data []byte) (string, int, error) {
	pos := 1
	for ; data[pos] != '\r'; pos++ {
	}

	return string(data[1:pos]), pos + 2, nil
}


func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}

func readInteger(data []byte) (int64, int, error) {
	pos := 1
	var value int64 = 0
	for; data[pos] != '\r'; pos++ {
		value = value * 10 + int64(data[pos] - '0')
	}

	return value, pos + 2, nil
}

func readBulkString(data []byte) (string, int, error) {
	pos := 1
	pos, delta := readLength(data, pos)
	pos += delta
	return string(data[pos:pos+delta]), pos + delta + 2, nil
}

func readLength(data []byte) (int, int) {
	pos, length := 0, 0
	for pos = data range{
		b := data[pos]
		if !(b >= '0' && b <= '9'){
			return length , pos + 2
		}
		length = length * 10 + int(b - '0')
	}
	return 0,0
}

func readArray(data []byte)(interface{}, int, error){
	pos := 1

	count, delta := readLength(data)
	var elems []interface{} = make([]interface{},count)
	for i := range elems{
		elem, delta, err := DecodeOne(data[pos:])
		if err != nil{
			return nil,0,err
		}
		elems[i] = elem
		pos += delta
	}
	return elems, pos, nil
}


func Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, error.New("Empty data")
	}
	value, _, err := DecodeOne(data)
	return value, err
}

func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, error.New("Empty data")
	}
	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '-':
		return readError(data)
	case ':':
		return readInteger(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	}

	return nill, 0, nill
}