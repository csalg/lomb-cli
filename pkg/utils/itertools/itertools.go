package itertools

func Chunk[T any](s []T, size int, callback func([]T, int) bool) {
	for i := 0; i < len(s); i += size {
		end := i + size
		if len(s) < end {
			end = len(s)
		}
		if !callback(s[i:end], i) {
			break
		}
	}
}
