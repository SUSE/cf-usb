package gocfbroker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/blang/semver"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

const (
	// RFC 6586 defines this for rate limiting but we'll use it in the case
	// where there's already an operation pending for the instanceID
	statusTooManyRequests     = 429
	statusUnprocessableEntity = 422
	paramInstanceID           = "instance_id"
	paramBindingID            = "binding_id"
	queryServiceID            = "service_id"
	queryPlanID               = "plan_id"
)

const (
	errMsgRequireAppGUID = `{"error":"RequiresApp","description":"This service supports generation of credentials through binding an application only."}`
)

var (
	logOutput io.Writer = os.Stdout
	logger              = log.New(logOutput, "", log.Ldate|log.Ltime)
)

// Broker is a CF V2 Service broker, it has a service, DB and configuration
// in order to function.
type Broker struct {
	provisioner Provisioner
	db          JSONStorer
	Options     Options
	apiVersion  semver.Version

	// jobs is queue to write jobs to
	jobs chan *brokerJob
	// jobAccepts has nil written to it on success, or an error if the job could not be queued
	jobAccepts chan error
	// jobResultReqs is the queue to write result requests to
	jobResultReqs chan jobResultRequest
}

// New creates a new service broker with a database and configuration struct.
func New(provisioner Provisioner, storer Storer, options Options) (*Broker, error) {
	if err := options.validate(); err != nil {
		return nil, err
	}

	// Because the broker API is versioned without patch versions
	if strings.Count(options.APIVersion, ".") == 1 {
		options.APIVersion += ".0"
	}

	brokerSemver, err := semver.Make(options.APIVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to init broker (bad api_version in config): %v", err)
	}

	// Create the service broker.
	broker := Broker{
		provisioner:   provisioner,
		db:            NewJSONStorer(storer),
		Options:       options,
		apiVersion:    brokerSemver,
		jobs:          make(chan *brokerJob),
		jobAccepts:    make(chan error),
		jobResultReqs: make(chan jobResultRequest),
	}

	return &broker, nil
}

// LoadConfig takes a filename of a JSON file and decodes it into the provided
// configStruct.
func LoadConfig(configFile string, configStruct interface{}) (err error) {
	f, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer f.Close()

	return LoadConfigReader(f, configStruct)
}

// LoadConfigReader takes a reader full of JSON and decodes it into the provided
// configStruct.
func LoadConfigReader(configReader io.Reader, configStruct interface{}) (err error) {
	dec := json.NewDecoder(configReader)
	if err = dec.Decode(configStruct); err != nil {
		return err
	}

	return nil
}

// Start listens on the port configured, exposing the REST API.
// This should only be called once during the life-cycle of the program.
func (b *Broker) Start() {
	// Set up logging system

	handler := b.buildAPI(b.buildRoutes())

	log.Println("listening:", b.Options.Listen)
	go b.dispatcher()
	log.Println(http.ListenAndServe(b.Options.Listen, handler))
}

func (b *Broker) buildAPI(router http.Handler) http.Handler {
	mainchain := alice.New(loggingMW, recoverMW, jsonMiddleware, b.basicAuthMiddleware, b.apiVersionMiddleware)
	return mainchain.Then(router)
}

func (b *Broker) buildRoutes() http.Handler {
	router := httprouter.New()

	router.GET("/v2/catalog", errorMiddleware(b.fetchCatalog))
	router.PUT("/v2/service_instances/:instance_id", asyncMiddleware(errorMiddleware(b.provisionInstance)))
	router.DELETE("/v2/service_instances/:instance_id", asyncMiddleware(errorMiddleware(b.deprovisionInstance)))
	router.PATCH("/v2/service_instances/:instance_id", asyncMiddleware(errorMiddleware(b.updateInstance)))
	router.PUT("/v2/service_instances/:instance_id/service_bindings/:binding_id", errorMiddleware(b.bindInstance))
	router.DELETE("/v2/service_instances/:instance_id/service_bindings/:binding_id", errorMiddleware(b.unbindInstance))
	router.GET("/v2/service_instances/:instance_id/last_operation", errorMiddleware(b.lastOperation))

	//router.NotFound = func(w http.ResponseWriter, r *http.Request) {
	//	writeJSONError(w, http.StatusNotFound, "api endpoint does not exist")
	//}

	return router
}

