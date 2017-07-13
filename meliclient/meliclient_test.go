package meliclient

import (
	"testing"
)

func TestWhenCountryIsValidCountryCodeGetterReturnsCountryCode(t *testing.T) {

	client := meliclient.New()
	countryCode, err:= client.getCountryCode("MLA1234")

	if err != nil || countryCode != "MLA" {
		t.Fail()
	}

}

func TestWhenCountryIsInvalidCountryCodeGetterReturnsError(t *testing.T) {

	client := meliclient.New()
	countryCode, err:= client.getCountryCode("ML")

	if err == nil {
		t.Fail()
	}

	countryCode, err = client.getCountryCode("MNN1234")

	if err == nil {
		t.Fail()
	}

}