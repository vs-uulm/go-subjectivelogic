//Copyright 2024 Institute of Distributed Systems, Ulm University
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package subjectivelogic

import (
	"math"
)

/*
multiSourceSums contains the sums over the opinions of n sources C from which the multi-source cumulative, averaging and
weighted fusion operators are computed. With U = prod_C u_C and P_C = prod_{j != C} u_j (the uncertainty of all other
sources):

	uProd      = U
	s          = sum_C P_C
	b, d       = sum_C b_C * P_C,  sum_C d_C * P_C
	bw, dw, sw = sum_C b_C * (1 - u_C) * P_C,  sum_C d_C * (1 - u_C) * P_C,  sum_C (1 - u_C) * P_C
	aw         = sum_C a_C * (1 - u_C) * P_C

These sums are products with one factor per source, so they may be scaled by a common factor, and only their ratios are
used. They are rescaled by a power of two so that they do not underflow for many sources (only uncertainties that are
themselves subnormal numbers can still make s zero, which results in an error).
*/
type multiSourceSums struct {
	n                    int
	uProd, s, b, d       float64
	bw, dw, sw, aw       float64
	sumU, sumA, sumConfA float64 // sum_C u_C, sum_C a_C and sum_C a_C * (1 - u_C), not rescaled
	vacuous              int     // number of vacuous opinions
	dogmatic             int     // number of dogmatic opinions
	dogmaticB, dogmaticD float64 // mean belief and disbelief of the dogmatic opinions
	dogmaticA            float64 // mean base rate of the dogmatic opinions
}

// rescaleBelow is the value below which the sums are rescaled by a power of two
const rescaleBelow = 0x1p-500

// computeMultiSourceSums accumulates the sums in a single pass with the recurrence
// sum_{k+1} = sum_k * u_{k+1} + U_k * x_{k+1}, i.e. every existing summand gains the factor u_{k+1}.
func computeMultiSourceSums(opinions []Opinion) multiSourceSums {
	sums := multiSourceSums{n: len(opinions), uProd: 1}
	for _, opinion := range opinions {
		b, d, u, a := opinion.belief, opinion.disbelief, opinion.uncertainty, opinion.baseRate
		c := 1 - u

		sums.s = sums.s*u + sums.uProd
		sums.b = sums.b*u + sums.uProd*b
		sums.d = sums.d*u + sums.uProd*d
		sums.bw = sums.bw*u + sums.uProd*b*c
		sums.dw = sums.dw*u + sums.uProd*d*c
		sums.sw = sums.sw*u + sums.uProd*c
		sums.aw = sums.aw*u + sums.uProd*a*c
		sums.uProd *= u

		// s is the largest sum, since P_C >= U and all other factors are at most 1
		if sums.s > 0 && sums.s < rescaleBelow {
			_, exp := math.Frexp(sums.s)
			sums.uProd, sums.s = math.Ldexp(sums.uProd, -exp), math.Ldexp(sums.s, -exp)
			sums.b, sums.d = math.Ldexp(sums.b, -exp), math.Ldexp(sums.d, -exp)
			sums.bw, sums.dw = math.Ldexp(sums.bw, -exp), math.Ldexp(sums.dw, -exp)
			sums.sw, sums.aw = math.Ldexp(sums.sw, -exp), math.Ldexp(sums.aw, -exp)
		}

		sums.sumU += u
		sums.sumA += a
		sums.sumConfA += a * c
		if u == 1 {
			sums.vacuous++
		}
		if u == 0 {
			sums.dogmatic++
			sums.dogmaticB += b
			sums.dogmaticD += d
			sums.dogmaticA += a
		}
	}
	if k := float64(sums.dogmatic); k > 0 {
		sums.dogmaticB, sums.dogmaticD, sums.dogmaticA = sums.dogmaticB/k, sums.dogmaticD/k, sums.dogmaticA/k
	}
	return sums
}
