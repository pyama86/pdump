package main

type counter struct {
	current uint
	sums    []uint
	len     uint
	capa    uint
}

func (c *counter) resetAll() {
	c.current = 0
	c.sums = []uint{}
	c.len = 0
}

func (c *counter) reset() {
	c.current = 0
}

func (c *counter) avg() uint {
	var s uint
	for _, n := range c.sums {
		s += n
	}
	if s == 0 {
		return 0
	}
	return s / c.len
}

func (c *counter) increment() {
	c.current++
}

func (c *counter) included() {
	if c.current == 0 {
		c.resetAll()
		return
	}
	c.sums = append(c.sums, c.current)
	c.len++
	if c.len > c.capa {
		c.sums = c.sums[1:]
		c.len--
	}
}
