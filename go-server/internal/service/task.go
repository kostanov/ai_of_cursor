package service

// SlowTask — тяжёлая CPU-задача (аналог _slow_task в api.py).
func SlowTask(iterations int) int {
	x := 0
	for i := 0; i < iterations; i++ {
		x += i * (i + 1) / 2
	}
	return x
}
