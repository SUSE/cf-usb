package gocfbroker

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/kat-co/vala"
)

var (
	rgxAPIVersion = regexp.MustCompile(`^2\.[0-9]+$`)
)

func validateServiceID(c Catalog, serviceID string) vala.Checker {
	return func() (valid bool, msg string) {
		msg = "service_id does not reference a service from the catalog"
		for _, service := range c.Services {
			if service.ID == serviceID {
				valid = true
				break
			}
		}
		return valid, msg
	}
}

func validatePlanID(c Catalog, serviceID, planID string) vala.Checker {
	return func() (valid bool, msg string) {
		msg = "plan_id does not reference a services plan from the catalog"
		for _, service := range c.Services {
			if service.ID != serviceID {
				continue
			}
			for _, plan := range service.Plans {
				if plan.ID == planID {
					valid = true
					break
				}
			}
			break
		}

		return valid, msg
	}
}

func validateServiceUpdatable(c Catalog, serviceID string) vala.Checker {
	return func() (valid bool, msg string) {
		msg = "service plan is not updatable"
		for _, service := range c.Services {
			if service.ID == serviceID {
				valid = service.PlanUpdatable
				break
			}
		}

		return valid, msg
	}
}

func validateAPIVersion(apiVersion string) vala.Checker {
	return func() (valid bool, msg string) {
		msg = "api_version must be in the form: 2.Y (Y is the minor version number)"
		valid = rgxAPIVersion.MatchString(apiVersion)
		return valid, msg
	}
}

func validateListen(listen string) vala.Checker {
	return func() (valid bool, msg string) {
		msg = "listen must be in form host:port where host is optional, and 0 < port <= 65535"

		_, portString, err := net.SplitHostPort(listen)
		if err != nil {
			return false, msg
		}

		if port, err := strconv.ParseUint(portString, 10, 16); err != nil {
			return false, msg
		} else if port == 0 {
			return false, msg
		}

		return true, msg
	}
}

func validateLenAtLeast(param interface{}, minLength int, paramName string) vala.Checker {
	if minLength < 0 {
		panic("minLength cannot be less than 0")
	}

	return func() (hasLen bool, errMsg string) {
		hasLen = reflect.ValueOf(param).Len() >= minLength
		errMsg = fmt.Sprintf("%s must contain at least %d entries", paramName, minLength)
		return hasLen, errMsg
	}
}

func validateStrNotEmpty(param string, paramName string) vala.Checker {
	return func() (notEmpty bool, msg string) {
		notEmpty = len(param) > 0
		msg = fmt.Sprintf("%s must not be blank", paramName)
		return notEmpty, msg
	}
}

// validate all important fields in the config
func (o *Options) validate() error {
	valid := vala.BeginValidation().Validate(
		validateAPIVersion(o.APIVersion),
		validateStrNotEmpty(o.AuthUser, "auth_user"),
		validateStrNotEmpty(o.AuthPassword, "auth_password"),
		validateListen(o.Listen),

		// Encryption key
		validateStrNotEmpty(o.EncryptionKey, "db_encryption_key"),

		// Catalog
		validateLenAtLeast(o.Catalog.Services, 1, "catalog.services"),
	)

	// Services in catalog
	for i, service := range o.Catalog.Services {
		context := fmt.Sprintf("catalog.services[%d].", i)

		valid = valid.Validate(
			validateStrNotEmpty(service.ID, context+"id"),
			validateStrNotEmpty(service.Name, context+"name"),
			validateStrNotEmpty(service.Description, context+"description"),
			validateLenAtLeast(service.Plans, 1, context+"plans"),
		)

		// Plans in the service
		for j, plan := range service.Plans {
			planContext := fmt.Sprintf("%splans[%d].", context, j)
			valid = valid.Validate(
				validateStrNotEmpty(plan.ID, planContext+"id"),
				validateStrNotEmpty(plan.Name, planContext+"name"),
				validateStrNotEmpty(plan.Description, planContext+"description"),
			)
		}
	}

	return validationErrors(valid)
}

// validationErrors returns the errors from the validation step in a single error
// ready to display to the user.
func validationErrors(v *vala.Validation) (err error) {
	if v != nil {
		err = errors.New(strings.Join(v.Errors, "\n"))
	}
	return err
}
