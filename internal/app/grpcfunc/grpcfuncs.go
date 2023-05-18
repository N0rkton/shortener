package grpcfunc

import (
	"context"
	"github.com/N0rkton/shortener/internal/app/config"
	"github.com/N0rkton/shortener/internal/app/handlers"
	"github.com/N0rkton/shortener/internal/app/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"strings"

	pb "github.com/N0rkton/shortener/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const urlLen = 5

var store storage.Storage

// Init - selects type of storage from config
func Init() {
	var err error
	store, err = config.GetStorage()
	if err != nil {
		store = storage.NewMemoryStorage()
	}
}

// ShortenerServer поддерживает все необходимые методы сервера.
type ShortenerServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedShortenerServer
}

// GetUserId - search UserId key in metadata, if none generates new one
func GetUserId(ctx context.Context) (string, bool) {
	var userId string
	var value []string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		value = md.Get("UserId")
		if len(value) > 0 {
			// ключ содержит слайс строк, получаем первую строку
			userId = value[0]
			return userId, true
		}
	} else {
		userId = handlers.GenerateRandomString(3)
	}
	return userId, false
}

// UserIDInterceptor - unary Interceptor which checks for UserId key in metadata, and if it doesn't exist adds new one
func UserIDInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	token, ok := GetUserId(ctx)
	if ok {
		return handler, nil
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		md.Append("UserId", token)
		err := grpc.SetHeader(ctx, md)
		if err != nil {
			return nil, status.Error(codes.Internal, "SetHeader err")
		}
		return handler, nil
	}
	md2 := metadata.New(map[string]string{"UserId": token})
	metadata.NewIncomingContext(context.Background(), md2)
	err := grpc.SetHeader(ctx, md2)
	if err != nil {
		return nil, status.Error(codes.Internal, "SetHeader err")
	}
	h, err := handler(ctx, req)
	return h, err
}

// IndexPage - shortens originalURL
func (s *ShortenerServer) IndexPage(ctx context.Context, in *pb.IndexPageRequest) (*pb.IndexPageResponse, error) {
	var response pb.IndexPageResponse

	userId, ok := GetUserId(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "metadata err")
	}
	if !handlers.IsValidURL(in.OriginalUrl) {
		err := status.Error(codes.InvalidArgument, "Invalid URL")
		return nil, err
	}
	code := handlers.GenerateRandomString(urlLen)
	err := store.AddURL(userId, code, in.OriginalUrl)
	if err != nil {
		err = status.Error(codes.Internal, "Err occurred adding url")
		return nil, err
	}
	response.ShortUrl = code
	return &response, nil
}

// RedirectTo - returns original url by its shorted url if exists
func (s *ShortenerServer) RedirectTo(ctx context.Context, in *pb.RedirectToRequest) (*pb.RedirectToResponse, error) {
	var response pb.RedirectToResponse
	originalURL, err := store.GetURL(in.ShortURL)
	if err != nil {
		return nil, status.Error(codes.Internal, "Err occurred getting url")
	}
	response.OriginalURL = originalURL
	return &response, nil
}

// ListURLs - returns all shorted urls by user
func (s *ShortenerServer) ListURLs(ctx context.Context, in *pb.ListURLsRequest) (*pb.ListURLsResponse, error) {
	var response pb.ListURLsResponse
	userId, ok := GetUserId(ctx)
	if !ok {

		return nil, status.Error(codes.Internal, "metadata err")

	}
	resp, err := store.GetURLByID(userId)
	if err != nil {
		err = status.Error(codes.Internal, "Err occurred getting url")
		return nil, err
	}
	var originalShorts []*pb.OriginalShort

	for k, v := range resp {
		originalShorts = append(originalShorts, &pb.OriginalShort{
			OriginalUrl: v,
			ShortUrl:    config.GetBaseURL() + "/" + k,
		})
	}
	response.Urls = originalShorts
	return &response, nil
}

// DeleteUrl - delete urls if it was shorted by user
func (s *ShortenerServer) DeleteUrl(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	var response pb.DeleteResponse

	userId, ok := GetUserId(ctx)
	if !ok {

		return nil, status.Error(codes.Internal, "metadata err")

	}
	urls := in.UrlsToDelete
	for i := 0; i < len(urls); i++ {
		job := handlers.DeleteURLJob{URL: urls[i], UserID: userId}
		handlers.JobCh <- job
	}
	return &response, nil
}

// Batch - shortens batch of urls
func (s *ShortenerServer) Batch(ctx context.Context, in *pb.BatchRequest) (*pb.BatchResponse, error) {
	var response pb.BatchResponse
	var originalShorts []*pb.OriginalShort

	req := in.Req
	for k := range req {
		if !handlers.IsValidURL(req[k].OriginalUrl) {
			err := status.Error(codes.InvalidArgument, "Invalid URL")
			return nil, err
		}
		code := handlers.GenerateRandomString(urlLen)
		originalShorts = append(originalShorts, &pb.OriginalShort{
			OriginalUrl: req[k].CorrelationId,
			ShortUrl:    config.GetBaseURL() + "/" + code})
		ok := store.AddURL(req[k].CorrelationId, code, req[k].OriginalUrl)
		if ok != nil {
			err := status.Error(codes.Internal, "Err occurred adding url")
			return nil, err
		}
	}
	response.Resp = originalShorts
	return &response, nil
}

// Stats - returns amount of shorted URLS and users
func (s *ShortenerServer) Stats(ctx context.Context, in *emptypb.Empty) (*pb.StatsResponse, error) {
	var response pb.StatsResponse
	var err error
	subnet := config.GetTrustedSubnet()
	if subnet == "" {
		err = status.Error(codes.NotFound, "flag is not provided")
		return nil, err
	}
	_, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		err = status.Error(codes.Internal, "ParseCIDR error")
		return nil, err
	}
	var realIPStr string
	var value []string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		value = md.Get("X-Real-IP")
		if len(value) > 0 {
			realIPStr = value[0]
		}
	}

	ip := strings.Replace(ipNet.String(), ":", ".", -1)
	realIP := strings.Replace(realIPStr, ":", ".", -1)
	ipSplit := strings.Split(ip, ".")
	realIPSplit := strings.Split(realIP, ".")
	for i := 0; i < len(realIPSplit)-1; i++ {
		if ipSplit[i] != realIPSplit[i] {
			err = status.Error(codes.PermissionDenied, "untrusted IP")
			return nil, err
		}
	}
	domainRange := strings.Split(ipSplit[len(ipSplit)-1], "/")
	ipIndex := realIPSplit[len(realIPSplit)-1]
	if ipIndex <= domainRange[0] && domainRange[1] <= ipIndex {
		err = status.Error(codes.PermissionDenied, "untrusted IP")
		return nil, err
	}
	response.UrlCount, response.UserCount, err = store.GetStats()
	if err != nil {
		err = status.Error(codes.Internal, "Err occurred collecting stats")
		return nil, err
	}
	return &response, nil
}
