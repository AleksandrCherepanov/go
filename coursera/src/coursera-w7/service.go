package main

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type ServerStruct struct {
	acl          map[string][]string
	logServices  map[chan Event]Admin_LoggingServer
	statServices map[chan Stat]*StatInterval
	stat         map[chan Stat]*Stat
	muLog        *sync.Mutex
	muStat       *sync.Mutex
	startTime    time.Time
}

type NullString struct {
	value  string
	isNull bool
}

type RequestData struct {
	consumer NullString
	host     string
	method   string
}

func (s *ServerStruct) addToACL(key string, value string) {
	_, ok := s.acl[key]
	if !ok {
		s.acl[key] = make([]string, 0)
	}

	s.acl[key] = append(s.acl[key], value)
}

func getRequestData(ctx context.Context) (*RequestData, error) {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("can't get context metadata")
	}

	method, ok := grpc.Method(ctx)
	if !ok {
		return nil, errors.New("can't get method of context")
	}

	consumer := metadata.Get("consumer")

	consumerData := NullString{}
	if len(consumer) == 0 {
		consumerData.value = ""
		consumerData.isNull = true
	} else {
		consumerData.value = consumer[0]
		consumerData.isNull = false
	}

	peer, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("can't get peer of context")
	}

	return &RequestData{
		consumer: consumerData,
		method:   method,
		host:     peer.Addr.String(),
	}, nil
}

func checkACL(ctx context.Context, acl map[string][]string) error {
	contextData, err := getRequestData(ctx)
	if err != nil {
		return err
	}

	if contextData.consumer.isNull {
		return status.Error(codes.Unauthenticated, "empty consumer")
	}

	aclData, ok := acl[contextData.consumer.value]
	if !ok {
		return status.Error(codes.Unauthenticated, "unknown consumer")
	}

	for _, acl := range aclData {
		ok, _ := regexp.Match(acl, []byte(contextData.method))
		if ok {
			return nil
		}
	}

	return status.Error(codes.Unauthenticated, "method is forbidden")
}

func logger(server *ServerStruct, requestData *RequestData) {
	event := new(Event)
	event.Consumer = requestData.consumer.value
	event.Method = requestData.method
	event.Host = requestData.host

	server.muLog.Lock()
	for logChan := range server.logServices {
		logChan <- *event
	}
	server.muLog.Unlock()
}

func (s *ServerStruct) Logging(req *Nothing, srv Admin_LoggingServer) error {
	logChannel := make(chan Event)
	s.muLog.Lock()
	s.logServices[logChannel] = srv
	s.muLog.Unlock()

	for {
		select {
		case event := <-logChannel:
			srv.Send(&event)
		}
	}
	return nil
}

func statTimer(server *ServerStruct, req *StatInterval, statChan chan Stat) {
	for {
		time.Sleep(time.Duration(req.IntervalSeconds * uint64(time.Second)))
		server.muStat.Lock()
		stat, ok := server.stat[statChan]
		if ok {
			for key := range server.statServices {
				if key == statChan {
					key <- *stat
				}
			}

			server.stat[statChan] = &Stat{
				ByMethod:   make(map[string]uint64, 0),
				ByConsumer: make(map[string]uint64, 0),
			}
		}

		server.muStat.Unlock()
	}
}

func statistic(server *ServerStruct, requestData *RequestData) {
	server.muStat.Lock()
	if len(server.statServices) <= 0 {
		server.muStat.Unlock()
		return
	}

	for _, stat := range server.stat {
		stat.ByMethod[requestData.method]++
		stat.ByConsumer[requestData.consumer.value]++
	}

	server.muStat.Unlock()
}

func (s *ServerStruct) checkAuth(ctx context.Context) error {
	return checkACL(ctx, s.acl)
}

func (s *ServerStruct) Statistics(req *StatInterval, srv Admin_StatisticsServer) error {
	statChannel := make(chan Stat)

	s.muStat.Lock()
	s.statServices[statChannel] = req
	s.stat[statChannel] = &Stat{
		ByMethod:   make(map[string]uint64),
		ByConsumer: make(map[string]uint64),
	}
	s.muStat.Unlock()

	go statTimer(s, req, statChannel)

	for {
		select {
		case stat := <-statChannel:
			srv.Send(&stat)
		}
	}
	return nil
}

func (s *ServerStruct) Check(ctx context.Context, req *Nothing) (*Nothing, error) {
	return req, nil
}

func (s *ServerStruct) Add(ctx context.Context, req *Nothing) (*Nothing, error) {
	return req, nil
}

func (s *ServerStruct) Test(ctx context.Context, req *Nothing) (*Nothing, error) {
	return req, nil
}

func stopServer(ctx context.Context, server *grpc.Server) {
	for {
		select {
		case <-ctx.Done():
			server.Stop()
			return
		}
	}
}

func checkAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	server, ok := info.Server.(*ServerStruct)
	if !ok {
		return nil, errors.New("Unexpected server handler")
	}

	err := server.checkAuth(ctx)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

func loggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	server, ok := info.Server.(*ServerStruct)
	if !ok {
		return nil, errors.New("Unexpected server handler")
	}

	requestData, err := getRequestData(ctx)
	if err != nil {
		return nil, err
	}

	logger(server, requestData)
	statistic(server, requestData)
	return handler(ctx, req)
}

func checkAuthStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	server, ok := srv.(*ServerStruct)
	if !ok {
		return errors.New("Unexpected server handler")
	}

	err := server.checkAuth(ss.Context())
	if err != nil {
		return err
	}

	return handler(srv, ss)
}

func loggerStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	server, ok := srv.(*ServerStruct)
	if !ok {
		return errors.New("Unexpected server handler")
	}

	requestData, err := getRequestData(ss.Context())
	if err != nil {
		return err
	}

	logger(server, requestData)
	statistic(server, requestData)
	return handler(srv, ss)
}

func StartMyMicroservice(ctx context.Context, listenAddr, ACLData string) error {
	aclData := make(map[string][]string, 0)
	err := json.Unmarshal([]byte(ACLData), &aclData)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	interceptors := make([]grpc.UnaryServerInterceptor, 0)
	interceptors = append(interceptors, checkAuthInterceptor)
	interceptors = append(interceptors, loggerInterceptor)

	streamInterceptors := make([]grpc.StreamServerInterceptor, 0)
	streamInterceptors = append(streamInterceptors, checkAuthStreamInterceptor)
	streamInterceptors = append(streamInterceptors, loggerStreamInterceptor)

	so := []grpc.ServerOption{grpc.ChainUnaryInterceptor(interceptors...), grpc.ChainStreamInterceptor(streamInterceptors...)}

	server := grpc.NewServer(so...)

	s := &ServerStruct{
		acl:          make(map[string][]string, 0),
		muLog:        &sync.Mutex{},
		muStat:       &sync.Mutex{},
		logServices:  make(map[chan Event]Admin_LoggingServer, 0),
		statServices: make(map[chan Stat]*StatInterval, 0),
		stat:         make(map[chan Stat]*Stat, 0),
		startTime:    time.Now(),
	}

	for k, v := range aclData {
		for _, aclValue := range v {
			route := strings.ReplaceAll(aclValue, "*", ".*")
			s.addToACL(k, route)
		}
	}

	RegisterAdminServer(server, s)
	RegisterBizServer(server, s)

	go server.Serve(lis)
	go stopServer(ctx, server)

	return nil
}
