package gocfbroker

import "errors"

const (
	testUnprocessablePlanID = "unprocessable-plan-id"
)

type testService struct {
	dashboardURL  string
	provisioned   bool
	deprovisioned bool
	updated       bool
	bound         bool
	unbound       bool
}

func (t *testService) Provision(instanceID string, provisionRequest ProvisionRequest) (ProvisionResponse, error) {
	t.provisioned = true
	return ProvisionResponse{DashboardURL: t.dashboardURL}, nil
}

func (t *testService) Deprovision(instanceID, serviceID, planID string) error {
	t.deprovisioned = true
	return nil
}

func (t *testService) Update(instanceID string, req UpdateProvisionRequest) error {
	t.updated = true
	if req.PlanID == testUnprocessablePlanID {
		return errors.New("Failed to process service bits and returning an error")
	}
	return nil
}

func (t *testService) Bind(instanceID, bindingID string, req BindingRequest) (res BindingResponse, err error) {
	t.bound = true
	return res, err
}

func (t *testService) Unbind(instanceID, bindingID, serviceID, planID string) error {
	t.unbound = true
	return nil
}
