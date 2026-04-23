package main

func LinearKernel() func([]float64, []float64) float64 {
	return func(X []float64, x []float64) float64 {
		// don't throw error but fail peacefully
		//
		// returning "not at all similar", basically
		if len(X) != len(x) {
			return 0.0
		}

		var dot float64

		for i := range X {
			dot += X[i] * x[i]
		}

		return dot
	}
}
