package main

import (
	"errors"
)

type block struct {
	ecCount    int
	dataCount  int
	dataWords  []uint8
	errorWords []uint8
}

// Retunrs the combined data words and error words
func generateErrorWords(codewords []uint8, version int, ecLevel string) []uint8 {
	var blocks []*block

	// Split codewords into blocks
	var c int
	for _, blockType := range codeWordTable[version][ecLevel].blocks {
		for i := 0; i < blockType.count; i++ {
			blocks = append(blocks, &block{
				dataWords: codewords[c : c+blockType.dataWords],
				ecCount:   codeWordTable[version][ecLevel].ecWordsPerBlock,
				dataCount: blockType.dataWords,
			})
			c += blockType.dataWords
		}
	}

	// Generate error codes for each block
	for _, b := range blocks {
		b.errorWords = rsEncode(b.dataWords, b.ecCount)
	}

	// Assemble the final sequence, taking each block in turn
	var out []uint8
	for i := 0; i < len(blocks[len(blocks)-1].dataWords); i++ {
		for _, v := range blocks {
			if i < len(v.dataWords) {
				out = append(out, v.dataWords[i])
			}
		}
	}
	for i := 0; i < len(blocks[len(blocks)-1].errorWords); i++ {
		for _, v := range blocks {
			if i < len(v.errorWords) {
				out = append(out, v.errorWords[i])
			}
		}
	}

	return out
}

func correctDataWords(allWords []uint8, version int, ecLevel string) ([]uint8, error) {
	// Generate the empty blocks
	var blocks []*block
	for _, blockType := range codeWordTable[version][ecLevel].blocks {
		for i := 0; i < blockType.count; i++ {
			blocks = append(blocks, &block{
				ecCount:   codeWordTable[version][ecLevel].ecWordsPerBlock,
				dataCount: blockType.dataWords,
			})
		}
	}

	// Disassmble the sequence, putting to each block in turn
	var i int
	for {
		for k, v := range blocks {
			// Fill up data blocks, then ec blocks
			if len(v.dataWords) < v.dataCount {
				v.dataWords = append(v.dataWords, allWords[i])
				i++
			} else if len(v.errorWords) < v.ecCount {
				v.errorWords = append(v.errorWords, allWords[i])
				i++
			} else if k == len(blocks)-1 {
				// Last block is filled, we are done
				goto ec
			}
		}
	}

ec:
	// Perform error correction
	var out []uint8
	for _, v := range blocks {
		var msgInt []int
		for _, v := range v.dataWords {
			msgInt = append(msgInt, int(v))
		}
		for _, v := range v.errorWords {
			msgInt = append(msgInt, int(v))
		}

		corrected, _, err := correctMessage(msgInt, v.ecCount)
		if err != nil {
			return []uint8{}, err
		}

		correctedByt := make([]uint8, len(corrected))
		for k := range corrected {
			correctedByt[k] = uint8(corrected[k])
		}

		out = append(out, correctedByt...)
	}

	return out, nil
}

func rsGenerator(symbols int) []int {
	g := []int{1}
	for i := 0; i < symbols; i++ {
		g = gfPolyMul(g, []int{1, gfPow(2, i)})
	}
	return g
}

// Returns the ec codewords
func rsEncode(msg []uint8, symbols int) []uint8 {
	gen := rsGenerator(symbols)

	m := make([]int, len(msg))
	for k, v := range msg {
		m[k] = int(v)
	}

	_, remainder := gfPolyDiv(append(m, make([]int, len(gen)-1)...), gen)

	r := make([]uint8, len(remainder))
	for k, v := range remainder {
		r[k] = uint8(v)
	}
	return r
}
func calcSyndromes(msg []int, nsym int) []int {
	// Given the received codeword msg and the number of error correcting symbols (nsym), computes the syndromes polynomial.
	// Mathematically, it's essentially equivalent to a Fourrier Transform (Chien search being the inverse).

	// Note the "[0] +" : we add a 0 coefficient for the lowest degree (the constant). This effectively shifts the syndrome, and will shift every computations depending on the syndromes (such as the errors locator polynomial, errors evaluator polynomial, etc. but not the errors positions).
	// This is not necessary, you can adapt subsequent computations to start from 0 instead of skipping the first iteration (ie, the often seen range(1, n-k+1)),
	synd := make([]int, nsym)
	for i := 0; i < nsym; i++ {
		synd[i] = gfPolyEval(msg, gfPow(2, i))
	}

	return append([]int{0}, synd...) // pad with one 0 for mathematical precision (else we can end up with weird calculations sometimes)
}

