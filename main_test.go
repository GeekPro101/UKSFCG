package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testChangelog Changelog = Changelog{}
)

func TestGetChanges(t *testing.T) {
	fakefile := "# Changes from release 2022/06 to 2022/07\n2. Bug - Corrected Alderney (EGJA) runway coords - thanks to @sdkjsdklfj (John Doe)\n3. AIRAC (2207) - Updated Cranfield (EGTC) SMR - thanks to @sdfsdf (Doe John)\nfakeline"
	changes := GetChanges([]byte(fakefile))
	intendedChanges := []string{
		"Bug - Corrected Alderney (EGJA) runway coords - thanks to @sdkjsdklfj (John Doe)",
		"AIRAC (2207) - Updated Cranfield (EGTC) SMR - thanks to @sdfsdf (Doe John)",
	}
	assert.ElementsMatch(t, intendedChanges, changes)
}

func TestChangesSorter(t *testing.T) {
	changes := []string{
		"Bug - Corrected Alderney (EGJA) runway coords - thanks to @sdkjsdklfj (John Doe)",
		"AIRAC (2207) - Updated Cranfield (EGTC) SMR - thanks to @sdfsdf (Doe John)",
		"AIRAC (2207) - Updated Inverness (EGPE) RWYs 05/23 and 11/29 coords - this was very silly - thanks to @sdfsdf (Smith)",
		"Enhancement - Added missing heli points and holds to Gloucestershire (EGBJ) SMR",
	}
	testChangelog = Changelog{
		Changes: changes,
	}
	airacList, otherList := testChangelog.ChangesSorter()
	intendedAiracList := []string{
		"AIRAC (2207) - Updated Cranfield (EGTC) SMR",
		"AIRAC (2207) - Updated Inverness (EGPE) RWYs 05/23 and 11/29 coords - this was very silly",
	}
	intendedOtherList := []string{
		"Bug - Corrected Alderney (EGJA) runway coords",
		"Enhancement - Added missing heli points and holds to Gloucestershire (EGBJ) SMR",
	}
	assert.ElementsMatch(t, intendedAiracList, airacList)
	assert.ElementsMatch(t, intendedOtherList, otherList)
}

func TestAiracMapGen(t *testing.T) {
	airacChanges := []string{
		"AIRAC (2207) - Updated Cranfield (EGTC) SMR",
		"AIRAC (2206) - Updated Inverness (EGPE) RWYs 05/23 and 11/29 coords - this was very silly",
		"AIRAC (2205)",
	}
	testChangelog = Changelog{
		AIRACList: airacChanges,
	}
	airacMap, airacs := testChangelog.AIRACMapGen()
	intendedAiracMap := map[string][]string{
		"2206": {"Updated Inverness (EGPE) RWYs 05/23 and 11/29 coords - this was very silly"},
		"2207": {"Updated Cranfield (EGTC) SMR"},
	}
	intendedAiracs := []int{2207, 2206}
	assert.Equal(t, intendedAiracMap, airacMap)
	assert.Equal(t, intendedAiracs, airacs)
}

func TestContribGen(t *testing.T) {
	changes := []string{
		"Bug - Corrected Alderney (EGJA) runway coords - thanks to @sdkjsdklfj (John Doe)",
		"AIRAC (2207) - Updated Cranfield (EGTC) SMR - thanks to @sdfsdf (Doe John)",
		"AIRAC (2207) - Updated Inverness (EGPE) RWYs 05/23 and 11/29 coords - this was very silly - thanks to @sdfsdf (Smith)",
		"Enhancement - Added missing heli points and holds to Gloucestershire (EGBJ) SMR",
	}
	testChangelog = Changelog{
		Changes: changes,
	}
	contribs := testChangelog.ContribGen()
	intendedContribs := []string{
		"Smith",
		"John Doe",
		"Doe John",
	}
	assert.ElementsMatch(t, intendedContribs, contribs)
}
