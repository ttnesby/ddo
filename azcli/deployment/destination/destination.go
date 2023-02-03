package destination

import (
	"github.com/google/uuid"
)

type Destination interface {
	AzCli() (destination string, parameters []string)
}

type base struct {
	id           uuid.UUID
	rgOrLocation string
}

func (b base) AzCli(dest, idTag, tag string) (destination string, parameters []string) {
	return dest, []string{idTag, b.id.String(), tag, b.rgOrLocation}
}

type resourceGroup base

func (rg resourceGroup) AzCli() (destination string, parameters []string) {
	return base(rg).AzCli("group", "--subscription", "--resource-group")
}

type subscription base

func (s subscription) AzCli() (destination string, parameters []string) {
	return base(s).AzCli("sub", "--subscription", "--location")
}

type managementGroup base

func (mg managementGroup) AzCli() (destination string, parameters []string) {
	return base(mg).AzCli("mg", "--management-group-id", "--location")
}

type creator func(id uuid.UUID, rgOrLocation string) Destination

func newDestination(id string, rgOrLocation string, create creator) (Destination, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return create(uid, rgOrLocation), nil
}
func ResourceGroup(name, inSubscriptionId string) (Destination, error) {
	create := func(id uuid.UUID, rgOrLocation string) Destination {
		return resourceGroup{id, rgOrLocation}
	}
	return newDestination(inSubscriptionId, name, create)
}

func Subscription(id, location string) (Destination, error) {
	create := func(id uuid.UUID, rgOrLocation string) Destination {
		return subscription{id, rgOrLocation}
	}
	return newDestination(id, location, create)
}

func ManagementGroup(id, location string) (Destination, error) {
	create := func(id uuid.UUID, rgOrLocation string) Destination {
		return managementGroup{id, rgOrLocation}
	}
	return newDestination(id, location, create)
}
