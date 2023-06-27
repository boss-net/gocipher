package gocipher

import (
	"math"
)

// Gocipher cipher implementation from masscan

type GoCipher struct {
	Rounds int64
	Seed   int64
	Range  int64
	A      int64
	B      int64
}

func New(rangez, seed int64) *GoCipher {
	split := int64(math.Floor(math.Sqrt(float64(rangez))))
	var gocipher GoCipher
	gocipher.Rounds = 3
	gocipher.Seed = seed
	gocipher.Range = rangez
	gocipher.A = split - 1
	gocipher.B = split + 1

	if gocipher.A <= 0 {
		gocipher.A = 1
	}

	for gocipher.A*gocipher.B <= rangez {
		gocipher.B++
	}

	return &gocipher
}

// Inner permutation function
func (gocipher *GoCipher) F(j, r, seed int64) int64 {
	var primes = []int64{961752031, 982324657, 15485843, 961752031}
	r = (r << (r & 0x4)) + r + seed
	return int64(math.Abs(float64((((primes[j]*r + 25) ^ r) + j))))
}

// Outer feistal construction
func (gocipher *GoCipher) Fe(r, a, b, m, seed int64) int64 {
	var (
		L, R int64
		j    int64
		tmp  int64
	)

	L = m % a
	R = m / a

	for j = 1; j <= r; j++ {
		if j&1 == 1 {
			tmp = (L + gocipher.F(j, R, seed)) % a
		} else {
			tmp = (L + gocipher.F(j, R, seed)) % b
		}
		L = R
		R = tmp
	}

	if r&1 == 1 {
		return a*L + R
	}
	return a*R + L
}

// Outer reverse feistal construction
func (gocipher *GoCipher) Unfe(r, a, b, m, seed int64) int64 {
	var (
		L, R int64
		j    int64
		tmp  int64
	)

	if r&1 == 1 {
		R = m % a
		L = m / a
	} else {
		L = m % a
		R = m / a
	}

	for j = r; j >= 1; j-- {
		if j&1 == 1 {
			tmp = gocipher.F(j, L, seed)
			if tmp > R {
				tmp -= -R
				tmp = a - (tmp % a)
				if tmp == a {
					tmp = 0
				}
			} else {
				tmp = R - tmp
				tmp %= a
			}
		} else {
			tmp = gocipher.F(j, L, seed)
			if tmp > R {
				tmp = (tmp - R)
				tmp = b - (tmp % b)
				if tmp == b {
					tmp = 0
				}
			} else {
				tmp = R - tmp
				tmp %= b
			}
		}
		R = L
		L = tmp
	}

	return a*R + L
}

func (gocipher *GoCipher) Shuffle(m int64) int64 {
	c := gocipher.Fe(gocipher.Rounds, gocipher.A, gocipher.B, m, gocipher.Seed)

	for c >= gocipher.Range {
		c = gocipher.Fe(gocipher.Rounds, gocipher.A, gocipher.B, c, gocipher.Seed)
	}

	return c
}

func (gocipher *GoCipher) UnShuffle(m int64) int64 {
	c := gocipher.Unfe(gocipher.Rounds, gocipher.A, gocipher.B, m, gocipher.Seed)
	for c >= gocipher.Range {
		c = gocipher.Unfe(gocipher.Rounds, gocipher.A, gocipher.B, c, gocipher.Seed)
	}

	return c
}