// Returns true if message is ok, false if errors/erasures are present
func checkMessage(msg []int, nsym int) bool {
	for _, v := range calcSyndromes(msg, nsym) {
		if v != 0 {
			return false
		}
	}

	return true
}

func findErrataLocator(pos []int) []int {
	// Compute the erasures/errors/errata locator polynomial from the erasures/errors/errata positions
	// (the positions must be relative to the x coefficient, eg: "hello worldxxxxxxxxx" is tampered to "h_ll_ worldxxxxxxxxx"
	// with xxxxxxxxx being the ecc of length n-k=9, here the string positions are [1, 4], but the coefficients are reversed
	// since the ecc characters are placed as the first coefficients of the polynomial, thus the coefficients of the
	// erased characters are n-1 - [1, 4] = [18, 15] = erasures_loc to be specified as an argument.'''

	loc := []int{1} // just to init because we will multiply, so it must be 1 so that the multiplication starts correctly without nulling any term
	// erasures_loc = product(1 - x*alpha**i) for i in erasures_pos and where alpha is the alpha chosen to evaluate polynomials.
	for _, v := range pos {
		loc = gfPolyMul(loc, gfPolyAdd([]int{1}, append([]int{gfPow(2, v)}, 0)))
	}

	return loc
}

func findErrorEvaluator(synd, err_loc []int, nsym int) []int {
	// Compute the error (or erasures if you supply sigma=erasures locator polynomial, or errata) evaluator polynomial Omega
	// from the syndrome and the error/erasures/errata locator Sigma.

	// Omega(x) = [ Synd(x) * Error_loc(x) ] mod x^(n-k+1)
	_, remainder := gfPolyDiv(gfPolyMul(synd, err_loc), append([]int{1}, make([]int, nsym+1)...)) // first multiply syndromes * errata_locator, then do a
	// polynomial division to truncate the polynomial to the
	// required length

	// Faster way that is equivalent
	//remainder = gf_poly_mul(synd, err_loc) // first multiply the syndromes with the errata locator polynomial
	//remainder = remainder[len(remainder)-(nsym+1):] // then slice the list to truncate it (which represents the polynomial), which
	// is equivalent to dividing by a polynomial of the length we want

	return remainder
}

func copyAndReverse(in []int) []int {
	out := make([]int, len(in))
	for k, v := range in {
		out[len(in)-1-k] = v
	}
	return out
}

