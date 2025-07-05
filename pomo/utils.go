package pomo

func dictLike(dict map[any]any, key, elseCond any) any {
	ans, ok := dict[key]
	if !ok {
		ans = elseCond
	}
	return ans
}
