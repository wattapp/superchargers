package web

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"github.com/wattapp/superchargers/pkg/supercharger"
	"golang.org/x/net/context"
)

var nodeDefinitions *relay.NodeDefinitions

// Each top level type
var locationType *graphql.Object
var emailType *graphql.Object
var phoneType *graphql.Object

// Custom node field types
var enumChargerType = graphql.NewEnum(graphql.EnumConfig{
	Name: "Type",
	Values: graphql.EnumValueConfigMap{
		"SUPERCHARGER": &graphql.EnumValueConfig{
			Value: "SUPERCHARGER",
		},
		"DESTINATION": &graphql.EnumValueConfig{
			Value: "DESTINATION",
		},
	},
})

var typeFields = graphql.FieldConfigArgument{
	"type": &graphql.ArgumentConfig{
		Type: graphql.NewList(enumChargerType),
	},
}
var typeArgs = relay.NewConnectionArgs(typeFields)

func BuildSchema() (graphql.Schema, error) {
	nodeDefinitions = relay.NewNodeDefinitions(relay.NodeDefinitionsConfig{
		IDFetcher: func(id string, info graphql.ResolveInfo, ctx context.Context) (interface{}, error) {
			resolvedID := relay.FromGlobalID(id)

			switch resolvedID.Type {
			case "Location":
				// tripID, _ := strconv.ParseInt(resolvedID.ID, 10, 64)
				// return vehicle.GetTrip(tripID)
				return nil, errors.New("Not implemented")
			default:
				return nil, errors.New("Unknown node type")
			}
		},
		TypeResolve: func(p graphql.ResolveTypeParams) *graphql.Object {
			switch p.Value.(type) {
			default:
				return locationType
			}
		},
	})

	emailType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Email",
		Fields: graphql.Fields{
			"email": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					e := p.Source.(supercharger.Email)
					return e.Email, nil
				},
			},
			"label": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					e := p.Source.(supercharger.Email)
					return e.Label, nil
				},
			},
		},
	})

	phoneType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Phone",
		Fields: graphql.Fields{
			"number": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					e := p.Source.(supercharger.Phone)
					return e.Number, nil
				},
			},
			"label": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					e := p.Source.(supercharger.Phone)
					return e.Label, nil
				},
			},
		},
	})

	locationType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Location",
		Description: "A vehicle contains trips and charges",
		Fields: graphql.Fields{
			"id": relay.GlobalIDField("Location", nil),
			"address": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.Address, nil
				},
			},
			"addressLine1": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.AddressLine1, nil
				},
			},
			"addressLine2": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.AddressLine2, nil
				},
			},
			"addressNotes": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.AddressNotes, nil
				},
			},
			"amentities": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.Amenities, nil
				},
			},
			"baiduLat": &graphql.Field{
				Type: graphql.Float,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.BaiduLat, nil
				},
			},
			"baiduLng": &graphql.Field{
				Type: graphql.Float,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.BaiduLng, nil
				},
			},
			"chargers": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.Chargers, nil
				},
			},
			"city": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.City, nil
				},
			},
			"commonName": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.CommonName, nil
				},
			},
			"country": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.Country, nil
				},
			},
			"destinationChargerLogo": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.DestinationChargerLogo, nil
				},
			},
			"destinationWebsite": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.DestinationWebsite, nil
				},
			},
			"directionsLink": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.DirectionsLink, nil
				},
			},
			"emails": &graphql.Field{
				Type: graphql.NewList(emailType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.Emails, nil
				},
			},
			"geocode": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.Geocode, nil
				},
			},
			"hours": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.Hours, nil
				},
			},
			"isGallery": &graphql.Field{
				Type: graphql.Boolean,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.IsGallery, nil
				},
			},
			"kioskPinX": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.KioskPinX, nil
				},
			},
			"kioskPinY": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.KioskPinY, nil
				},
			},
			"kioskZoomPinX": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.KioskZoomPinX, nil
				},
			},
			"kioskZoomPinY": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.KioskZoomPinY, nil
				},
			},
			"latitude": &graphql.Field{
				Type: graphql.Float,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.Latitude, nil
				},
			},
			"longitude": &graphql.Field{
				Type: graphql.Float,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.Longitude, nil
				},
			},
			"locationId": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.LocationID, nil
				},
			},
			"locationType": &graphql.Field{
				Type: graphql.NewList(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.LocationType, nil
				},
			},
			"nid": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.Nid, nil
				},
			},
			"openSoon": &graphql.Field{
				Type: graphql.Boolean,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.OpenSoon, nil
				},
			},
			"path": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.Path, nil
				},
			},
			"postalCode": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.PostalCode, nil
				},
			},
			"provinceState": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.ProvinceState, nil
				},
			},
			"region": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.Region, nil
				},
			},
			"salesPhone": &graphql.Field{
				Type: graphql.NewList(phoneType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.SalesPhone, nil
				},
			},
			"salesRepresentative": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.SalesRepresentative, nil
				},
			},
			"subRegion": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.SubRegion, nil
				},
			},
			"title": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					v := p.Source.(supercharger.Location)
					return v.Title, nil
				},
			},
		},
		Interfaces: []*graphql.Interface{
			nodeDefinitions.NodeInterface,
		},
	})

	locationsConnectionDefinition := relay.ConnectionDefinitions(relay.ConnectionConfig{
		Name:     "Location",
		NodeType: locationType,
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"locations": &graphql.Field{
				Type: locationsConnectionDefinition.ConnectionType,
				Args: typeArgs,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					superchargers, err := supercharger.Superchargers()
					if err != nil {
						return nil, err
					}

					args := relay.NewConnectionArguments(p.Args)

					locations := []interface{}{}
					for _, location := range superchargers {
						locations = append(locations, location)
					}

					return relay.ConnectionFromArray(locations, args), nil
				},
			},
			"node": nodeDefinitions.NodeField,
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
}
