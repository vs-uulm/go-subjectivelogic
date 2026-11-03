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

import "errors"

// checkMultiSourceInput checks that at least one opinion is given and that no opinion is a null opinion
func checkMultiSourceInput(opinions []Opinion, operator string) error {
	if len(opinions) == 0 {
		return errors.New(operator + ": At least one Opinion required")
	}
	nullChecker := Opinion{belief: 0, disbelief: 0, uncertainty: 0, baseRate: 0}
	for _, opinion := range opinions {
		if opinion == nullChecker {
			return errors.New(operator + ": Inputs cannot be null opinions")
		}
	}
	return nil
}
