package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"github.com/DEMAxx/demin/hw12_13_14_15_calendar/events/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"log"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
)

type Server struct {
	httpServer      *http.Server
	grpcServer      *grpc.Server
	grpcHostAndPort string
	logger          Logger
	app             Application
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type LoggerConf struct {
	Level string
}

type Application interface {
	CreateEvent(ctx context.Context, id uuid.UUID, title string, date time.Time, description string, user uuid.UUID) error
}

type EventServiceServer struct {
	pb.UnimplementedEventsServer
	Server *Server
}

type UserServiceServer struct {
	pb.UnimplementedUsersServer
	Server *Server
}

func (s *UserServiceServer) UserCreate(ctx context.Context, req *pb.UserCreateRequest) (*pb.UserCreateResponse, error) {

	return nil, nil
}

func (s *EventServiceServer) EventCreate(ctx context.Context, req *pb.EventCreateRequest) (*pb.EventCreateResponse, error) {
	event := req.GetEvent()

	if event == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "event is required")
	}

	// Parse the date
	date, err := time.Parse("2006-01-02", event.GetDate())
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	if date.Before(time.Now()) {
		return nil, fmt.Errorf("date must be in the future")
	}

	if event.GetDuration() == nil {
		return nil, fmt.Errorf("duration is required")
	}

	eventID := uuid.New()

	userID := uuid.MustParse(event.GetUser())

	// Call the application layer to create the event
	err = s.Server.app.CreateEvent(ctx, eventID, event.GetTitle(), date, event.GetDescription(), userID)

	if err != nil {
		return nil, err
	}

	return &pb.EventCreateResponse{
		Id: eventID.String(),
	}, nil
}

func (s *EventServiceServer) EventUpdate(ctx context.Context, req *pb.EventUpdateRequest) (*pb.EventUpdateResponse, error) {
	// Implement the EventUpdate logic here
	return nil, nil
}

func (s *EventServiceServer) EventRemove(ctx context.Context, req *pb.EventRemoveRequest) (*pb.EventRemoveResponse, error) {
	// Implement the EventRemove logic here
	return nil, nil
}

func (s *EventServiceServer) EventWeekList(ctx context.Context, req *pb.EventWeekListRequest) (*pb.EventWeekListResponse, error) {
	return nil, nil
}

func (s *EventServiceServer) EventMonthList(ctx context.Context, req *pb.EventMonthListRequest) (*pb.EventMonthListResponse, error) {
	return nil, nil
}

func NewServer(logger Logger, hostAndPort string, grpcHostAndPort string, app Application) *Server {
	mux := http.NewServeMux()

	mux.Handle("/hello", LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		dateTime := time.Now().Format(time.RFC3339)
		method := r.Method
		path := r.URL.Path
		httpVersion := r.Proto
		userAgent := r.Header.Get("User-Agent")

		logger.Info(
			fmt.Sprintf(
				"Client IP: %s, DateTime: %s, Method: %s, Path: %s, HTTP Version: %s, User Agent: %s",
				clientIP, dateTime, method, path, httpVersion, userAgent,
			),
		)

		write, err := w.Write([]byte("Hello, World!"))
		if err != nil {
			return
		}
		logger.Info(fmt.Sprintf("response: %d", write))
	}), logger))

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(ValidationInterceptor),
	)

	return &Server{
		httpServer: &http.Server{
			Addr:              hostAndPort,
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
		},
		grpcServer:      grpcServer,
		grpcHostAndPort: grpcHostAndPort,
		logger:          logger,
		app:             app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("Starting HTTP and gRPC servers...")

	// Start HTTP server
	go func() {
		s.logger.Info("HTTP server start...")

		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("HTTP server ListenAndServe: " + err.Error())
		}

		s.logger.Info("HTTP server started")
	}()

	// Start gRPC server
	go func() {
		s.logger.Info("gRPC server start...")

		lsn, err := net.Listen("tcp", s.grpcHostAndPort)

		if err != nil {
			s.logger.Error(fmt.Sprintf("Failed to start gRPC server: %s", err.Error()))
			return
		}

		s.logger.Info(fmt.Sprintf("starting server on %s", lsn.Addr().String()))

		server := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				ValidationInterceptor,
			),
		)

		pb.RegisterEventsServer(server, new(EventServiceServer))

		if err := server.Serve(lsn); err != nil {
			log.Fatal(err)
		}

		s.logger.Info("gRPC server started")
	}()

	<-ctx.Done()
	return s.Stop(ctx)
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping HTTP and gRPC servers...")

	// Stop HTTP server
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		s.logger.Error("HTTP server Shutdown: " + err.Error())
		return err
	}

	// Stop gRPC server
	s.grpcServer.GracefulStop()

	s.logger.Info("HTTP and gRPC servers stopped")
	return nil
}
