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
MultiSourceWeightedFusion fuses an arbitrary number of opinions with Weighted Belief Fusion for multiple sources as
defined by van der Heijden, Kopp and Kargl, "Multi-Source Fusion Operations in Subjective Logic", FUSION 2018, which is
also implemented by subjective-logic-java (weightedCollectionFuse). It corresponds to the confidence-weighted average of
the evidence of all sources, where the confidence of a source is 1 - u. The fused base rate is the confidence-weighted
average base rate.
Weighted fusion is not associative, so fusing three or more opinions by chaining WeightedFusion gives order-dependent
results. If any opinion is dogmatic, the result is the average of the dogmatic opinions (including their base rates)
with equal weights, as in Case II of Eq. 56 of Jøsang, "Categories of Belief Fusion", Journal of Advances in
Information Fusion 13(2), 2018. For a single dogmatic opinion, the base rate therefore differs from WeightedFusion,
which uses the confidence-weighted base rate of both opinions. If all
opinions are vacuous, the result is vacuous with the average base rate.
For two opinions the result equals that of WeightedFusion apart from that base rate, and a single opinion is returned
unchanged.
*/
func MultiSourceWeightedFusion(opinions []Opinion) (Opinion, error) {
	if err := checkMultiSourceInput(opinions, "MultiSourceWeightedFusion"); err != nil {
		return Opinion{}, err
	}
	if len(opinions) == 1 {
		return opinions[0], nil
	}

	sums := computeMultiSourceSums(opinions)
	n := float64(sums.n)
	if sums.dogmatic > 0 {
		return NewOpinionWithTolerance(sums.dogmaticB, sums.dogmaticD, 0, sums.dogmaticA)
	}
	if sums.vacuous == sums.n {
		return NewOpinionWithTolerance(0, 0, 1, sums.sumA/n)
	}
	a := sums.sumConfA / (n - sums.sumU)
	return NewOpinionWithTolerance(sums.bw/sums.sw, sums.dw/sums.sw, (n-sums.sumU)*sums.uProd/sums.sw, a)
}
