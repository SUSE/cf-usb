package gocfbroker

import (
	"fmt"
	"strings"
	"testing"
)

func TestValidate_ServiceID(t *testing.T) {
	t.Parallel()

	valid, msg := validateServiceID(testConfig.Catalog, testConfig.Services[0].ID)()
	if len(msg) == 0 {
		t.Error("Message must never be empty")
	} else if !valid {
		t.Error("It should be valid.")
	}

	valid, msg = validateServiceID(testConfig.Catalog, "junk")()
	if msg != "service_id does not reference a service from the catalog" {
		t.Error("Wrong msg:", msg)
	} else if valid {
		t.Error("It should not be valid.")
	}
}

func TestValidate_PlanID(t *testing.T) {
	t.Parallel()

	valid, msg := validatePlanID(testConfig.Catalog, testConfig.Services[0].ID, testConfig.Services[0].Plans[0].ID)()
	if len(msg) == 0 {
		t.Error("Message must never be empty")
	} else if !valid {
		t.Error("It should be valid.")
	}

	valid, msg = validatePlanID(testConfig.Catalog, testConfig.Services[0].ID, testConfig.Services[1].Plans[0].ID)()
	if msg != "plan_id does not reference a services plan from the catalog" {
		t.Error("Wrong message:", msg)
	} else if valid {
		t.Error("It should not be valid.")
	}
}

func TestValidate_ServiceUpdatable(t *testing.T) {
	t.Parallel()

	valid, msg := validateServiceUpdatable(testConfig.Catalog, testConfig.Services[0].ID)()
	if len(msg) == 0 {
		t.Error("Message must never be empty")
	} else if !valid {
		t.Error("It should be valid.")
	}

	valid, msg = validateServiceUpdatable(testConfig.Catalog, testConfig.Services[1].ID)()
	if msg != "service plan is not updatable" {
		t.Error("Wrong msg:", msg)
	} else if valid {
		t.Error("It should not be valid.")
	}
}

func TestValidate_StrNotEmpty(t *testing.T) {
	t.Parallel()

	if valid, msg := validateStrNotEmpty("a", "name")(); len(msg) == 0 {
		t.Error("Message must never be empty")
	} else if !valid {
		t.Error("It should be valid.")
	}

	if valid, msg := validateStrNotEmpty("", "name")(); msg != "name must not be blank" {
		t.Error("Wrong message:", msg)
	} else if valid {
		t.Error("It should not be valid.")
	}
}

func TestValidate_APIVersion(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		APIVersion string
		Valid      bool
	}{
		{"2.5", true},
		{"2.4", true},

		{"1.5", false},
		{"3.5", false},
		{".5", false},
		{".", false},
		{"7.", false},
		{"2.", false},
		{"3.5", false},
		{"2.5.0", false},
		{"asdf", false},
	}

	for i, test := range tests {
		valid, msg := validateAPIVersion(test.APIVersion)()

		if msg != `api_version must be in the form: 2.Y (Y is the minor version number)` {
			t.Errorf("%d) Wrong message: %s", i, msg)
		}

		if valid != test.Valid {
			t.Errorf("%d) %v Want: %v got: %v", i, test.APIVersion, test.Valid, valid)
		}
	}
}

func TestValidate_Listen(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		Listen string
		Valid  bool
	}{
		{":80", true},
		{":65535", true},
		{":1", true},
		{"192.168.1.1:80", true},
		{"[2001:db8::ff00:42:8329]:80", true},
		{"localhost:1", true},

		{":0", false},
		{"192.168.1.1:", false},
		{":66666", false},
		{":-66", false},
	}

	for i, test := range tests {
		valid, msg := validateListen(test.Listen)()

		if msg != `listen must be in form host:port where host is optional, and 0 < port <= 65535` {
			t.Errorf("%d) Wrong message: %s", i, msg)
		}

		if valid != test.Valid {
			t.Errorf("%d) %v Want: %v got: %v", i, test.Listen, test.Valid, valid)
		}
	}
}

func TestValidate_LenAtLeast(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		Arr     []int
		AtLeast int
		Valid   bool
	}{
		{[]int{5, 6}, 0, true},
		{[]int{5, 6}, 1, true},
		{[]int{5, 6}, 2, true},
		{[]int{5, 6}, 3, false},
	}

	for i, test := range tests {
		valid, msg := validateLenAtLeast(test.Arr, test.AtLeast, "arr")()

		if msg != fmt.Sprintf("arr must contain at least %d entries", test.AtLeast) {
			t.Errorf("%d) Wrong message: %v", i, msg)
		}

		if valid != test.Valid {
			t.Errorf("%d) %v Want: %v got: %v", i, test.Arr, test.Valid, valid)
		}
	}
}

func TestValidate_LenAtLeastPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("should have panic'd")
		}
	}()
	validateLenAtLeast([]int{}, -5, "arr")
}

func TestValidate_OptionErrors(t *testing.T) {
	t.Parallel()

	// Deep copy config
	config := *testConfig
	config.Services = make([]Service, 0, len(testConfig.Services))
	for i, s := range testConfig.Services {
		s.Plans = make([]Plan, 0, len(testConfig.Services[i].Plans))
		for _, p := range testConfig.Services[i].Plans {
			p.Metadata = nil
			s.Plans = append(s.Plans, p)
		}
		// Copy Services.Requires & Tags
		s.Requires = nil
		s.Tags = nil
		config.Services = append(config.Services, s)
	}

	err := config.validate()
	if err != nil {
		t.Error(err)
	}

	// Any one of these should trigger failures, check t.Log(err) output
	// to see all the error messages.
	config.Services[1].ID = ""
	config.Listen = "junk"
	config.Services[1].Plans[0].Description = ""

	err = config.validate()
	if err == nil {
		t.Fatal("expected a failure")
	} else if str := err.Error(); !strings.Contains(str, "listen must be in form") {
		t.Error("Expect a failure about listen:", str)
	} else if !strings.Contains(str, "catalog.services[1].id must not be blank") {
		t.Error("Expect a failure about service id:", str)
	} else if !strings.Contains(str, "catalog.services[1].plans[0].description must not be blank") {
		t.Error("Expect a failure about plan description:", str)
	}
}
