package clients

//type GrpcClient struct{}
//
//func NewGrpcClient() *GrpcClient {
//	conn, err := grpc.NewClient("delivery-service:50052",
//		grpc.WithTransportCredentials(insecure.NewCredentials()))
//	if err != nil {
//		log.Fatalf("new client: %v", err)
//		return err
//	}
//	defer conn.Close()
//
//	client := pb.NewPostsServiceClient(conn)
//}
