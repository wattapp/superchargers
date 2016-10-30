package web

import (
	"errors"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"github.com/wattapp/superchargers/pkg/location"
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
				locationID, _ := strconv.ParseInt(resolvedID.ID, 10, 64)
				return location.GetLocation(locationID)
			default:
				return nil, errors.New("Unknown node type")
			}
		},
		TypeResolve: func(p graphql.ResolveTypeParams) *graphql.Object {
			switch p.Value.(type) {
			case *location.Location:
				return locationType
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
					l := p.Source.(*location.Location)
					return l.Address, nil
				},
			},
			"addressLine1": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					raw := *l.AddressLine1
					if raw == "" {
						return nil, nil
					}
					return raw, nil
				},
			},
			"addressLine2": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					raw := *l.AddressLine2
					if raw == "" {
						return nil, nil
					}
					return raw, nil
				},
			},
			"addressNotes": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					raw := *l.AddressNotes
					if raw == "" {
						return nil, nil
					}
					return raw, nil
				},
			},
			"amentities": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.Amenities, nil
				},
			},
			"baiduLat": &graphql.Field{
				Type: graphql.Float,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.BaiduLat, nil
				},
			},
			"baiduLng": &graphql.Field{
				Type: graphql.Float,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.BaiduLng, nil
				},
			},
			"chargers": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					raw := *l.Chargers
					if raw == "" {
						return nil, nil
					}
					return raw, nil
				},
			},
			"city": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.City, nil
				},
			},
			"commonName": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.CommonName, nil
				},
			},
			"country": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.Country, nil
				},
			},
			"destinationChargerLogo": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					raw := *l.DestinationChargerLogo
					if raw == "" {
						return nil, nil
					}
					return raw, nil
				},
			},
			"destinationWebsite": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					raw := *l.DestinationWebsite
					if raw == "" {
						return nil, nil
					}
					return raw, nil
				},
			},
			"directionsLink": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.DirectionsLink, nil
				},
			},
			"emails": &graphql.Field{
				Type: graphql.NewList(emailType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.Emails, nil
				},
			},
			"geocode": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.Geocode, nil
				},
			},
			"hours": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					raw := *l.Hours
					if raw == "" {
						return nil, nil
					}
					return raw, nil
				},
			},
			"isGallery": &graphql.Field{
				Type: graphql.Boolean,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.IsGallery, nil
				},
			},
			"kioskPinX": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.KioskPinX, nil
				},
			},
			"kioskPinY": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.KioskPinY, nil
				},
			},
			"kioskZoomPinX": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.KioskZoomPinX, nil
				},
			},
			"kioskZoomPinY": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.KioskZoomPinY, nil
				},
			},
			"latitude": &graphql.Field{
				Type: graphql.Float,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.Latitude, nil
				},
			},
			"longitude": &graphql.Field{
				Type: graphql.Float,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.Longitude, nil
				},
			},
			"locationId": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.LocationID, nil
				},
			},
			"locationType": &graphql.Field{
				Type: graphql.NewList(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.LocationType, nil
				},
			},
			"nid": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.Nid, nil
				},
			},
			"openSoon": &graphql.Field{
				Type: graphql.Boolean,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.OpenSoon, nil
				},
			},
			"path": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.Path, nil
				},
			},
			"postalCode": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					raw := *l.PostalCode
					if raw == "" {
						return nil, nil
					}
					return raw, nil
				},
			},
			"provinceState": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					raw := *l.ProvinceState
					if raw == "" {
						return nil, nil
					}
					return raw, nil
				},
			},
			"region": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.Region, nil
				},
			},
			"salesPhone": &graphql.Field{
				Type: graphql.NewList(phoneType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.SalesPhone, nil
				},
			},
			"salesRepresentative": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.SalesRepresentative, nil
				},
			},
			"subRegion": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					raw := *l.SubRegion
					if raw == "" {
						return nil, nil
					}
					return raw, nil
				},
			},
			"title": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.Title, nil
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
					locations, err := location.Locations()
					if err != nil {
						return nil, err
					}

					args := relay.NewConnectionArguments(p.Args)

					l := []interface{}{}
					for _, location := range locations {
						l = append(l, location)
					}

					return relay.ConnectionFromArray(l, args), nil
				},
			},
			"node": nodeDefinitions.NodeField,
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
}
