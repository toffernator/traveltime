package googleroutes

import (
	"context"
	"time"

	routing "cloud.google.com/go/maps/routing/apiv2"
	"cloud.google.com/go/maps/routing/apiv2/routingpb"
	"google.golang.org/grpc/metadata"
)

func CalculateTravelTime(ctx context.Context, origin *routingpb.Waypoint, destination *routingpb.Waypoint) (time.Duration, error) {
	routesClient, err := routing.NewRoutesClient(ctx)
	if err != nil {
		return time.Duration(0), err
	}
	defer routesClient.Close()

	req := &routingpb.ComputeRoutesRequest{
		Origin:      origin,
		Destination: destination,
		TravelMode:  routingpb.RouteTravelMode_TRANSIT,
		TransitPreferences: &routingpb.TransitPreferences{
			// TODO: Calculate and report both
			RoutingPreference: routingpb.TransitPreferences_FEWER_TRANSFERS,
			// TODO: Allowed travelmodes
		},
	}

	// ctx = metadata.AppendToOutgoingContext(ctx, "X-Goog-FieldMask", "routes.localizedValues")
	ctx = metadata.AppendToOutgoingContext(ctx, "X-Goog-FieldMask", "*")
	resp, err := routesClient.ComputeRoutes(ctx, req)
	if err != nil {
		return time.Duration(0), err
	}

	duration := resp.Routes[0].Duration.AsDuration()
	return duration, nil
}