// err_pos is a list of the positions of the errors/erasures/errata
func correctErrata(msgIn, synd, err_pos []int) ([]int, error) {
	// Forney algorithm, computes the values (error magnitude) to correct the input message.
	// calculate errata locator polynomial to correct both errors and erasures (by combining the errors positions given by the error locator polynomial found by BM with the erasures positions given by caller)

	var coef_pos []int
	for _, p := range err_pos {
		coef_pos = append(coef_pos, len(msgIn)-1-p) // need to convert the positions to coefficients degrees for the errata locator algo to work (eg: instead of [0, 1, 2] it will become [len(msg)-1, len(msg)-2, len(msg) -3])
	}

	err_loc := findErrataLocator(coef_pos)
	// calculate errata evaluator polynomial (often called Omega or Gamma in academic papers)
	err_eval := copyAndReverse(findErrorEvaluator(copyAndReverse(synd), err_loc, len(err_loc)-1))

	// Second part of Chien search to get the error location polynomial X from the error positions in err_pos (the roots of the error locator polynomial, ie, where it evaluates to 0)
	var X []int // will store the position of the errors
	for i := 0; i < len(coef_pos); i++ {
		l := 255 - coef_pos[i]
		X = append(X, gfPow(2, -l))
	}

	// Forney algorithm: compute the magnitudes
	E := make([]int, len(msgIn)) // will store the values that need to be corrected (substracted) to the message containing errors. This is sometimes called the error magnitude polynomial.
	Xlength := len(X)
	for i, Xi := range X {

		XiInv := gfInverse(Xi)

		// Compute the formal derivative of the error locator polynomial (see Blahut, Algebraic codes for data transmission, pp 196-197).
		// the formal derivative of the errata locator is used as the denominator of the Forney Algorithm, which simply says that the ith error value is given by error_evaluator(gf_inverse(Xi)) / error_locator_derivative(gf_inverse(Xi)). See Blahut, Algebraic codes for data transmission, pp 196-197.
		var err_loc_prime_tmp []int
		for j := 0; j < Xlength; j++ {
			if j != i {
				err_loc_prime_tmp = append(err_loc_prime_tmp, 1^gfQuickMul(XiInv, X[j]))
			}
		}

		// compute the product, which is the denominator of the Forney algorithm (errata locator derivative)
		err_loc_prime := 1
		for _, coef := range err_loc_prime_tmp {
			err_loc_prime = gfQuickMul(err_loc_prime, coef)
		}
		// equivalent to: err_loc_prime = functools.reduce(gf_mul, err_loc_prime_tmp, 1)

		// Compute y (evaluation of the errata evaluator polynomial)
		// This is a more faithful translation of the theoretical equation contrary to the old forney method. Here it is an exact reproduction:
		// Yl = omega(Xl.inverse()) / prod(1 - Xj*Xl.inverse()) for j in len(X)
		y := gfPolyEval(copyAndReverse(err_eval), XiInv) // numerator of the Forney algorithm (errata evaluator evaluated)
		y = gfQuickMul(gfPow(Xi, 1), y)

		// Check: err_loc_prime (the divisor) should not be zero.
		if err_loc_prime == 0 {
			return []int{}, errors.New("could not find error magnitude") // Could not find error magnitude
		}

		// Compute the magnitude
		magnitude := gfQuickDiv(y, err_loc_prime) // magnitude value of the error, calculated by the Forney algorithm (an equation in fact): dividing the errata evaluator with the errata locator derivative gives us the errata magnitude (ie, value to repair) the ith symbol
		E[err_pos[i]] = magnitude                 // store the magnitude for this error into the magnitude polynomial
	}

	// Apply the correction of values to get our message corrected! (note that the ecc bytes also gets corrected!)
	// (this isn't the Forney algorithm, we just apply the result of decoding here)
	msgIn = gfPolyAdd(msgIn, E) // equivalent to Ci = Ri - Ei where Ci is the correct message, Ri the received (senseword) message, and Ei the errata magnitudes (minus is replaced by XOR since it's equivalent in GF(2^p)). So in fact here we substract from the received message the errors magnitude, which logically corrects the value to what it should be.
	return msgIn, nil
}

