package tinygo

func EnsureInstalled(opts ...Option) (string, error) {
	c := newConfig(opts...)

	p, err := getPath(c)
	if err == nil {
		return p, nil
	}

	if err := Install(opts...); err != nil {
		return "", err
	}

	return getPath(c)
}
