package debugstream

const (
	unmatched          = iota
	potentialDuplicate = iota
	transposed         = iota
	matched            = iota
)

const (
	prefixLen = 4
)

// jaroDecreased is a lightly modified version of Jaro
// While it takes inspiration from JaroWinkler this seemed more fun.
// Since this is intended for commands the lengths of the strings will be short.
// Presuppose that users will miss the end of a command  or misappend extra data.
// Modified approach for domain that diverges from Jaro-Winkler's prefix strategy.
func jaroDecreased(candidate, registered string) (float64, float64) {
	if len(candidate) == 0 {
		return 0, 0 //  probably shouldnt let it get to this step due to upstream constraints but adding for completeness
	}

	totalMatches := 0.0
	transposed := 0.0

	// denoted as match if within
	matchingDist := len(candidate)
	if matchingDist > len(registered) {
		matchingDist = len(registered)
	}
	matchingDist = (matchingDist - 3) / 2
	if matchingDist < 0 {
		matchingDist = 1
	}

	candidateCharStates := make([]int, len(candidate))
	registeredCharStates := make([]int, len(registered))
	// check for a potential match of every character in the candidate
	for i := range candidate {

		start := i - matchingDist
		if start < 0 {
			start = 0
		}
		end := i + matchingDist
		if end >= len(registered) {
			end = len(registered) - 1
		}

		// for our purposes lets be less mean to extra duplicates

		for j := start; j <= end; j++ {

			if registeredCharStates[j] == matched {
				continue
			}
			if candidate[i] == registered[j] {
				if registeredCharStates[j] > unmatched {
					candidateCharStates[i] = potentialDuplicate
					transposed++
					continue
				}

				if candidateCharStates[i] == potentialDuplicate {
					transposed--
				}
				candidateCharStates[i] = matched
				registeredCharStates[j] = matched
				totalMatches++
				break
			}
		}
	}
	if totalMatches == 0 {
		return 0, 0
	}

	pseudoJaroVal := (totalMatches/float64(len(candidate)) + totalMatches/float64(len(registered)) +
		(totalMatches-transposed)/totalMatches) / 3.0

	prefixChecking := prefixLen
	if len(candidate) < prefixLen {
		prefixChecking = len(candidate)
	}
	boost := 0
	for i := 0; i < prefixChecking; i++ {

		if candidateCharStates[i] > potentialDuplicate {
			boost++
		}
	}
	prefixBoostJaro := pseudoJaroVal

	// dont boost if its super low otherwise ful vs help will have a high boost.
	if boost >= prefixChecking/2.0 {
		boostFac := 1.0 + float64(boost)/10.0

		prefixBoostJaro = (pseudoJaroVal * boostFac)

	}

	lengthBoost := 1.0 + float64(len(registered))/100.0
	prefixAndLengthBoosted := prefixBoostJaro * lengthBoost
	if prefixAndLengthBoosted > 1 {
		prefixAndLengthBoosted = 1
	}

	return pseudoJaroVal, prefixAndLengthBoosted
}
