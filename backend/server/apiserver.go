package server

// APIServer defines the server structure for the API server.
type APIServer struct {
	*server
}

func NewAPIServer() (*APIServer, error) {
	base, err := newBaseServer()
	if err != nil {
		return nil, err
	}

	return &APIServer{base}, nil
}
