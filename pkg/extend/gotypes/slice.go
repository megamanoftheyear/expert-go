package gotypes

func SliceOf[Src, Dst any](src []Src, convert func(Src) Dst) (dst []Dst) {
	for _, element := range src {
		dst = append(dst, convert(element))
	}

	return dst
}
