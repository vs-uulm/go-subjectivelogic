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
MultiSourceCumulativeEpistemicFusion fuses an arbitrary number of opinions with Epistemic Cumulative Belief Fusion for
multiple sources: the opinions are fused with MultiSourceCumulativeFusion, and the result is uncertainty-maximised once.
This is CBF-UM as described by Jøsang, "Categories of Belief Fusion", Journal of Advances in Information Fusion 13(2),
2018 (Sections 6.2 and 6.5), and as argued for by van der Heijden, Kopp and Kargl, "Multi-Source Fusion Operations in
Subjective Logic", FUSION 2018 (Section III-A). Chaining CumulativeEpistemicFusion instead maximises the uncertainty after every pairwise
fusion, which gives order-dependent results.
For two opinions the result equals that of CumulativeEpistemicFusion. A single opinion is uncertainty-maximised, so that
adding a vacuous opinion, which is neutral in cumulative fusion, does not change the result.
*/
func MultiSourceCumulativeEpistemicFusion(opinions []Opinion) (Opinion, error) {
	if err := checkMultiSourceInput(opinions, "MultiSourceCumulativeEpistemicFusion"); err != nil {
		return Opinion{}, err
	}
	fused, err := MultiSourceCumulativeFusion(opinions)
	if err != nil {
		return Opinion{}, err
	}
	return uncertaintyMaximised(fused, "MultiSourceCumulativeEpistemicFusion")
}