// fetchCatalog from the broker
// GET /v2/catalog
func (b *Broker) fetchCatalog(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	writeJSON(w, http.StatusOK, b.Options.Catalog)
	return nil
}

// provisionInstance of the service
// PUT /v2/service_instances/:instance_id
func (b *Broker) provisionInstance(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	instanceID, ok := ensureParam(w, r, p, paramInstanceID)
	if !ok {
		return nil
	}

	var provisionReq ProvisionRequest
	if err := readJSON(r, &provisionReq); err != nil {
		writeJSONError(w, http.StatusBadRequest, "could not parse request")
		return nil
	}

	if errors := provisionReq.validate(b.Options.Catalog); errors != nil {
		writeJSONError(w, http.StatusBadRequest, validationErrors(errors).Error())
		return nil
	}
	if r.URL.Query().Get(queryAcceptsIncomplete) != "true" {

		result := b.provisionSync(instanceID, provisionReq)
		if result.Status != jobStatusFailed {
			writeJSON(w, http.StatusCreated, nil)
			return nil
		} else {
			return errors.New(result.ErrorMsg)
		}
	} else {
		err := b.queueJob(brokerJob{
			JobKind:          jobKindProvision,
			InstanceID:       instanceID,
			ServiceID:        provisionReq.ServiceID,
			PlanID:           provisionReq.PlanID,
			OrganizationGUID: provisionReq.OrganizationGUID,
			SpaceGUID:        provisionReq.SpaceGUID,
			JobResult: jobResult{
				Status: jobStatusInProgress,
			},
		})
		if err == errJobExists {
			writeJSONError(w, statusTooManyRequests, err.Error())
			return nil
		} else if err != nil {
			return err
		}

		writeJSON(w, http.StatusAccepted, nil)
	}
	return nil
}

// deprovisionInstance of the service
// DELETE /v2/service_instances/:instance_id
func (b *Broker) deprovisionInstance(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	instanceID, ok := ensureParam(w, r, p, paramInstanceID)
	if !ok {
		return nil
	}

	// Optional service / plan id query parameters can be provided as hints to the broker.
	query := r.URL.Query()
	serviceID := query.Get("service_id")
	planID := query.Get("plan_id")

	var deprovisionReq deprovisionRequest
	deprovisionReq.planID = planID
	deprovisionReq.serviceID = serviceID

	if r.URL.Query().Get(queryAcceptsIncomplete) != "true" {
		result := b.deprovisionSync(instanceID, deprovisionReq)
		if result.Status != jobStatusFailed {
			writeJSON(w, http.StatusOK, nil)
			return nil
		} else {
			return errors.New(result.ErrorMsg)
		}
	} else {
		err := b.queueJob(brokerJob{
			JobKind:    jobKindDeprovision,
			InstanceID: instanceID,
			ServiceID:  serviceID,
			PlanID:     planID,
			JobResult: jobResult{
				Status: jobStatusInProgress,
			},
		})
		if err == errJobExists {
			writeJSONError(w, statusTooManyRequests, err.Error())
			return nil
		} else if err != nil {
			return err
		}

		writeJSON(w, http.StatusAccepted, nil)
	}
	return nil
}

// updateInstance plan
// PATCH /v2/service_instances/:instance_id
func (b *Broker) updateInstance(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	instanceID, ok := ensureParam(w, r, p, paramInstanceID)
	if !ok {
		return nil
	}

	var updateProvisionReq UpdateProvisionRequest
	if err := readJSON(r, &updateProvisionReq); err != nil {
		writeJSONError(w, http.StatusBadRequest, "could not parse request")
		return nil
	}

	if errors := updateProvisionReq.validate(b.Options.Catalog); errors != nil {
		writeJSONError(w, http.StatusBadRequest, validationErrors(errors).Error())
		return nil
	}
	if r.URL.Query().Get(queryAcceptsIncomplete) != "true" {
		result := b.updateSync(instanceID, updateProvisionReq)
		if result.Status != jobStatusFailed {
			writeJSON(w, http.StatusCreated, nil)
			return nil
		} else {
			return errors.New(result.ErrorMsg)
		}
	} else {

		err := b.queueJob(brokerJob{
			JobKind:    jobKindUpdate,
			InstanceID: instanceID,
			PlanID:     updateProvisionReq.PlanID,
			JobResult: jobResult{
				Status: jobStatusInProgress,
			},
		})
		if err == errJobExists {
			writeJSONError(w, statusTooManyRequests, err.Error())
			return nil
		} else if err != nil {
			return err
		}

		writeJSON(w, http.StatusAccepted, nil)
	}
	return nil
}