func findErrorLocator(synd []int, nsym int, erase_loc []int, erase_count int) ([]int, error) {
	// Find error/errata locator and evaluator polynomials with Berlekamp-Massey algorithm'''
	// The idea is that BM will iteratively estimate the error locator polynomial.
	// To do this, it will compute a Discrepancy term called Delta, which will tell us if the error locator polynomial needs an update or not
	// (hence why it's called discrepancy: it tells us when we are getting off board from the correct value).

	// Init the polynomials
	var err_loc []int
	var old_loc []int
	if len(erase_loc) > 0 { // if the erasure locator polynomial is supplied, we init with its value, so that we include erasures in the final locator polynomial
		err_loc = make([]int, len(erase_loc))
		copy(err_loc, erase_loc)
		old_loc = make([]int, len(erase_loc))
		copy(old_loc, erase_loc)
	} else {
		err_loc = []int{1} // This is the main variable we want to fill, also called Sigma in other notations or more formally the errors/errata locator polynomial.
		old_loc = []int{1} // BM is an iterative algorithm, and we need the errata locator polynomial of the previous iteration in order to update other necessary variables.
	}
	//L = 0 // update flag variable, not needed here because we use an alternative equivalent way of checking if update is needed (but using the flag could potentially be faster depending on if using length(list) is taking linear time in your language, here in Python it's constant so it's as fast.

	// Fix the syndrome shifting: when computing the syndrome, some implementations may prepend a 0 coefficient for the lowest degree term (the constant). This is a case of syndrome shifting, thus the syndrome will be bigger than the number of ecc symbols (I don't know what purpose serves this shifting). If that's the case, then we need to account for the syndrome shifting when we use the syndrome such as inside BM, by skipping those prepended coefficients.
	// Another way to detect the shifting is to detect the 0 coefficients: by definition, a syndrome does not contain any 0 coefficient (except if there are no errors/erasures, in this case they are all 0). This however doesn't work with the modified Forney syndrome, which set to 0 the coefficients corresponding to erasures, leaving only the coefficients corresponding to errors.
	synd_shift := len(synd) - nsym

	for i := 0; i < nsym-erase_count; i++ { // generally: nsym-erase_count == len(synd), except when you input a partial erase_loc and using the full syndrome instead of the Forney syndrome, in which case nsym-erase_count is more correct (len(synd) will fail badly with IndexError).
		var K int
		if len(erase_loc) > 0 { // if an erasures locator polynomial was provided to init the errors locator polynomial, then we must skip the FIRST erase_count iterations (not the last iterations, this is very important!)
			K = erase_count + i + synd_shift
		} else { // if erasures locator is not provided, then either there's no erasures to account or we use the Forney syndromes, so we don't need to use erase_count nor erase_loc (the erasures have been trimmed out of the Forney syndromes).
			K = i + synd_shift
		}

		// Compute the discrepancy Delta
		// Here is the close-to-the-books operation to compute the discrepancy Delta: it's a simple polynomial multiplication of error locator with the syndromes, and then we get the Kth element.
		//delta = gf_poly_mul(err_loc[::-1], synd)[K] // theoretically it should be gf_poly_add(synd[::-1], [1])[::-1] instead of just synd, but it seems it's not absolutely necessary to correctly decode.
		// But this can be optimized: since we only need the Kth element, we don't need to compute the polynomial multiplication for any other element but the Kth. Thus to optimize, we compute the polymul only at the item we need, skipping the rest (avoiding a nested loop, thus we are linear time instead of quadratic).
		// This optimization is actually described in several figures of the book "Algebraic codes for data transmission", Blahut, Richard E., 2003, Cambridge university press.
		delta := synd[K]
		for j := 1; j < len(err_loc); j++ {
			delta ^= gfQuickMul(err_loc[len(err_loc)-(j+1)], synd[K-j]) // delta is also called discrepancy. Here we do a partial polynomial multiplication (ie, we compute the polynomial multiplication only for the term of degree K). Should be equivalent to brownanrs.polynomial.mul_at().
		}
		//print "delta", K, delta, list(gf_poly_mul(err_loc[::-1], synd)) // debugline

		// Shift polynomials to compute the next degree
		old_loc = append(old_loc, 0)

		// Iteratively estimate the errata locator and evaluator polynomials
		if delta != 0 { // Update only if there's a discrepancy
			if len(old_loc) > len(err_loc) { // Rule B (rule A is implicitly defined because rule A just says that we skip any modification for this iteration)
				//if 2*L <= K+erase_count: // equivalent to len(old_loc) > len(err_loc), as long as L is correctly computed
				// Computing errata locator polynomial Sigma
				new_loc := gfPolyScale(old_loc, delta)
				old_loc = gfPolyScale(err_loc, gfInverse(delta)) // effectively we are doing err_loc * 1/delta = err_loc // delta
				err_loc = new_loc
				// Update the update flag
				//L = K - L // the update flag L is tricky: in Blahut's schema, it's mandatory to use `L = K - L - erase_count` (and indeed in a previous draft of this function, if you forgot to do `- erase_count` it would lead to correcting only 2*(errors+erasures) <= (n-k) instead of 2*errors+erasures <= (n-k)), but in this latest draft, this will lead to a wrong decoding in some cases where it should correctly decode! Thus you should try with and without `- erase_count` to update L on your own implementation and see which one works OK without producing wrong decoding failures.
			}
			// Update with the discrepancy
			err_loc = gfPolyAdd(err_loc, gfPolyScale(old_loc, delta))
		}
	}

	// Check if the result is correct, that there's not too many errors to correct
	for len(err_loc) > 0 && err_loc[0] == 0 {
		err_loc = err_loc[1:] // drop leading 0s, else errs will not be of the correct size
	}
	errs := len(err_loc) - 1
	if (errs-erase_count)*2+erase_count > nsym {
		return []int{}, errors.New("too many errors to correct")
	}

	return err_loc, nil
}

