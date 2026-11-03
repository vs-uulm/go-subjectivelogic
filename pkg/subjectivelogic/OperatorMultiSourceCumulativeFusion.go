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

/*
MultiSourceCumulativeFusion fuses an arbitrary number of opinions with Cumulative Belief Fusion for multiple sources as
defined by Jøsang, Wang and Zhang, "Multi-Source Fusion in Subjective Logic", FUSION 2017 (Eq. 16-19). It corresponds to
summing the evidence of all sources. The fused base rate is the evidence-weighted average base rate as given by Jøsang,
"Categories of Belief Fusion", Journal of Advances in Information Fusion 13(2), 2018 (Eq. 23).
If any opinion is dogmatic, the result is the average of the dogmatic opinions (including their base rates) with equal
weights, as in Case II of Eq. 24 of "Categories of Belief Fusion".
For two opinions the result equals that of CumulativeFusion, and a single opinion is returned unchanged.
*/
func MultiSourceCumulativeFusion(opinions []Opinion) (Opinion, error) {
	if err := checkMultiSourceInput(opinions, "MultiSourceCumulativeFusion"); err != nil {
		return Opinion{}, err
	}
	if len(opinions) == 1 {
		return opinions[0], nil
	}

	sums := computeMultiSourceSums(opinions)
	if sums.dogmatic > 0 {
		return NewOpinionWithTolerance(sums.dogmaticB, sums.dogmaticD, 0, sums.dogmaticA)
	}
	// sw + uProd equals sum_C P_C - (n - 1) U, the denominator of Eq. 16
	denominator := sums.sw + sums.uProd
	a := sums.sumA / float64(sums.n)
	if sums.sw > 0 { // otherwise all opinions are vacuous and the base rates are averaged
		a = sums.aw / sums.sw
	}
	return NewOpinionWithTolerance(sums.b/denominator, sums.d/denominator, sums.uProd/denominator, a)
}
