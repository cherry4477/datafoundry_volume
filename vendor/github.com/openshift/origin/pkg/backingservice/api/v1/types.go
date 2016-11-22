package v1

import (
	"k8s.io/kubernetes/pkg/api/unversioned"
	kapi "k8s.io/kubernetes/pkg/api/v1"
)

// BackingService describe a BackingService
type BackingService struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	kapi.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the behavior of the Namespace.
	Spec BackingServiceSpec `json:"spec,omitempty" description:"specification of the desired behavior for a BackingService"`

	// Status describes the current status of a Namespace
	Status BackingServiceStatus `json:"status,omitempty" description:"status describes the current status of a BackingService"`
}

// BackingServiceList describe a list of BackingService
type BackingServiceList struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	unversioned.ListMeta `json:"metadata,omitempty"`

	// Items is a list of routes
	Items []BackingService `json:"items" description:"list of backingservice"`
}

// BackingServiceSpec describe the attributes on a Backingservice
type BackingServiceSpec struct {
	// name of backingservice
	Name string `json:"name" description:"name of backingservice"`
	// id of backingservice
	Id string `json:"id" description:"id of backingservice"`
	// description of a backingservice
	Description string `json:"description" description:"description of a backingservice"`
	// is backingservice bindable
	Bindable bool `json:"bindable" description:"is backingservice bindable?"`
	// is  backingservice plan updateable
	PlanUpdateable bool `json:"plan_updateable, omitempty" description:"is  backingservice plan updateable"`
	// list of backingservice tags of BackingService
	Tags []string `json:"tags, omitempty" description:"list of backingservice tags of BackingService"`
	// require condition of backingservice
	Requires []string `json:"requires, omitempty" description:"require condition of backingservice"`

	// metadata of backingservice
	Metadata map[string]string `json:"metadata, omitempty" description:"metadata of backingservice"`
	// plans of a backingservice
	Plans []ServicePlan `json:"plans" description:"plans of a backingservice"`
	// DashboardClient of backingservic
	DashboardClient map[string]string `json:"dashboard_client" description:"DashboardClient of backingservice"`
}

// ServiceMetadata describe a ServiceMetadata
type ServiceMetadata struct {
	// displayname of a ServiceMetadata
	DisplayName string `json:"displayName, omitempty"`
	// imageurl of a ServiceMetadata
	ImageUrl string `json:"imageUrl, omitempty"`
	// long description of a ServiceMetadata
	LongDescription string `json:"longDescription, omitempty"`
	// providerdisplayname of a ServiceMetadata
	ProviderDisplayName string `json:"providerDisplayName, omitempty"`
	// documrntation url of a ServiceMetadata
	DocumentationUrl string `json:"documentationUrl, omitempty"`
	// support url of a ServiceMetadata
	SupportUrl string `json:"supportUrl, omitempty"`
}

// ServiceDashboardClient describe a ServiceDashboardClient
type ServiceDashboardClient struct {
	// id of a ServiceDashboardClient
	Id string `json:"id, omitempty"`
	// secret of a ServiceDashboardClient
	Secret string `json:"secret, omitempty"`
	//redirect uri of a ServiceDashboardClient
	RedirectUri string `json:"redirect_uri, omitempty"`
}

// ServicePlan describe a ServicePlan
type ServicePlan struct {
	// name of a ServicePlan
	Name string `json:"name"`
	//id of a ServicePlan
	Id string `json:"id"`
	// description of a ServicePlan
	Description string `json:"description"`
	// metadata of a ServicePlan
	Metadata ServicePlanMetadata `json:"metadata, omitempty"`
	// is this plan free or not
	Free bool `json:"free, omitempty"`
}

// ServicePlanMetadata describe a ServicePlanMetadata
type ServicePlanMetadata struct {
	// bullets of a ServicePlanMetadata
	Bullets []string `json:"bullets, omitempty"`
	// costs of a ServicePlanMetadata
	Costs []ServicePlanCost `json:"costs, omitempty"`
	// displayname of a ServicePlanMetadata
	DisplayName string `json:"displayName, omitempty"`
}

//TODO amount should be a array object...

// ServicePlanCost describe a ServicePlanCost
type ServicePlanCost struct {
	// amount of a ServicePlanCost
	Amount map[string]float64 `json:"amount, omitempty"`
	// unit of a ServicePlanCost
	Unit string `json:"unit, omitempty"`
}

// ProjectStatus is information about the current status of a Project
type BackingServiceStatus struct {
	// phase is the current lifecycle phase of the servicebroker
	Phase BackingServicePhase `json:"phase,omitempty" description:"phase is the current lifecycle phase of the servicebroker"`
}

type BackingServicePhase string

const (
	BackingServicePhaseActive   BackingServicePhase = "Active"
	BackingServicePhaseInactive BackingServicePhase = "Inactive"
)
