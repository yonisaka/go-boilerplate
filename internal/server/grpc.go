package server

type grpcServer struct {
}

// NewGRPCServer is a constructor of grpcServer.
func NewGRPCServer() Server {
	return &grpcServer{}
}

func (s *grpcServer) Run() error {
	//TODO implement me
	panic("implement me")
}

func (s *grpcServer) GracefulStop() {
	//TODO implement me
	panic("implement me")
}