// nmess is len(msg_in)
func findErrors(err_loc []int, nmess int) ([]int, error) {
	// Find the roots (ie, where evaluation = zero) of error polynomial by brute-force trial, this is a sort of Chien's search
	// (but less efficient, Chien's search is a way to evaluate the polynomial such that each evaluation only takes constant time).
	errs := len(err_loc) - 1
	var err_pos []int
	for i := 0; i < nmess; i++ { // normally we should try all 2^8 possible values, but here we optimize to just check the interesting symbols
		if gfPolyEval(err_loc, gfPow(2, i)) == 0 { // It's a 0? Bingo, it's a root of the error locator polynomial,
			err_pos = append(err_pos, nmess-1-i) // in other terms this is the location of an error
		}
	}

	// Sanity check: the number of errors/errata positions found should be exactly the same as the length of the errata locator polynomial
	if len(err_pos) != errs {
		return []int{}, errors.New("too many (or few) errors found by Chien Search for the errata locator polynomial")
	}

	return err_pos, nil
}

func forneySyndromes(synd, pos []int, nmess int) []int {
	// Compute Forney syndromes, which computes a modified syndromes to compute only errors (erasures are trimmed out). Do not confuse this with Forney algorithm, which allows to correct the message based on the location of errors.
	var erase_pos_reversed []int
	for p := range pos {
		erase_pos_reversed = append(erase_pos_reversed, nmess-1-p) // prepare the coefficient degree positions (instead of the erasures positions)
	}

	// Optimized method, all operations are inlined
	fsynd := make([]int, len(synd)-1) // make a copy and trim the first coefficient which is always 0 by definition
	copy(fsynd, synd[1:])
	for i := 0; i < len(pos); i++ {
		x := gfPow(2, erase_pos_reversed[i])
		for j := 0; j < len(fsynd)-1; j++ {
			fsynd[j] = gfQuickMul(fsynd[j], x) ^ fsynd[j+1]
		}
	}

	// Equivalent, theoretical way of computing the modified Forney syndromes: fsynd = (erase_loc * synd) % x^(n-k)
	// See Shao, H. M., Truong, T. K., Deutsch, L. J., & Reed, I. S. (1986, April). A single chip VLSI Reed-Solomon decoder. In Acoustics, Speech, and Signal Processing, IEEE International Conference on ICASSP'86. (Vol. 11, pp. 2151-2154). IEEE.ISO 690
	//erase_loc = rs_find_errata_locator(erase_pos_reversed, generator=generator) // computing the erasures locator polynomial
	//fsynd = gf_poly_mul(erase_loc[::-1], synd[1:]) // then multiply with the syndrome to get the untrimmed forney syndrome
	//fsynd = fsynd[len(pos):] # then trim the first erase_pos coefficients which are useless. Seems to be not necessary, but this reduces the computation time later in BM (thus it's an optimization).

	return fsynd
}

