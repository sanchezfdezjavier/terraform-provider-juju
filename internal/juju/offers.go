package juju

import (
	"fmt"

	"github.com/juju/juju/api/client/applicationoffers"
	"github.com/juju/juju/core/crossmodel"
	"github.com/juju/juju/rpc/params"
)

type offersClient struct {
	ConnectionFactory
}

type CreateOfferInput struct {
	ApplicationName string
	Endpoint        string
	ModelName       string
	ModelUUID       string
	Name            string
}

type CreateOfferResponse struct {
	Name     string
	OfferURL string
}

type DestroyOfferInput struct {
	OfferURL string
}

func newOffersClient(cf ConnectionFactory) *offersClient {
	return &offersClient{
		ConnectionFactory: cf,
	}
}

func (c offersClient) CreateOffer(input *CreateOfferInput) (*CreateOfferResponse, []error) {
	var errs []error

	conn, err := c.GetConnection(nil)
	if err != nil {
		return nil, append(errs, err)
	}

	client := applicationoffers.NewClient(conn)
	defer client.Close()

	offerName := input.Name
	if offerName == "" {
		offerName = input.ApplicationName
	}

	result, err := client.Offer(input.ModelUUID, input.ApplicationName, []string{input.Endpoint}, "admin", offerName, "")
	if err != nil {
		return nil, append(errs, err)
	}

	for _, v := range result {
		var result = params.ErrorResult{}
		if v == result {
			continue
		} else {
			errs = append(errs, v.Error)
		}
	}
	if len(errs) != 0 {
		return nil, errs
	}

	filter := crossmodel.ApplicationOfferFilter{
		OfferName: offerName,
		ModelName: input.ModelName,
	}

	offer, err := findApplicationOffers(client, filter)
	if err != nil {
		return nil, append(errs, err)
	}

	resp := CreateOfferResponse{
		Name:     offer.OfferName,
		OfferURL: offer.OfferURL,
	}
	return &resp, nil
}

func (c offersClient) DestroyOffer(input *DestroyOfferInput) error {
	conn, err := c.GetConnection(nil)
	if err != nil {
		return err
	}

	client := applicationoffers.NewClient(conn)
	defer client.Close()

	//TODO: verify destruction after attaching
	forceDestroy := false
	err = client.DestroyOffers(forceDestroy, input.OfferURL)
	if err != nil {
		return err
	}

	return nil
}

func findApplicationOffers(client *applicationoffers.Client, filter crossmodel.ApplicationOfferFilter) (*crossmodel.ApplicationOfferDetails, error) {
	offers, err := client.FindApplicationOffers(filter)
	if err != nil {
		return nil, err
	}

	if len(offers) > 1 || len(offers) == 0 {
		return nil, fmt.Errorf("unable to find offer after creation")
	}

	return offers[0], nil
}
