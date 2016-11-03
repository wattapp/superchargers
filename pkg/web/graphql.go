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
		Type: graphql.NewList(enumLocationType),
	},
	"region": &graphql.ArgumentConfig{
		Type: graphql.NewList(enumRegion),
	},
	"country": &graphql.ArgumentConfig{
		Type: graphql.NewList(enumCountry),
	},
	"openSoon": &graphql.ArgumentConfig{
		Type: graphql.Boolean,
	},
	"boundingBox": &graphql.ArgumentConfig{
		Type: graphql.NewList(graphql.Float),
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
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					if l.AddressNotes == nil {
						return nil, nil
					}
					return *l.AddressNotes, nil
				},
			},
			"amentities": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					if l.Amenities == nil {
						return nil, nil
					}
					return *l.Amenities, nil
				},
			},
			"baiduLat": &graphql.Field{
				Type: graphql.Float,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					if l.BaiduGeo == nil {
						return nil, nil
					}
					return l.BaiduGeo.Lat, nil
				},
			},
			"baiduLng": &graphql.Field{
				Type: graphql.Float,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					if l.BaiduGeo == nil {
						return nil, nil
					}
					return l.BaiduGeo.Lng, nil
				},
			},
			"chargers": &graphql.Field{
				Type: graphql.String,
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
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					if l.DestinationChargerLogo == nil {
						return nil, nil
					}
					return *l.DestinationChargerLogo, nil
				},
			},
			"destinationWebsite": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					if l.DestinationWebsite == nil {
						return nil, nil
					}
					return *l.DestinationWebsite, nil
				},
			},
			"directionsLink": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					if l.DirectionsLink == nil {
						return nil, nil
					}
					return *l.DirectionsLink, nil
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
					if l.Hours == nil {
						return nil, nil
					}
					return *l.Hours, nil
				},
			},
			"isGallery": &graphql.Field{
				Type: graphql.Boolean,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return bool(l.IsGallery), nil
				},
			},
			"kioskPinX": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return *l.KioskPinX, nil
				},
			},
			"kioskPinY": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return *l.KioskPinY, nil
				},
			},
			"kioskZoomPinX": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return *l.KioskZoomPinX, nil
				},
			},
			"kioskZoomPinY": &graphql.Field{
				Type: graphql.Int,
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
					return bool(l.OpenSoon), nil
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
					if l.PostalCode == nil {
						return nil, nil
					}
					return *l.PostalCode, nil
				},
			},
			"provinceState": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					if l.ProvinceState == nil {
						return nil, nil
					}
					return *l.ProvinceState, nil
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
				Type: graphql.Boolean,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					l := p.Source.(*location.Location)
					return bool(l.SalesRepresentative), nil
				},
			},
			"subRegion": &graphql.Field{
				Type: graphql.String,
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

	locationsConnectionDefinition := relay.ConnectionDefinitions(relay.ConnectionConfig{
		Name:     "Location",
		NodeType: locationType,
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"locations": &graphql.Field{
				Type: locationsConnectionDefinition.ConnectionType,
				Args: locationFieldArguments,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					scope := database.NewGraphQLScopeWithFilters(p.Args)

					// Include all results by default
					if scope.First == -1 && scope.Last == -1 {
						scope.Limit = -1
					}

					data := []database.GraphQLCursor{}
					locations, err := location.Locations(scope)
					if err != nil {
						return nil, err
					}

					for _, l := range locations {
						data = append(data, l)
					}

					return database.GraphQLConnection(data, scope), nil
				},
			},
			"node": nodeDefinitions.NodeField,
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
}
