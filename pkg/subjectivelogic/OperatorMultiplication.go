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

func Multiplication(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error) {
	if opinion1 == nil || opinion2 == nil {
		return Opinion{}, errors.New("Multiplication: Input cannot be nil")
	}
	if opinion1.baseRate == 1 && opinion2.baseRate == 1 {
		return Opinion{}, errors.New("Multiplication: Base rates cannot both be 1")
	}

	b1 := opinion1.belief
	d1 := opinion1.disbelief
	u1 := opinion1.uncertainty
	a1 := opinion1.baseRate

	b2 := opinion2.belief
	d2 := opinion2.disbelief
	u2 := opinion2.uncertainty
	a2 := opinion2.baseRate

	b := b1*b2 + ((1-a1)*a2*b1*u2+a1*(1-a2)*u1*b2)/(1-a1*a2)
	d := d1 + d2 - d1*d2
	u := u1*u2 + ((1-a2)*b1*u2+(1-a1)*u1*b2)/(1-a1*a2)
	a := a1 * a2

	return NewOpinion(b, d, u, a)
}
