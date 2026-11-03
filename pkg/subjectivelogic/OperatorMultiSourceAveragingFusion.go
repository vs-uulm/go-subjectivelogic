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
MultiSourceAveragingFusion fuses an arbitrary number of opinions with Averaging Belief Fusion for multiple sources as
defined by Jøsang, Wang and Zhang, "Multi-Source Fusion in Subjective Logic", FUSION 2017 (Eq. 32-35). It corresponds to
averaging the evidence of all sources. The fused base rate is the average base rate.
Averaging fusion is not associative, so fusing three or more opinions by chaining AveragingFusion gives order-dependent
results. If any opinion is dogmatic, the result is the average of the dogmatic opinions (including their base rates)
with equal weights, as in Case II of Eq. 39 of Jøsang, "Categories of Belief Fusion", Journal of Advances in
Information Fusion 13(2), 2018. For a single dogmatic opinion, the base rate therefore differs from AveragingFusion,
which averages the base rates of both opinions.
For two opinions the result equals that of AveragingFusion apart from that base rate, and a single opinion is returned
unchanged.
*/
func MultiSourceAveragingFusion(opinions []Opinion) (Opinion, error) {
	if err := checkMultiSourceInput(opinions, "MultiSourceAveragingFusion"); err != nil {
		return Opinion{}, err
	}
	if len(opinions) == 1 {
		return opinions[0], nil
	}

	sums := computeMultiSourceSums(opinions)
	if sums.dogmatic > 0 {
		return NewOpinionWithTolerance(sums.dogmaticB, sums.dogmaticD, 0, sums.dogmaticA)
	}
	a := sums.sumA / float64(sums.n)
	return NewOpinionWithTolerance(sums.b/sums.s, sums.d/sums.s, float64(sums.n)*sums.uProd/sums.s, a)
}