// bindInstance adds a binding to this service instance
// PUT /v2/service_instances/:instance_id/service_bindings/:binding_id
func (b *Broker) bindInstance(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	instanceID, ok := ensureParam(w, r, p, paramInstanceID)
	if !ok {
		return nil
	}
	bindingID, ok := ensureParam(w, r, p, paramBindingID)
	if !ok {
		return nil
	}

	var bindingReq BindingRequest
	if err := readJSON(r, &bindingReq); err != nil {
		writeJSONError(w, http.StatusBadRequest, "could not parse request")
		return nil
	}

	if errors := bindingReq.validate(b.Options.Catalog); errors != nil {
		writeJSONError(w, http.StatusBadRequest, validationErrors(errors).Error())
		return nil
	}

	if len(bindingReq.AppGUID) == 0 && b.Options.RequireAppGUID {
		w.WriteHeader(statusUnprocessableEntity)
		_, err := io.WriteString(w, errMsgRequireAppGUID)
		return err
	}

	key := []byte(b.Options.EncryptionKey)
	// Check if a service binding already exists
	dbBindReq, dbBindRes, err := getBinding(b.db, key, instanceID, bindingID)
	if err == nil {
		if dbBindReq.Equal(bindingReq) {
			writeJSON(w, http.StatusOK, dbBindRes)
		} else {
			writeJSONError(w, http.StatusConflict, "a different binding exists for this binding_id")
		}
		return nil
	} else if !IsKeyNotExist(err) {
		return err
	}

	// Bind the service instance
	dbBindRes, err = b.provisioner.Bind(instanceID, bindingID, bindingReq)
	if err != nil {
		return err
	}

	// Store the binding in the database
	if err = putBinding(b.db, key, instanceID, bindingID, bindingReq, dbBindRes); err != nil {
		return err
	}

	writeJSON(w, http.StatusCreated, dbBindRes)
	return nil
}

// unbindInstance
// DELETE /v2/service_instances/:instance_id/service_bindings/:binding_id
func (b *Broker) unbindInstance(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	instanceID, ok := ensureParam(w, r, p, paramInstanceID)
	if !ok {
		return nil
	}
	bindingID, ok := ensureParam(w, r, p, paramBindingID)
	if !ok {
		return nil
	}

	// Optional service / plan id query parameters can be provided as hints to the broker.
	query := r.URL.Query()
	serviceID := query.Get(queryServiceID)
	planID := query.Get(queryPlanID)

	// Check if the binding exists in the broker database first.
	err := delBinding(b.db, instanceID, bindingID)
	if IsKeyNotExist(err) {
		writeJSON(w, http.StatusGone, nil)
		return nil
	} else if err != nil {
		return err
	}

	if err = b.provisioner.Unbind(instanceID, bindingID, serviceID, planID); err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, nil)
	return nil
}

// lastOperation
// GET /v2/service_instances/:instance_id/last_operation
func (b *Broker) lastOperation(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	instanceID, ok := ensureParam(w, r, p, paramInstanceID)
	if !ok {
		return nil
	}

	result := b.jobResult(instanceID)

	if result.ErrorMsg == errMsgNoJobExists {
		writeJSONError(w, http.StatusNotFound, errMsgNoJobExists)
		return nil
	}

	var resp = struct {
		State       string `json:"state"`
		Description string `json:"description,omitempty"`
	}{
		State:       string(result.Status),
		Description: result.ErrorMsg,
	}

	writeJSON(w, http.StatusOK, resp)
	return nil
}

// ensureParam writes out an error response if the route parameter is not present.
func ensureParam(w http.ResponseWriter, r *http.Request, p httprouter.Params, param string) (value string, present bool) {
	value = p.ByName(param)
	if len(value) == 0 {
		writeJSONError(w, http.StatusBadRequest, "must provide parameter: "+param)
		return value, present
	}
	present = true
	return value, present
}
