package main

type Calculator struct {
	id  int
	sum int
}

func (c *Calculator) CalculateSum(start, end int) {
	for i := start; i <= end; i++ {
		c.sum += i
	}
}
