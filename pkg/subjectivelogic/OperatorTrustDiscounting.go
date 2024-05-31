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
	"errors"
)

func TrustDiscounting(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error) {
	// Checking if the opinion pointers are empty
	if opinion1 == nil || opinion2 == nil {
		return Opinion{}, errors.New("OpTrustDisc: Input cannot be nil")
	}

	// Checking if the opinion values are null values
	nullChecker := Opinion{belief: 0, disbelief: 0, uncertainty: 0, baseRate: 0}
	if *opinion1 == nullChecker || *opinion2 == nullChecker {
		return Opinion{}, errors.New("Addition: Inputs cannot be null opinions")
	}

	b1 := opinion1.belief
	u1 := opinion1.uncertainty
	a1 := opinion1.baseRate

	b2 := opinion2.belief
	d2 := opinion2.disbelief
	a2 := opinion2.baseRate

	p1 := b1 + u1*a1
	b := p1 * b2
	d := p1 * d2
	u := 1 - b - d
	a := a2

	return NewOpinion(b, d, u, a)
}

func MultiEdgeTrustDisc(opinions []Opinion) (Opinion, error) {

	if opinions == nil {
		return Opinion{}, errors.New("MultiEdgeTrustDisc: Input cannot be nil")
	}
	n := len(opinions)
	if n < 2 {
		return Opinion{}, errors.New("MultiEdgeTrustDisc: At least two Opinions required")
	}

	P_acc := 1.0
	for i := 0; i < (n - 1); i++ {
		P_acc *= opinions[i].ProjProb()
	}

	nth_Opinion := opinions[n-1]
	b := P_acc * nth_Opinion.belief
	d := P_acc * nth_Opinion.disbelief
	u := 1 - b - d
	a := nth_Opinion.baseRate

	return NewOpinion(b, d, u, a)
}
