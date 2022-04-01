package sequence

// Task: Fibonacci numbers
//
// fibonacci(n) returns the n-th Fibonacci number, and is defined by the
// recurrence relation F_n = F_n-1 + F_n-2, with seed values F_0=0 and F_1=1.
func fibonacci(n uint) uint {
	var f0 uint = 0
	var f1 uint = 1
	var i uint

	for i = 0; i < n; i++ {
		temp := f0
		f0 = f1
		f1 = temp + f0
	}
	return f0
}
