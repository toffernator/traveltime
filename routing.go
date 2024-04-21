package main

import (
	"time"

	routing "cloud.google.com/go/maps/routing/apiv2"
	"cloud.google.com/go/maps/routing/apiv2/routingpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TravelMode int
type RoutingPreference int

const (
	Rail TravelMode = iota
	Bus

	LessWalking RoutingPreference = iota
	LessTransfers
)

type computeTransitTravelTimeOptions struct {
	allowedTravelModes *map[TravelMode]bool
	routingPreference  *RoutingPreference
	departure          *time.Time
}

type ComputeTransitTravelTimeOptions func(options *computeTransitTravelTimeOptions) error

func WithAllowedTravelMode(mode TravelMode) ComputeTransitTravelTimeOptions {
	return func(options *computeTransitTravelTimeOptions) error {
		if options.allowedTravelModes == nil {
			options.allowedTravelModes = &map[TravelMode]bool{}
		}

		dereferencedAllowedTravelModes := *options.allowedTravelModes
		dereferencedAllowedTravelModes[mode] = true
		return nil
	}
}

func WithRoutingPreference(preference RoutingPreference) ComputeTransitTravelTimeOptions {
	return func(options *computeTransitTravelTimeOptions) error {
		options.routingPreference = &preference
		return nil
	}
}

func WithArrivalTime(arrival time.Time) ComputeTransitTravelTimeOptions {
	return func(options *computeTransitTravelTimeOptions) error {
		options.departure = &arrival
		return nil
	}
}

type ComputeTravelTimeResult struct {
	Origin      Address
	Destination Address
	Duration    time.Duration
	Options     computeTransitTravelTimeOptions
}

type Address string

func ComputeTravelTime(ctx context.Context, origin Address, destination Address, opts ...ComputeTransitTravelTimeOptions) (ComputeTravelTimeResult, error) {
	var options computeTransitTravelTimeOptions
	for _, opt := range opts {
		if err := opt(&options); err != nil {
			return ComputeTravelTimeResult{}, nil
		}
	}

	if options.routingPreference == nil {
		lessTransfers := LessTransfers
		options.routingPreference = &lessTransfers
	}
	if options.allowedTravelModes == nil {
		options.allowedTravelModes = &map[TravelMode]bool{
			Rail: true,
			Bus:  true,
		}
	}
	if options.departure == nil {
		nearestTuesdayAt1000Utc := nearestTuesdayAt1000Utc()
		options.departure = &nearestTuesdayAt1000Utc
	}

	return computeTravelTime(ctx, origin, destination, options)
}

func computeTravelTime(ctx context.Context, origin Address, destination Address, opts computeTransitTravelTimeOptions) (ComputeTravelTimeResult, error) {
	originWaypoint := &routingpb.Waypoint{
		LocationType: &routingpb.Waypoint_Address{
			Address: string(origin),
		},
	}

	destinationWaypoint := &routingpb.Waypoint{
		LocationType: &routingpb.Waypoint_Address{
			Address: string(destination),
		},
	}

	routesClient, err := routing.NewRoutesClient(ctx)
	if err != nil {
		return ComputeTravelTimeResult{}, err
	}
	defer routesClient.Close()

	transitPreferences := &routingpb.TransitPreferences{
		RoutingPreference:  toRoutingpbRoutingPrefernce[*opts.routingPreference],
		AllowedTravelModes: make([]routingpb.TransitPreferences_TransitTravelMode, 0),
	}
	for travelMode, isIncluded := range *opts.allowedTravelModes {
		if !isIncluded {
			continue
		}
		asRouingpbTravelMode := toRoutingpbTravelMode[travelMode]
		transitPreferences.AllowedTravelModes = append(transitPreferences.AllowedTravelModes, asRouingpbTravelMode)
	}

	req := &routingpb.ComputeRoutesRequest{
		Origin:             originWaypoint,
		Destination:        destinationWaypoint,
		TravelMode:         routingpb.RouteTravelMode_TRANSIT,
		TransitPreferences: transitPreferences,
		ArrivalTime:        timestamppb.New(*opts.departure),
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "X-Goog-FieldMask", "routes.duration")
	resp, err := routesClient.ComputeRoutes(ctx, req)
	if err != nil {
		return ComputeTravelTimeResult{}, err
	}

	return ComputeTravelTimeResult{Origin: origin, Destination: destination, Duration: resp.Routes[0].Duration.AsDuration(), Options: opts}, nil
}

var toRoutingpbRoutingPrefernce map[RoutingPreference]routingpb.TransitPreferences_TransitRoutingPreference = map[RoutingPreference]routingpb.TransitPreferences_TransitRoutingPreference{
	LessWalking:   routingpb.TransitPreferences_LESS_WALKING,
	LessTransfers: routingpb.TransitPreferences_FEWER_TRANSFERS,
}

var toRoutingpbTravelMode map[TravelMode]routingpb.TransitPreferences_TransitTravelMode = map[TravelMode]routingpb.TransitPreferences_TransitTravelMode{
	Rail: routingpb.TransitPreferences_RAIL,
	Bus:  routingpb.TransitPreferences_BUS,
}

func nearestTuesdayAt1000Utc() time.Time {
	now := time.Now()
	todayAt1000Utc := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, time.UTC)

	switch todayAt1000Utc.Weekday() {
	case time.Wednesday:
		return todayAt1000Utc.AddDate(0, 0, 6)
	case time.Thursday:
		return todayAt1000Utc.AddDate(0, 0, 5)
	case time.Friday:
		return todayAt1000Utc.AddDate(0, 0, 4)
	case time.Saturday:
		return todayAt1000Utc.AddDate(0, 0, 3)
	case time.Sunday:
		return todayAt1000Utc.AddDate(0, 0, 2)
	case time.Monday:
		return todayAt1000Utc.AddDate(0, 0, 1)
	default:
		return todayAt1000Utc
	}
}
