# go-subjectivelogic



This Go library implements a number of Subjective Logic operators. Subjective Logic is a type of probabilistic logic that uses subjective opinions instead of probabilities. To learn more about Subjective Logic please refer to "Subjective Logic, A Formalism for Reasoning under Uncertainty" book by Audun Jøsang. 

## Table of Contents

- [Usage](#usage)
	* [Opinion](#opinion)
	* [Addition](#addition)
	* [Complement](#complement)
	* [Binomial Multiplication](#binomial-multiplication)
	* [Binomial Comultiplication](#binomial-comultiplication)
	* [Belief Constraint Fusion](#belief-constraint-fusion)
	* [Cumulative Fusion](#cumulative-fusion)
	* [Averaging Fusion](#averaging-fusion)
	* [Weighted Fusion](#weighted-fusion)
	* [Trust Discounting](#trust-discounting)
	* [Consecutive Trust Discounting](#consecutive-trust-discounting)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## Usage

### Opinion

Subjective Logic is a type of probabilistic logic that uses subjective opinions instead of probabilities. A binomial opinion about the truth of the proposition $x$ is $\omega_x = (b_x, d_x, u_x, a_x)$, where:

- $b_x$: belief mass in support of the proposition $x$ being true.
- $d_x$: disbelief mass in support of the proposition $x$ being false.
- $u_x$: uncertainty mass representing lack of evidence.
- $a_x$: base rate i.e. the prior probability of $x$ being true without any evidence.

```go
type Opinion struct {
	belief      float64
	disbelief   float64
	uncertainty float64
	baseRate    float64
}
```

The following methods have been implemented for the opinion:
- [NewOpinion()](#new-opinion)
- [GetBelief()]()
- [GetDisbelief()]()
- [GetUncertainty()]()
- [GetBaseRate()]()
- [Modify()]()
- [ProjProb()]()
- [Compare()]()
- [ToString]()
- [ToStringE()]()


#### NewOpinion() method

The `NewOpinion()` method takes in four `float64` values and outputs an `Opinion` as well as an Error. In case a valid `Opinion` can be formed, it will be returned and the error will be `nil`. If the input values violate the requirements for a valid `Opinion`, the `Opinion` would be initialized with zeroes and an error will be returned.
For a valid `Opinion`, each input value $i$ must fulfill $0 \leq i \leq 1$ and the following statement $b + d + u = 1$ must hold.

```go
opinion, err := subjectivelogic.NewOpinion(.5, .25, .25, .5)

if err != nil {
	println(err.Error())
} else {
	belief := opinion.GetBelief()
	disbelief := opinion.GetDisbelief()
	uncertainty := opinion.GetUncertainty()
	baseRate := opinion.GetBaseRate()
	
	println(fmt.Sprintf("Opinion: %.2f, %.2f, %.2f, %.2f", belief, disbelief, uncertainty, baseRate))
}
```

This code generates the following output:

```
Opinion: 0.50, 0.25, 0.25, 0.50
```
---

### Addition
This implements the Addition Operator as defined in Subjective Logic:

$$
\omega_{(x\cup y)} = 
\begin{cases}
b_{x\cup y} = b_x + b_y \\
d_{x\cup y} = \frac{a_x(d_x - b_y) + a_y(d_y - b_x)}{a_x+a_y} \\
u_{x\cup y} = \frac{a_xu_x + a_yu_y}{a_x+a_y} \\
a_{x\cup y} = 0
\end{cases}
$$

#### API Reference

```go
func Addition(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error)
```

#### Problematic Inputs
Inputs $\omega_{x} = (b_x, d_x, u_x, a_x)$ and $\omega_{y} = (b_y, d_y, u_y, a_y)$ are problematic and will result in an error, if:

$$
\begin{split}
b_x + b_y &> 1 \text{, or} \\
a_x + a_y &> 1 
\end{split}
$$

Moreover, the Addition Operator will attempt to divide through $0$, if $a_x = a_y = 0$, hence this is not allowed and an error will be returned.

#### Example

```go
func main() {

	opinion1, _ := subjectivelogic.NewOpinion(0.2, 0.8, 0, 0.5)
	opinion2, _ := subjectivelogic.NewOpinion(0.4, 0, 0.6, 0.5) 
	opinion3, _ := subjectivelogic.NewOpinion(1, 0, 0, 0.5) 
	
	out1, err1 := subjectivelogic.Addition(&opinion1, &opinion2) //Case 1 
	out2, err2 := subjectivelogic.Addition(&opinion1, &opinion3) //Case 2 
	
	fmt.Println("Case 1:", "Opinion = ", out1.ToString(), "Error:", err1) 
	fmt.Println("Case 2:", "Opinion = ", out2.ToString(), "Error:", err2) 

}
```

The above code snippet shows the usage of the Addition operator. Case 1 uses two Opinions that are not problematic for the Addition operator, resulting in a valid output Opinion and no error:

```
Case 1: Opinion =  0.6000000000000001, 0.1, 0.3, 1 Error: <nil>
```

 Case 2 uses two valid Opinions that are problematic for the Addition operator as the sum of their belief masses exceeds $1$. This results in the Addition operator returning a zeroed Opinion and an error. 

```
Case 2: Opinion =  0, 0, 0, 0 Error: Addition: Check the validity of your input values
```
---

### Complement
This implements the Complement Operator as defined in Subjective Logic:
```math
	\omega_{\overline{x}}  :
	\begin{cases}
		b_{\overline{x}} = d_x \\
		d_{\overline{x}} = b_x \\
		u_{\overline{x}} = u_x \\
		a_{\overline{x}} = 1 - a_x
	\end{cases}
```

#### API Reference

```go
func Complement(opinion *Opinion) (Opinion, error)
```

#### Problematic Inputs
There are no problematic inputs for this operator, as long as they are valid opinions.

#### Example

```go
func main() {

	opinion, _ := subjectivelogic.NewOpinion(0.2, 0.8, 0, 0.5) 
	
	out, err := subjectivelogic.Complement(&opinion) 

	if err != nil { 
		fmt.Println("Error:", err) 
	} else { 
		fmt.Println("Output:", out, err) 
	}
}
```
The code snippet above shows the usage of the Complement operator. This specific example will result in the following output:

```
Output: {0.8 0.2 0 0.5} <nil>
```
---

### Binomial Multiplication
This implements the Binomial Multiplication Operator as defined in Subjective Logic:
$$
	\omega_{x \wedge y}  :
	\begin{cases}
		b_{x \wedge y} = b_x b_y + \frac{(1-a_x) a_y b_x u_y + a_x (1-a_y) b_y u_x}{1 - a_x a_y} \\
		d_{x \wedge y} = d_x + d_y - d_x d_y \\
		u_{x \wedge y} = u_x u_y + \frac{(1-a_y) b_x u_y + (1-a_x) u_x b_y}{1 - a_x a_y} \\
		a_{x \wedge y} = a_x a_y
	\end{cases}       
$$

#### API Reference

```go
func Multiplication(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error)
```

#### Problematic Inputs
The Binomial Multiplication Operator will try to divide through $0$, if $a_x = a_y = 1$, hence this is not allowed and an error is returned.

#### Example

```go
func main() {

	opinion1, _ := subjectivelogic.NewOpinion(0.2, 0.8, 0, 0.5)
	opinion2, _ := subjectivelogic.NewOpinion(0.4, 0, 0.6, 0.5) 

	out, err := subjectivelogic.Multiplication(&opinion1, &opinion2)

	if err != nil { 
		fmt.PrintlnImplementation("Error:", err)
	} else {
		fmt.Println("Output:", out, err) 
	}	
}
``` 

The code snippet above shows the usage of the Multiplication operator. This specific example will result in the following output:

```
Output: {0.12000000000000002 0.8 0.08 0.25} <nil>
```
---

### Binomial Comultiplication
This implements the Comultiplication Operator as defined in Subjective Logic:

$$
	\omega_{x \wedge y}  :
	\begin{cases}
		b_{x \wedge y} = b_x + b_y - b_x b_y \\
		d_{x \wedge y} = d_x d_y + \frac{a_x (1-a_y) d_x u_y + (1-a_x) a_y u_x d_y}{a_x + a_y - a_x a_y} \\
		u_{x \wedge y} = u_x u_y + \frac{a_y d_x u_y + a_x u_x d_y}{a_x + a_y - a_x a_y} \\
		a_{x \wedge y} = a_x + a_y - a_x a_y
	\end{cases}       
$$

#### API Reference

```go
func Comultiplication(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error)
```

#### Problematic Inputs
The Binomial Comultiplication Operator will try to divide through $0$, if $a_x = a_y = 0$, hence this is not allowed and an error will be returned.

#### Example 

```go
func main() {

	opinion1, _ := subjectivelogic.NewOpinion(0.2, 0.8, 0, 0.5)
	opinion2, _ := subjectivelogic.NewOpinion(0.4, 0, 0.6, 0.5)

	out, err := subjectivelogic.Comultiplication(&opinion1, &opinion2)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Output:", out, err)
	}
}
```

The code snippet above shows the usage of the Comultiplication operator. This specific example will result in the following output:

```
Output: {0.52 0.16 0.32 0.75} <nil>
```
---

### Belief Constraint Fusion
==This implements the Belief Constraint Fusion Operator as defined in Subjective Logic:==


NEEDS TO BE EDITED
$$
	\omega_{X}^{(A\&B)}  :
	\begin{cases}
		b_{X}^{(A\&B)} = \frac{b_{X}^{A}u_{X}^{B} + b_{X}^{B}u_{X}^{A} + b_{X}^{A}b_{X}^{B}}{1 - Con} \\
		d_{X}^{(A\&B)} = \frac{d_{X}^{A}u_{X}^{B} + d_{X}^{B}u_{X}^{A} + d_{X}^{A}d_{X}^{B}}{1 - Con} \\
		u_{X}^{(A\&B)} = \frac{u_{X}^{A}u_{X}^{B}}{1-Con} \\
		a_{X}^{(A\&B)} = \frac{a_{X}^{A}(1-u_{X}^{A})+a_{X}^{B}(1-u_{X}^{B})}{2-u_{X}^{A}-u_{X}^{B}} & \text{for} \hspace{2mm} u_{X}^{A} + u_{X}^{B} < 2 \\
		a_{X}^{(A\&B)} = \frac{a_{X}^{A}+a_{X}^{B}}{2} & \text{for} \hspace{2mm}u_{X}^{A} = u_{X}^{B} = 1
	\end{cases}       
$$

$$
	Con = b_{X}^{A}d_{X}^{B} + d_{X}^{A}b_{X}^{B}
$$

#### API Reference

```go
func ConstraintFusion(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error)
```

#### Problematic Inputs
The Belief Constraint Fusion Operator will try to divide through $0$, if the conflict variable $Con = 1$, and an error will be returned.

#### Example

```go
func main() {

	opinion1, _ := subjectivelogic.NewOpinion(0.2, 0.8, 0, 0.5) 
	opinion2, _ := subjectivelogic.NewOpinion(0.4, 0, 0.6, 0.5)

	out, err := subjectivelogic.ConstraintFusion(&opinion1, &opinion2)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Output:", out, err)
	}
}
```

The code snippet above shows the usage of the Belief Constraint operator. This specific example will result in the following output:

```
Output: {0.2941176470588236 0.7058823529411764 0 0.5} <nil>
```
---

### Cumulative Fusion
This implements the Aleatory Cumulative Fusion Operator as defined in Subjective Logic:

Case I: For $u_{X}^{A} \neq 0 \vee u_{X}^{B} \neq 0$:
$$
	\omega_{X}^{(A\diamond B)}  :
	\begin{cases}
		b_{X}^{(A\diamond B)}(x) = \frac{b_{X}^{A}(x)u_{X}^{B} + b_{X}^{B}(x)u_{X}^{A}}{u_{X}^{A} + u_{X}^{B} - u_{X}^{A} u_{X}^{B}} \\
		u_{X}^{(A\diamond B)} = \frac{u_{X}^{A}u_{X}^{B}}{u_{X}^{A} + u_{X}^{B} - u_{X}^{A} u_{X}^{B}} \\
		a_{X}^{(A\diamond B)}(x) = \frac{a_{X}^{A}(x)u_{X}^{B} + a_{X}^{B}(x)u_{X}^{A} - (a_{X}^{A}(x) +  a_{X}^{B}(x))u_{X}^{B}u_{X}^{A}}{u_{X}^{A} + u_{X}^{B} - 2 u_{X}^{A} u_{X}^{B}} & \text{if} \hspace{2mm} u_{X}^{A} \neq 1 \vee u_{X}^{B} \neq 1 \\
		a_{X}^{(A\diamond B)}(x) = \frac{a_{X}^{A}(x)+a_{X}^{B}(x)}{2} & \text{if} \hspace{2mm}u_{X}^{A} = u_{X}^{B} = 1
	\end{cases}       
$$


Case II: For $u_{X}^{A} = u_{X}^{B} = 0$:
$$
	\omega_{X}^{(A\diamond B)}  :
	\begin{cases}
		b_{X}^{(A\diamond B)}(x) = \gamma_{X}^{A} b_{X}^{A}(x) + \gamma_{X}^{B} b_{X}^{B}(x) \\
		u_{X}^{(A\diamond B)} = 0 \\
		a_{X}^{(A\diamond B)}(x) = \gamma_{X}^{A} a_{X}^{A}(x) + \gamma_{X}^{B} a_{X}^{B}(x)
	\end{cases}       
$$
$$
	 \text{where} \hspace{2mm} \gamma_{X}^{A} = \gamma_{X}^{B} = 0.5     
$$

#### API Reference

```go
func CumulativeFusion(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error)
```

#### Problematic Inputs
There are no problematic inputs for this operator, as long as they are valid opinions.

#### Example

```go
func main() {

	opinion1, _ := subjectivelogic.NewOpinion(0.2, 0.8, 0, 0.5) 
	opinion2, _ := subjectivelogic.NewOpinion(0.4, 0, 0.6, 0.5)

	out, err := subjectivelogic.CumulativeFusion(&opinion1, &opinion2)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Output:", out, err)
	}
}
```
The code snippet above shows the usage of the Cumulative fusion operator. This specific example will result in the following output:

```
Output: {0.2 0.8 0 0.5} <nil>
```
---

### Averaging Fusion
This implements the Averaging Fusion Operator as defined in Subjective Logic:

Case I: For $u_{X}^{A} \neq 0 \vee u_{X}^{B} \neq 0$:
$$
	\omega_{X}^{(A\underline{\diamond} B)}  :
	\begin{cases}
		b_{X}^{(A\underline{\diamond} B)}(x) = \frac{b_{X}^{A}(x)u_{X}^{B} + b_{X}^{B}(x)u_{X}^{A}}{u_{X}^{A} + u_{X}^{B}} \\
		u_{X}^{(A\underline{\diamond} B)} = \frac{2 u_{X}^{A}u_{X}^{B}}{u_{X}^{A} + u_{X}^{B}} \\
		a_{X}^{(A\underline{\diamond} B)}(x) = \frac{a_{X}^{A}(x)+a_{X}^{B}(x)}{2}
	\end{cases}       
$$


Case II: For $u_{X}^{A} = u_{X}^{B} = 0$:
$$
	\omega_{X}^{(A\underline{\diamond} B)}  :
	\begin{cases}
		b_{X}^{(A\underline{\diamond} B)}(x) = \gamma_{X}^{A} b_{X}^{A}(x) + \gamma_{X}^{B} b_{X}^{B}(x) \\
		u_{X}^{(A\underline{\diamond} B)} = 0 \\
		a_{X}^{(A\underline{\diamond} B)}(x) = \gamma_{X}^{A} a_{X}^{A}(x) + \gamma_{X}^{B} a_{X}^{B}(x)
	\end{cases}       
$$
$$
	\text{where} \hspace{2mm} \gamma_{X}^{A} = \gamma_{X}^{B} = 0.5     
$$

#### API Reference

```go
func AveragingFusion(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error)
```

#### Problematic Inputs
There are no problematic inputs for this operator, as long as they are valid opinions.

#### Example

```go
func main() {

	opinion1, _ := subjectivelogic.NewOpinion(0.2, 0.8, 0, 0.5) 
	opinion2, _ := subjectivelogic.NewOpinion(0.4, 0, 0.6, 0.5)

	out, err := subjectivelogic.AveragingFusion(&opinion1, &opinion2)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Output:", out, err)
	}
}
```

The code snippet above shows the usage of the Averaging fusion operator. This specific example will result in the following output:

```
Output: {0.2 0.8 0 0.5} <nil>
```

---

### Weighted Fusion
This implements the Weighted Fusion Operator as defined in Subjective Logic:

Case I: For $(u_{X}^{A} \neq 0 \vee u_{X}^{B} \neq 0) \wedge (u_{X}^{A} \neq 1 \vee u_{X}^{B} \neq 1)$:
$$
	\omega_{X}^{(A\diamond B)}  :
	\begin{cases}
		b_{X}^{(A\diamond B)}(x) = \frac{b_{X}^{A}(x)(1-u_{X}^{A})u_{X}^{B} + b_{X}^{B}(x)(1-u_{X}^{B})u_{X}^{A}}{u_{X}^{A} + u_{X}^{B} - 2u_{X}^{A} u_{X}^{B}} \\
		u_{X}^{(A\diamond B)} = \frac{(2 - u_{X}^{A} - u_{X}^{B})u_{X}^{A}u_{X}^{B}}{u_{X}^{A} + u_{X}^{B} - 2u_{X}^{A} u_{X}^{B}} \\
		a_{X}^{(A\diamond B)}(x) = \frac{a_{X}^{A}(x)(1-u_{X}^{A}) + a_{X}^{B}(x)(1-u_{X}^{B})}{2 - u_{X}^{A} - u_{X}^{B}}
	\end{cases}       
$$


Case II: For $u_{X}^{A} = u_{X}^{B} = 0$:
$$
	\omega_{X}^{(A\diamond B)}  :
	\begin{cases}
		b_{X}^{(A\diamond B)}(x) = \gamma_{X}^{A} b_{X}^{A}(x) + \gamma_{X}^{B} b_{X}^{B}(x) \\
		u_{X}^{(A\diamond B)} = 0 \\
		a_{X}^{(A\diamond B)}(x) = \gamma_{X}^{A} a_{X}^{A}(x) + \gamma_{X}^{B} a_{X}^{B}(x)
	\end{cases}       
$$
$$
	\text{where} \hspace{2mm} \gamma_{X}^{A} = \gamma_{X}^{B} = 0.5     
$$


Case III: $u_{X}^{A} = u_{X}^{B} = 1$:
$$
	\omega_{X}^{(A\diamond B)}  :
	\begin{cases}
		b_{X}^{(A\diamond B)}(x) = 0 \\
		u_{X}^{(A\diamond B)} = 1 \\
		a_{X}^{(A\diamond B)}(x) = \frac{a_{X}^{A}(x)+a_{X}^{B}(x)}{2}
	\end{cases}       
$$

#### API Reference
```go
func WeightedFusion(opinion1 *Opinion, opinion2 *Opinion) (Opinion, error)
```

#### Problematic Inputs
There are no problematic inputs for this operator, as long as they are valid opinions.


#### Example
```go
func main() {

	opinion1, _ := subjectivelogic.NewOpinion(0.2, 0.8, 0, 0.5) 
	opinion2, _ := subjectivelogic.NewOpinion(0.4, 0, 0.6, 0.5)

	out, err := subjectivelogic.WeightedFusion(&opinion1, &opinion2)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Output:", out, err)
	}
}
```
The code snippet above shows the usage of the Weighted fusion operator. This specific example will result in the following output:

```
Output: {0.2 0.8 0 0.5} <nil>
```

---

### Trust Discounting
This implements the Trust Discounting Operator as defined in Subjective Logic:

$$
	\omega_{X}^{[A;B]}  :
	\begin{cases}
		b_{X}^{[A;B]}(x) & = P_B^{A}*b_{X}^B(x) \\
		u_{X}^{[A;B]} & = 1 - P_B^{A}*\sum_{x \in X}b_{X}^B(x) \\
		a_{X}^{[A;B]}(x) & = a_{X}^B(x)
	\end{cases}       
$$

#### API Reference

#### Problematic Inputs

#### Example
TBD 

### Trust Discounting for Multi-edge Path
This implements the Trust Discounting Operator for Multi-edge Paths as defined in Subjective Logic.

If $[A_1, ..., A_n]$ denotes the referral trust path and $[A_n, X]$ the functional trust, the Trust Discounting Operator for multi-edge paths is defined as:

$$
	\omega_{X}^{A_1}  :
	\begin{cases}
		b_{X}^{A_1}(x) & = P_{A_n}^{A_1}*b_{X}^{A_n}(x) \\
		u_{X}^{A_1} & = 1 - P_{A_n}^{A_1}*\sum_{x \in X}b_{X}^{A_n}(x) \\
		a_{X}^{A_1}(x) & = a_{X}^{A_n}(x)
	\end{cases}       
$$
and:
$$
	P_{A_n}^{A_1} =  \prod_{i=1}^{n-1} P_{A_{i+1}}^{A_i}     
$$

#### API Reference

#### Problematic Inputs

#### Example
TBD 

## Contributing
Contributions are very welcome! Please let us know if you find an issue and have ideas for improvement. Alternately, open an issue or submit a pull request on GitHub. 

## License
This project is licensed under the Apache License, Verion 2.0 - see the LICENSE file for details.

## Contact
In case of problems or issues, please open an issue on GitHub.   
