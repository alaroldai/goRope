package rope

import "math"

func phi() float64 {
	return (math.Sqrt(5) + 1) / 2
}

func inphi() float64 {
	return 1 / phi()
}

func fibIndex(fibNum uint) uint {
	return uint(math.Floor(math.Log(float64(fibNum)*math.Sqrt(5)+0.5) / math.Log(phi())))
}

func fib(idx uint) uint {
	var n float64 = float64(idx)
	return uint((math.Pow(phi(), n) - math.Pow(-1*inphi(), n)) / math.Sqrt(5))
}
