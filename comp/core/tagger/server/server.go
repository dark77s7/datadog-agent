// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package server implements a gRPC server that streams Tagger entities.
package server

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"

	tagger "github.com/DataDog/datadog-agent/comp/core/tagger/def"
	"github.com/DataDog/datadog-agent/comp/core/tagger/proto"
	"github.com/DataDog/datadog-agent/comp/core/tagger/types"
	pb "github.com/DataDog/datadog-agent/pkg/proto/pbgo/core"
	"github.com/DataDog/datadog-agent/pkg/util/grpc"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

const (
	taggerStreamSendTimeout = 1 * time.Minute
	streamKeepAliveInterval = 9 * time.Minute
)

// Server is a grpc server that streams tagger entities
type Server struct {
	taggerComponent tagger.Component
	maxEventSize    int

	// if set to true, the server loads one chunk at a time into memory, sends it on the stream, and then loads another chunk
	// if set to false, the server loads all chunks to memory at once (slice of chunks), and then sends them sequentially over the stream
	useLazyEventChunking bool
}

// NewServer returns a new Server
func NewServer(t tagger.Component, maxEventSize int, useLazyEventChunking bool) *Server {
	return &Server{
		taggerComponent: t,
		maxEventSize:    maxEventSize,
	}
}

// TaggerStreamEntities subscribes to added, removed, or changed entities in the Tagger
// and streams them to clients as pb.StreamTagsResponse events. Filtering is as
// of yet not implemented.
func (s *Server) TaggerStreamEntities(in *pb.StreamTagsRequest, out pb.AgentSecure_TaggerStreamEntitiesServer) error {
	cardinality, err := proto.Pb2TaggerCardinality(in.GetCardinality())
	if err != nil {
		return err
	}

	filterBuilder := types.NewFilterBuilder()
	for _, prefix := range in.GetPrefixes() {
		filterBuilder = filterBuilder.Include(types.EntityIDPrefix(prefix))
	}

	filter := filterBuilder.Build(cardinality)

	streamingID := in.GetStreamingID()
	if streamingID == "" {
		// this is done to preserve backward compatibility
		// if CLC runner is using an old version, the streaming ID would be an empty string,
		// and the server needs to auto-assign a unique id
		streamingID = uuid.New().String()
	}

	subscriptionID := fmt.Sprintf("streaming-client-%s", streamingID)
	subscription, err := s.taggerComponent.Subscribe(subscriptionID, filter)
	if err != nil {
		return err
	}

	defer subscription.Unsubscribe()

	ticker := time.NewTicker(streamKeepAliveInterval)
	defer ticker.Stop()
	for {
		select {
		case events, ok := <-subscription.EventsChan():
			if !ok {
				log.Warnf("subscriber channel closed, client will reconnect")
				return fmt.Errorf("subscriber channel closed")
			}

			ticker.Reset(streamKeepAliveInterval)

			responseEvents := make([]*pb.StreamTagsEvent, 0, len(events))
			for _, event := range events {
				e, err := proto.Tagger2PbEntityEvent(event)
				if err != nil {
					log.Warnf("can't convert tagger entity to protobuf: %s", err)
					continue
				}

				responseEvents = append(responseEvents, e)
			}

			// Split events into chunks and send each one
			chunks := splitEventsLazy(responseEvents, s.maxEventSize)
			for chunk := range chunks {
				if len(chunk) == 0 {
					continue
				}
				err = grpc.DoWithTimeout(func() error {
					return out.Send(&pb.StreamTagsResponse{
						Events: chunk,
					})
				}, taggerStreamSendTimeout)

				if err != nil {
					log.Warnf("error sending tagger event: %s", err)
					s.taggerComponent.GetTaggerTelemetryStore().ServerStreamErrors.Inc()
					return err
				}
			}

		case <-out.Context().Done():
			return nil

		// The remote tagger client has a timeout that closes the
		// connection after 10 minutes of inactivity (implemented in
		// comp/core/tagger/remote/tagger.go) In order to avoid closing the
		// connection and having to open it again, the server will send
		// an empty message after 9 minutes of inactivity. The goal is
		// only to keep the connection alive without losing the
		// protection against “half” closed connections brought by the
		// timeout.
		case <-ticker.C:
			err = grpc.DoWithTimeout(func() error {
				return out.Send(&pb.StreamTagsResponse{
					Events: []*pb.StreamTagsEvent{},
				})
			}, taggerStreamSendTimeout)

			if err != nil {
				log.Warnf("error sending tagger keep-alive: %s", err)
				s.taggerComponent.GetTaggerTelemetryStore().ServerStreamErrors.Inc()
				return err
			}
		}
	}
}

// TaggerFetchEntity fetches an entity from the Tagger with the desired cardinality tags.
//
//nolint:revive // TODO(CINT) Fix revive linter
func (s *Server) TaggerFetchEntity(_ context.Context, in *pb.FetchEntityRequest) (*pb.FetchEntityResponse, error) {
	if in.Id == nil {
		return nil, status.Errorf(codes.InvalidArgument, `missing "id" parameter`)
	}

	entityID := types.NewEntityID(types.EntityIDPrefix(in.Id.Prefix), in.Id.Uid)
	cardinality, err := proto.Pb2TaggerCardinality(in.GetCardinality())
	if err != nil {
		return nil, err
	}

	tags, err := s.taggerComponent.Tag(entityID, cardinality)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	return &pb.FetchEntityResponse{
		Id:          in.Id,
		Cardinality: in.GetCardinality(),
		Tags:        tags,
	}, nil
}
