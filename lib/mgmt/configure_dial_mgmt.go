package mgmt

import (
	"encoding/json"

	"github.com/frodenas/brokerapi"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/genmodel"
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/ccapi"
	. "github.com/hpcloud/cf-usb/lib/operations"
	"github.com/pivotal-golang/lager"
	//	"github.com/xeipuuv/gojsonschema"

	uuid "github.com/satori/go.uuid"
)

func ConfigureDialAPI(api *UsbMgmtAPI, auth authentication.AuthenticationInterface,
	configProvider config.ConfigProvider, ccServiceBroker ccapi.ServiceBrokerInterface,
	logger lager.Logger) {

	log := logger.Session("usb-mgmt-driver")

	api.CreateDialHandler = CreateDialHandlerFunc(func(params CreateDialParams, principal interface{}) middleware.Responder {
		log := log.Session("create-dial")
		log.Info("request", lager.Data{"driver-instance-id": params.Dial.InstanceID})

		var dial config.Dial

		dialID := uuid.NewV4().String()
		dialconfig, err := json.Marshal(params.Dial.Configuration)
		if err != nil {
			return &CreateDialInternalServerError{Payload: err.Error()}
		}

		configuration := json.RawMessage(dialconfig)
		dial.Configuration = &configuration

		_, _, err = configProvider.GetInstance(*params.Dial.InstanceID)
		if err != nil {
			return &CreateDialInternalServerError{Payload: err.Error()}
		}

		var defaultPlan brokerapi.ServicePlan
		var defaultMeta brokerapi.ServicePlanMetadata

		defaultPlan.ID = uuid.NewV4().String()
		defaultMeta.DisplayName = "Plan-" + defaultPlan.ID[:8]

		defaultPlan.Description = "N/A"
		defaultPlan.Name = "plan-" + defaultPlan.ID[:8]
		defaultPlan.Free = false

		defaultPlan.Metadata = &defaultMeta
		dial.Plan = defaultPlan

		err = configProvider.SetDial(*params.Dial.InstanceID, dialID, dial)
		if err != nil {
			return &CreateDialInternalServerError{Payload: err.Error()}
		}

		params.Dial.ID = dialID
		params.Dial.Plan = defaultPlan.ID

		instance, _, err := configProvider.GetInstance(*params.Dial.InstanceID)
		if err != nil {
			return &CreateDialInternalServerError{Payload: err.Error()}
		}

		err = ccServiceBroker.EnableServiceAccess(instance.Service.Name)
		if err != nil {
			log.Error("enable-service-access-failed", err)
			return &UpdateServiceInternalServerError{Payload: err.Error()}
		}

		return &CreateDialCreated{Payload: params.Dial}
	})

	api.UpdateDialHandler = UpdateDialHandlerFunc(func(params UpdateDialParams, principal interface{}) middleware.Responder {
		log := log.Session("update-dial")
		log.Info("request", lager.Data{"driver-instance-id": params.Dial.InstanceID})

		dialID := params.DialID

		dial, _, err := configProvider.GetDial(dialID)
		if err != nil {
			log.Error("update-dial", err, lager.Data{"dial-id": dialID})
			return &UpdateDialInternalServerError{Payload: err.Error()}
		}
		if dial == nil {
			return &UpdateDialNotFound{}
		}

		dialconfig, err := json.Marshal(params.Dial.Configuration)
		if err != nil {
			return &UpdateDialInternalServerError{Payload: err.Error()}
		}

		configuration := json.RawMessage(dialconfig)
		dial.Configuration = &configuration

		_, _, err = configProvider.GetInstance(*params.Dial.InstanceID)
		if err != nil {
			return &UpdateDialInternalServerError{Payload: err.Error()}
		}

		err = configProvider.SetDial(*params.Dial.InstanceID, dialID, *dial)
		if err != nil {
			return &UpdateDialInternalServerError{Payload: err.Error()}
		}

		return &UpdateDialOK{Payload: params.Dial}
	})

	api.GetAllDialsHandler = GetAllDialsHandlerFunc(func(params GetAllDialsParams, principal interface{}) middleware.Responder {
		log := log.Session("get-dials")
		log.Info("request", lager.Data{"driver-instance-id": params.InstanceID})

		var dials = make([]*genmodel.Dial, 0)
		if *params.InstanceID == "" {
			return &GetAllDialsInternalServerError{Payload: "Empty driver instance id in get all dials"}
		}
		instanceInfo, err := configProvider.LoadDriverInstance(*params.InstanceID)
		if err != nil {
			return &GetAllDialsInternalServerError{Payload: err.Error()}
		}

		for diaID, dia := range instanceInfo.Dials {
			dialID := diaID
			planID := dia.Plan.ID
			dial := &genmodel.Dial{
				Configuration: dia.Configuration,
				InstanceID:    params.InstanceID,
				ID:            dialID,
				Plan:          planID,
			}

			dials = append(dials, dial)
		}

		log.Debug("", lager.Data{"dials-found": len(dials)})

		return &GetAllDialsOK{Payload: dials}
	})

	api.DeleteDialHandler = DeleteDialHandlerFunc(func(params DeleteDialParams, principal interface{}) middleware.Responder {
		log := log.Session("delete-dial")
		log.Info("request", lager.Data{"dial-id": params.DialID})

		dial, _, err := configProvider.GetDial(params.DialID)
		if err != nil {
			return &DeleteDialInternalServerError{Payload: err.Error()}
		}
		if dial == nil {
			return &DeleteDialNotFound{}
		}

		err = configProvider.DeleteDial(params.DialID)
		if err != nil {
			return &DeleteDialInternalServerError{Payload: err.Error()}
		}
		return &DeleteDialNoContent{}
	})

	api.GetDialHandler = GetDialHandlerFunc(func(params GetDialParams, principal interface{}) middleware.Responder {
		log := log.Session("get-dial")
		log.Info("request", lager.Data{"dial-id": params.DialID})

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &GetDialInternalServerError{Payload: err.Error()}
		}

		for diID, di := range config.Instances {
			for diaID, dia := range di.Dials {
				if diaID == params.DialID {

					var conf map[string]interface{}
					err := json.Unmarshal(*dia.Configuration, &conf)
					if err != nil {
						return &GetDialInternalServerError{Payload: err.Error()}
					}

					dial := &genmodel.Dial{
						Configuration: conf,
						InstanceID:    &diID,
						ID:            diaID,
						Plan:          dia.Plan.ID,
					}

					return &GetDialOK{Payload: dial}
				}
			}
		}

		return &GetDialNotFound{}
	})

}
