package logger

func String(key, val string) Field {
	return Field{
		Key: key,
		Val: val,
	}
}

func Error(key string, val error) Field {
	return Field{
		Key: key,
		Val: val,
	}
}
