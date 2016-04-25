package mgmt

import (
	"github.com/frodenas/brokerapi"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/genmodel"
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/ccapi"
	. "github.com/hpcloud/cf-usb/lib/operations"
	"github.com/pivotal-golang/lager"
)

// This file is safe to edit. Once it exists it will not be overwritten

const defaultBrokerName string = "usb"

func ConfigureServicePlanAPI(api *UsbMgmtAPI, auth authentication.AuthenticationInterface,
	configProvider config.ConfigProvider, ccServiceBroker ccapi.ServiceBrokerInterface,
	logger lager.Logger) {
	log := logger.Session("usb-mgmt-serviceplan")

	api.GetServicePlanHandler = GetServicePlanHandlerFunc(func(params GetServicePlanParams, principal interface{}) middleware.Responder {
		log := log.Session("get-service-plan")
		log.Info("request", lager.Data{"plan-id": params.PlanID})

		planInfo, dialID, _, err := configProvider.GetPlan(params.PlanID)
		if err != nil {
			return &GetServicePlanInternalServerError{Payload: err.Error()}
		}
		if planInfo == nil {
			return &GetServicePlanNotFound{}
		}

		plan := &genmodel.Plan{
			Name:        &planInfo.Name,
			ID:          planInfo.ID,
			DialID:      &dialID,
			Description: planInfo.Description,
			Free:        planInfo.Free,
		}

		return &GetServicePlanOK{Payload: plan}

	})

	api.UpdateServicePlanHandler = UpdateServicePlanHandlerFunc(func(params UpdateServicePlanParams, principal interface{}) middleware.Responder {
		log := log.Session("update-service-plan")
		log.Info("request", lager.Data{"plan-id": params.PlanID})

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &UpdateServicePlanInternalServerError{Payload: err.Error()}
		}

		plan, dialID, _, err := configProvider.GetPlan(params.PlanID)
		if err != nil {
			return &UpdateServicePlanInternalServerError{Payload: err.Error()}
		}
		if plan == nil {
			return &UpdateServicePlanNotFound{}
		}

		brokerName := defaultBrokerName
		if len(config.ManagementAPI.BrokerName) > 0 {
			brokerName = config.ManagementAPI.BrokerName
		}

		guid, err := ccServiceBroker.GetServiceBrokerGuidByName(brokerName)
		if err != nil {
			log.Error("get-service-broker-failed", err)
			return &UpdateServiceInternalServerError{Payload: err.Error()}
		}

		err = ccServiceBroker.Update(guid, brokerName, config.BrokerAPI.ExternalUrl, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
		if err != nil {
			log.Error("update-service-broker-failed", err)
			return &UpdateServiceInternalServerError{Payload: err.Error()}
		}

		dial, instanceID, err := configProvider.GetDial(dialID)
		if err != nil {
			return &UpdateServicePlanInternalServerError{Payload: err.Error()}
		}
		if dial == nil {
			return &GetDialNotFound{}
		}

		instance, _, err := configProvider.GetDriverInstance(instanceID)
		if err != nil {
			return &UpdateServicePlanInternalServerError{Payload: err.Error()}
		}
		if instance == nil {
			return &GetDriverInstanceNotFound{}
		}

		var meta brokerapi.ServicePlanMetadata
		if params.Plan.Description != "" {
			plan.Description = params.Plan.Description
		}
		if params.Plan.Name != nil {
			plan.Name = *params.Plan.Name
			meta.DisplayName = *params.Plan.Name
		}
		plan.Free = params.Plan.Free

		plan.Metadata = &meta
		dial.Plan = *plan
		err = configProvider.SetDial(instanceID, dialID, *dial)

		if err != nil {
			return &UpdateServicePlanInternalServerError{Payload: err.Error()}
		}
		err = ccServiceBroker.EnableServiceAccess(instance.Service.Name)
		if err != nil {
			log.Error("enable-service-access-failed", err)
			return &UpdateServiceInternalServerError{Payload: err.Error()}
		}

		return &UpdateServicePlanOK{Payload: params.Plan}
	})
	api.GetServicePlansHandler = GetServicePlansHandlerFunc(func(params GetServicePlansParams, principal interface{}) middleware.Responder {
		log := log.Session("get-service-plans")
		log.Info("request", lager.Data{"driver-instance-id": params.DriverInstanceID})

		if *params.DriverInstanceID == "" {
			return &GetServicePlansInternalServerError{Payload: "Empty driver instance id in get service plans handler"}
		}
		var servicePlans = make([]*genmodel.Plan, 0)

		instanceInfo, err := configProvider.LoadDriverInstance(*params.DriverInstanceID)
		if err != nil {
			return &GetServicePlansInternalServerError{Payload: err.Error()}
		}

		for diaID, dia := range instanceInfo.Dials {
			planID := dia.Plan.ID
			description := dia.Plan.Description
			free := dia.Plan.Free
			plan := &genmodel.Plan{
				Name:        &dia.Plan.Name,
				ID:          planID,
				DialID:      &diaID,
				Description: description,
				Free:        free,
			}

			servicePlans = append(servicePlans, plan)
		}

		log.Debug("", lager.Data{"service-plans-found": len(servicePlans)})

		return &GetServicePlansOK{Payload: servicePlans}
	})

}
