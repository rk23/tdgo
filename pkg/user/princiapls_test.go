package user

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func assert(t *testing.T, a string, b string) {
	if a != b {
		t.Errorf("%s != %s", a, b)
	}
}

func TestMarshal(t *testing.T) {
	// Open our jsonFile
	golden, err := os.Open("testdata/principalsGolden.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		panic(err)
	}
	defer golden.Close()

	goldenBytes, err := ioutil.ReadAll(golden)
	if err != nil {
		panic(err)
	}

	var myPrincipals Principals
	json.Unmarshal(goldenBytes, &myPrincipals)
	maybeGoldenBytes, _ := json.Marshal(myPrincipals)

	var myPrincipalsRemarshalled Principals
	json.Unmarshal(maybeGoldenBytes, &myPrincipalsRemarshalled)

	// TODO: This doesn't work. Fix it.
	if !reflect.DeepEqual(myPrincipals, myPrincipalsRemarshalled) {
		t.Errorf("Structs not outputing correctly")
	}
}