// Reed-Solomon main decoding function
func correctMessage(msg_in []int, nsym int) ([]int, []int, error) {
	if len(msg_in) > 255 {
		return []int{}, []int{}, errors.New("message is too long")
	}

	msg_out := make([]int, len(msg_in)) // copy of message
	copy(msg_out, msg_in)
	// erasures: set them to null bytes for easier decoding (but this is not necessary, they will be corrected anyway, but debugging will be easier with null bytes because the error locator polynomial values will only depend on the errors locations, not their values)
	var erase_pos []int
	// check if there are too many erasures to correct (beyond the Singleton bound)
	if len(erase_pos) > nsym {
		return []int{}, []int{}, errors.New("too many erasures to correct")
	}
	// prepare the syndrome polynomial using only errors (ie: errors = characters that were either replaced by null byte
	// or changed to another character, but we don't know their positions)
	synd := calcSyndromes(msg_out, nsym)
	// check if there's any error/erasure in the input codeword. If not (all syndromes coefficients are 0), then just return the message as-is.
	if checkMessage(msg_in, nsym) {
		// no errors
		return msg_out[:len(msg_out)-nsym], msg_out[len(msg_out)-nsym:], nil
	}

	// compute the Forney syndromes, which hide the erasures from the original syndrome (so that BM will just have to deal with errors, not erasures)
	fsynd := forneySyndromes(synd, erase_pos, len(msg_out))
	// compute the error locator polynomial using Berlekamp-Massey
	err_loc, err := findErrorLocator(fsynd, nsym, []int{}, len(erase_pos))
	if err != nil {
		return []int{}, []int{}, err
	}

	// locate the message errors using Chien search (or brute-force search)
	err_pos, err := findErrors(copyAndReverse(err_loc), len(msg_out))
	if err != nil {
		return []int{}, []int{}, err
	}

	if len(err_pos) == 0 {
		return []int{}, []int{}, errors.New("could not locate error")
	}

	// Find errors values and apply them to correct the message
	// compute errata evaluator and errata magnitude polynomials, then correct errors and erasures
	msg_out, err = correctErrata(msg_out, synd, append(erase_pos, err_pos...)) // note that we here use the original syndrome, not the forney syndrome
	// (because we will correct both errors and erasures, so we need the full syndrome)
	if err != nil {
		return []int{}, []int{}, err
	}

	// check if the final message is fully repaired
	if !checkMessage(msg_out, nsym) {
		return []int{}, []int{}, errors.New("could not correct message")
	}

	// return the successfully decoded message
	return msg_out[:len(msg_out)-nsym], msg_out[len(msg_out)-nsym:], nil // also return the corrected ecc block so that the user can check()
}

// Positive modulo, returns non negative solution to x `%` d
func pmod(x, b int) int {
	x = x % b
	if x >= 0 {
		return x
	}
	if b < 0 {
		return x - b
	}
	return x + b
}

var gfExp = make([]int, 512)
var gfLog = make([]int, 256)

// Uses Russian Peasant Multiplication. Buggered if I know.
func gfMul(x int, y int, prim int, fieldCharacFull int, carryless bool) int {
	r := 0
	for y > 0 {
		if y&1 > 0 {
			if carryless {
				r = r ^ x
			} else {
				r += x
			}
		}

		y = y >> 1
		x = x << 1
		if prim > 0 && x&fieldCharacFull > 0 {
			x = x ^ prim
		}
	}

	return r
}

func gfQuickMul(x, y int) int {
	if x == 0 || y == 0 {
		return 0
	}

	return gfExp[gfLog[x]+gfLog[y]]
}

func gfQuickDiv(x, y int) int {
	if y == 0 {
		panic("cannot divide by 0")
	}

	if x == 0 {
		return 0
	}

	return gfExp[pmod(gfLog[x]+255-gfLog[y], 255)]
}

