package logger

func String(key, val string) Field {
	return Field{
		Key: key,
		Val: val,
	}
}

func Error(err error) Field {
	return Field{
		Key: "error",
		Val: err,
	}
}
