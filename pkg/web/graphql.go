package web

import (
	"errors"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"github.com/wattapp/superchargers/pkg/database"
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
var enumLocationType = graphql.NewEnum(graphql.EnumConfig{
	Name: "LocationType",
	Values: graphql.EnumValueConfigMap{
		"SUPERCHARGER": &graphql.EnumValueConfig{
			Value: "supercharger",
		},
		"STANDARD_CHARGER": &graphql.EnumValueConfig{
			Value: "standard charger",
		},
		"DESTINATION_CHARGER": &graphql.EnumValueConfig{
			Value: "destination charger",
		},
		"STORE": &graphql.EnumValueConfig{
			Value: "store",
		},
		"SERVICE": &graphql.EnumValueConfig{
			Value: "service",
		},
	},
})

var enumRegion = graphql.NewEnum(graphql.EnumConfig{
	Name: "Region",
	Values: graphql.EnumValueConfigMap{
		"NORTH_AMERICA": &graphql.EnumValueConfig{
			Value: "north_america",
		},
		"EUROPE": &graphql.EnumValueConfig{
			Value: "europe",
		},
		"ASIA_PACIFIC": &graphql.EnumValueConfig{
			Value: "asia_pacific",
		},
	},
})

var enumCountry = graphql.NewEnum(graphql.EnumConfig{
	Name: "Country",
	Values: graphql.EnumValueConfigMap{
		"ANDORRA": &graphql.EnumValueConfig{
			Value: "Andorra",
		},
		"AUSTRALIA": &graphql.EnumValueConfig{
			Value: "Australia",
		},
		"AUSTRIA": &graphql.EnumValueConfig{
			Value: "Austria",
		},
		"BELGIUM": &graphql.EnumValueConfig{
			Value: "Belgium",
		},
		"CANADA": &graphql.EnumValueConfig{
			Value: "Canada",
		},
		"CHINA": &graphql.EnumValueConfig{
			Value: "China",
		},
		"CROATIA": &graphql.EnumValueConfig{
			Value: "Croatia",
		},
		"CZECH_REPUBLIC": &graphql.EnumValueConfig{
			Value: "Czech Republic",
		},
		"DENMARK": &graphql.EnumValueConfig{
			Value: "Denmark",
		},
		"FINLAND": &graphql.EnumValueConfig{
			Value: "Finland",
		},
		"FRANCE": &graphql.EnumValueConfig{
			Value: "France",
		},
		"GERMANY": &graphql.EnumValueConfig{
			Value: "Germany",
		},
		"HONG_KONG": &graphql.EnumValueConfig{
			Value: "Hong Kong",
		},
		"ITALY": &graphql.EnumValueConfig{
			Value: "Italy",
		},
		"JAPAN": &graphql.EnumValueConfig{
			Value: "Japan",
		},
		"LIECHTENSTEIN": &graphql.EnumValueConfig{
			Value: "Liechtenstein",
		},
		"LUXEMBOURG": &graphql.EnumValueConfig{
			Value: "Luxembourg",
		},
		"MACAU": &graphql.EnumValueConfig{
			Value: "Macau",
		},
		"MEXICO": &graphql.EnumValueConfig{
			Value: "Mexico",
		},
		"NETHERLANDS": &graphql.EnumValueConfig{
			Value: "Netherlands",
		},
		"NORWAY": &graphql.EnumValueConfig{
			Value: "Norway",
		},
		"POLAND": &graphql.EnumValueConfig{
			Value: "Poland",
		},
		"SERBIA": &graphql.EnumValueConfig{
			Value: "Serbia",
		},
		"SLOVAKIA": &graphql.EnumValueConfig{
			Value: "Slovakia",
		},
		"SLOVENIA": &graphql.EnumValueConfig{
			Value: "Slovenia",
		},
		"SPAIN": &graphql.EnumValueConfig{
			Value: "Spain",
		},
		"SWEDEN": &graphql.EnumValueConfig{
			Value: "Sweden",
		},
		"SWITZERLAND": &graphql.EnumValueConfig{
			Value: "Switzerland",
		},
		"TAIWAN": &graphql.EnumValueConfig{
			Value: "Taiwan",
		},
		"UNITED_KINGDOM": &graphql.EnumValueConfig{
			Value: "United Kingdom",
		},
		"UNITED_STATES": &graphql.EnumValueConfig{
			Value: "United States",
		},
	},
})

var locationFieldArguments = relay.NewConnectionArgs(graphql.FieldConfigArgument{
	"type": &graphql.ArgumentConfig{
		Type:        graphql.NewList(enumLocationType),
		Description: "Each location may provide of 1 or many services such as supercharging, standard charging, destination charging, service, or a store.",
	},
	"region": &graphql.ArgumentConfig{
		Type: graphql.NewList(enumRegion),
	},
	"country": &graphql.ArgumentConfig{
		Type: graphql.NewList(enumCountry),
	},
	"openSoon": &graphql.ArgumentConfig{
		Type:        graphql.Boolean,
		Description: "Whether or not the location is opening soon.",
	},
	"isGallery": &graphql.ArgumentConfig{
		Type:        graphql.Boolean,
		Description: "Whether or not the location is a gallery.",
	},
	"boundingBox": &graphql.ArgumentConfig{
		Type:        graphql.NewList(graphql.Float),
		Description: "The 4 coordinates to make a bounding box in the following order: [North West Latitude, North West Longitude, South East Latitude, South East Longitude]",
	},
})

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
		Description: "A location can be a supercharger, standard charger, destination charger, service center, or a store.",
		Fields: graphql.Fields{
			"id": relay.GlobalIDField("Location", nil),
			"address": &graphql.Field{
				Type:        graphql.String,
				Description: "The precomputed address for this location including city, state, country, postal code, and region.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.Address, nil
				},
			},
			"addressLine1": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					if l.AddressLine1 == nil {
						return nil, nil
					}
					return *l.AddressLine1, nil
				},
			},
			"addressLine2": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					if l.AddressLine2 == nil {
						return nil, nil
					}
					return *l.AddressLine2, nil
				},
			},
			"addressNotes": &graphql.Field{
				Type:        graphql.String,
				Description: "Helpful human direction to find this location.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					if l.AddressNotes == nil {
						return nil, nil
					}
					return *l.AddressNotes, nil
				},
			},
			"amentities": &graphql.Field{
				Type:        graphql.String,
				Description: "The HTML representation of amentities provided by this location.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					if l.Amenities == nil {
						return nil, nil
					}
					return *l.Amenities, nil
				},
			},
			"chargers": &graphql.Field{
				Type:        graphql.String,
				Description: "The HTML representation of the chargers provided by this location.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					if l.Chargers == nil {
						return nil, nil
					}
					return *l.Chargers, nil
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
				Type:        graphql.String,
				Description: "The URL for the logo of the operators of the destination charger.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					if l.DestinationChargerLogo == nil {
						return nil, nil
					}
					return *l.DestinationChargerLogo, nil
				},
			},
			"destinationWebsite": &graphql.Field{
				Type:        graphql.String,
				Description: "The URL for the operators of the destination charger.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					if l.DestinationWebsite == nil {
						return nil, nil
					}
					return *l.DestinationWebsite, nil
				},
			},
			"directionsLink": &graphql.Field{
				Type:        graphql.String,
				Description: "Pre-generated link to the location using Google Maps. Recommended to build your own using provided latitude and longitude fields.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					if l.DirectionsLink == nil {
						return nil, nil
					}
					return *l.DirectionsLink, nil
				},
			},
			"emails": &graphql.Field{
				Type:        graphql.NewList(emailType),
				Description: "The list of e-mail contacts for the location.",
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
				Type:        graphql.String,
				Description: "The HTML representation of hours of operation for the location.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					if l.Hours == nil {
						return nil, nil
					}
					return *l.Hours, nil
				},
			},
			"isGallery": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Whether or not the location is a gallery.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return bool(l.IsGallery), nil
				},
			},
			"kioskPinX": &graphql.Field{
				Type:        graphql.Int,
				Description: "Unknown what this information serves for.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return *l.KioskPinX, nil
				},
			},
			"kioskPinY": &graphql.Field{
				Type:        graphql.Int,
				Description: "Unknown what this information serves for.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return *l.KioskPinY, nil
				},
			},
			"kioskZoomPinX": &graphql.Field{
				Type:        graphql.Int,
				Description: "Unknown what this information serves for.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return *l.KioskZoomPinX, nil
				},
			},
			"kioskZoomPinY": &graphql.Field{
				Type:        graphql.Int,
				Description: "Unknown what this information serves for.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return *l.KioskZoomPinY, nil
				},
			},
			"latitude": &graphql.Field{
				Type: graphql.Float,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.Geo.Lat, nil
				},
			},
			"longitude": &graphql.Field{
				Type: graphql.Float,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.Geo.Lng, nil
				},
			},
			"locationId": &graphql.Field{
				Type:        graphql.String,
				Description: "The URL friendly title which is used in the path field.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.LocationID, nil
				},
			},
			"locationType": &graphql.Field{
				Type:        graphql.NewList(graphql.String),
				Description: "Each location may provide of 1 or many services such as supercharging, standard charging, destination charging, service, or a store.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.LocationType, nil
				},
			},
			"nid": &graphql.Field{
				Type:        graphql.Int,
				Description: "Internal Tesla specific unique identifer for the location.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.Nid, nil
				},
			},
			"openSoon": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Whether or not the location is opening soon.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return bool(l.OpenSoon), nil
				},
			},
			"path": &graphql.Field{
				Type:        graphql.String,
				Description: "The URL path to the location on Tesla's website.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return l.Path, nil
				},
			},
			"postalCode": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					if l.PostalCode == nil {
						return nil, nil
					}
					return *l.PostalCode, nil
				},
			},
			"provinceState": &graphql.Field{
				Type:        graphql.String,
				Description: "The ISO ALPHA-2 code of the province.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					if l.ProvinceState == nil {
						return nil, nil
					}
					return *l.ProvinceState, nil
				},
			},
			"region": &graphql.Field{
				Type:        graphql.String,
				Description: "The geographical region where the location resides in.",
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
				Type: graphql.Boolean,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return bool(l.SalesRepresentative), nil
				},
			},
			"subRegion": &graphql.Field{
				Type:        graphql.String,
				Description: "The State for locations in North America, otherwise country for a specific region.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					if l.SubRegion == nil {
						return nil, nil
					}
					return *l.SubRegion, nil
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

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"locations": &graphql.Field{
				Type: graphql.NewList(locationType),
				Args: locationFieldArguments,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					scope := database.NewGraphQLScopeWithFilters(p.Args)
					scope.Limit = -1

					locations, err := location.Locations(scope)
					if err != nil {
						return nil, err
					}

					return locations, nil
				},
			},
			"near": &graphql.Field{
				Type: graphql.NewList(locationType),
				Args: relay.NewConnectionArgs(graphql.FieldConfigArgument{
					"latitude": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.Float),
						Description: "The latitude of the coordinate.",
					},
					"longitude": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.Float),
						Description: "The longitude of the coordinate.",
					},
				}),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					scope := database.NewGraphQLScopeWithFilters(p.Args)
					scope.Limit = -1

					locations, err := location.LocationsNear(scope)
					if err != nil {
						return nil, err
					}

					return locations, nil
				},
			},
			"node": nodeDefinitions.NodeField,
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
}
