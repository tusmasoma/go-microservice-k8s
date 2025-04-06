package testing

import "testing"

func ValidateErr(t *testing.T, expectedErr, actualErr error) {
	if (actualErr != nil) != (expectedErr != nil) {
		t.Errorf("actualError = %v, expectedError %v", actualErr, expectedErr)
	}
	if actualErr != nil && expectedErr != nil && actualErr.Error() != expectedErr.Error() {
		t.Errorf("actualError = %v, expectedError %v", actualErr, expectedErr)
	}
}