func gfPow(x, power int) int {
	t := pmod(gfLog[x]*power, 255)
	if t < 0 {
		return gfExp[len(gfExp)+t]
	}
	return gfExp[t]
}

func gfInverse(x int) int {
	return gfExp[255-gfLog[x]]
}

func gfPolyScale(p []int, x int) []int {
	r := make([]int, len(p))
	for i := 0; i < len(p); i++ {
		r[i] = gfQuickMul(p[i], x)
	}
	return r
}

func gfPolyAdd(p, q []int) []int {
	var r []int
	if len(p) > len(q) {
		r = make([]int, len(p))
	} else {
		r = make([]int, len(q))
	}

	for i := 0; i < len(p); i++ {
		r[i+len(r)-len(p)] = p[i]
	}

	for i := 0; i < len(q); i++ {
		r[i+len(r)-len(q)] ^= q[i]
	}

	return r
}

func gfPolyMul(p, q []int) []int {
	r := make([]int, len(p)+len(q)-1)

	for j := 0; j < len(q); j++ {
		for i := 0; i < len(p); i++ {
			r[i+j] ^= gfQuickMul(p[i], q[j])
		}
	}

	return r
}

func gfPolyDiv(dividend, divisor []int) ([]int, []int) {
	msgOut := make([]int, len(dividend))
	copy(msgOut, dividend)

	for i := 0; i < len(dividend)-len(divisor)+1; i++ {
		coef := msgOut[i]
		if coef != 0 {
			for j := 1; j < len(divisor); j++ {
				if divisor[j] != 0 {
					msgOut[i+j] ^= gfQuickMul(divisor[j], coef)
				}
			}
		}
	}

	separator := len(divisor) - 1
	return msgOut[:len(msgOut)-separator], msgOut[len(msgOut)-separator:]
}

func gfPolyEval(poly []int, x int) int {
	y := poly[0]
	for i := 1; i < len(poly); i++ {
		y = gfQuickMul(y, x) ^ poly[i]
	}

	return y
}

var _ = initGFTables()

func initGFTables() bool {
	prim := 0x11d

	x := 1
	for i := 0; i < 255; i++ {
		gfExp[i] = x
		gfLog[x] = i
		x = gfMul(x, 2, prim, 256, true)
	}

	for i := 255; i < 512; i++ {
		gfExp[i] = gfExp[i-255]
	}

	return true
}

func hammingWeight(x int) int {
	var weight int
	for x > 0 {
		weight += x & 1
		x >>= 1
	}
	return weight
}

func checkFormat(fmt int) int {
	g := 0b10100110111
	for i := 4; i >= 0; i-- {
		if fmt&(1<<(i+10)) != 0 {
			fmt ^= g << i
		}
	}

	return fmt
}

func checkVersion(fmt int) int {
	g := 0b1111100100101
	for i := 5; i >= 0; i-- {
		if fmt&(1<<(i+12)) != 0 {
			fmt ^= g << i
		}
	}

	return fmt
}

func decodeFormat(format int) int {
	bestFmt := -1
	bestDist := 15
	for testFmt := 0; testFmt < 32; testFmt++ {
		testCode := (testFmt << 10) ^ checkFormat(testFmt<<10)
		testDist := hammingWeight(format ^ testCode)
		if testDist < bestDist {
			bestDist = testDist
			bestFmt = testFmt
		} else if testDist == bestDist {
			bestFmt = -1
		}
	}

	return bestFmt
}

func decodeVersion(version int) int {
	bestFmt := -1
	bestDist := 18
	for testVer := 0; testVer < 32*2; testVer++ {
		testCode := (testVer << 12) ^ checkVersion(testVer<<12)
		testDist := hammingWeight(version ^ testCode)
		if testDist < bestDist {
			bestDist = testDist
			bestFmt = testVer
		} else if testDist == bestDist {
			bestFmt = -1
		}
	}

	return bestFmt
}

// I love plagiarism
