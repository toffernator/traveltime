package main

import (
	"context"
	"fmt"
	"os"

	routing "cloud.google.com/go/maps/routing/apiv2"
	"cloud.google.com/go/maps/routing/apiv2/routingpb"
	"google.golang.org/grpc/metadata"
)

func main() {
	ctx := context.Background()
	client, err := routing.NewRoutesClient(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	computeRoutesReq := routingpb.ComputeRoutesRequest{
		Origin: &routingpb.Waypoint{
			LocationType: &routingpb.Waypoint_PlaceId{
				PlaceId: "ChIJeRpOeF67j4AR9ydy_PIzPuM",
			},
		},
		Destination: &routingpb.Waypoint{
			LocationType: &routingpb.Waypoint_PlaceId{
				PlaceId: "ChIJG3kh4hq6j4AR_XuFQnV0_t8",
			},
		},
		RoutingPreference: routingpb.RoutingPreference_TRAFFIC_AWARE,
		TravelMode:        routingpb.RouteTravelMode_DRIVE,
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "X-Goog-FieldMask", "*")
	computeRoutesResponse, err := client.ComputeRoutes(ctx, &computeRoutesReq)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "response: %v\n", computeRoutesResponse)
}

