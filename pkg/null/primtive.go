package null

func String(v string) *string {
	value := v
	return &value
}

func Int32(v int32) *int32 {
	value := v
	return &value
}
