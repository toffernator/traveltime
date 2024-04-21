package googleroutes

import (
	"context"

	routing "cloud.google.com/go/maps/routing/apiv2"
	"cloud.google.com/go/maps/routing/apiv2/routingpb"
	"google.golang.org/grpc/metadata"
)

type RouteComputationOpts struct {
	RoutingPreference  routingpb.TransitPreferences_TransitRoutingPreference
	AllowedTravelModes []routingpb.TransitPreferences_TransitTravelMode
}

var DefaultRouteComputationOpts RouteComputationOpts = RouteComputationOpts{
	RoutingPreference: routingpb.TransitPreferences_FEWER_TRANSFERS,
	AllowedTravelModes: []routingpb.TransitPreferences_TransitTravelMode{
		routingpb.TransitPreferences_BUS,
		routingpb.TransitPreferences_LIGHT_RAIL,
		routingpb.TransitPreferences_RAIL,
	},
}

func ComputeRoute(ctx context.Context, origin *routingpb.Waypoint, destination *routingpb.Waypoint, opts RouteComputationOpts) (*routingpb.ComputeRoutesResponse, error) {
	routesClient, err := routing.NewRoutesClient(ctx)
	if err != nil {
		return nil, err
	}
	defer routesClient.Close()

	req := &routingpb.ComputeRoutesRequest{
		Origin:      origin,
		Destination: destination,
		TravelMode:  routingpb.RouteTravelMode_TRANSIT,
		TransitPreferences: &routingpb.TransitPreferences{
			RoutingPreference:  opts.RoutingPreference,
			AllowedTravelModes: opts.AllowedTravelModes,
		},
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "X-Goog-FieldMask", "*")
	resp, err := routesClient.ComputeRoutes(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
