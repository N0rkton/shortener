package grpcfunc

import (
	"context"
	"fmt"
	"github.com/N0rkton/shortener/internal/app/config"
	"github.com/N0rkton/shortener/internal/app/handlers"
	"github.com/N0rkton/shortener/internal/app/storage"
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

// IndexPage - shortens originalURL
func (s *ShortenerServer) IndexPage(ctx context.Context, in *pb.IndexPageRequest) (*pb.IndexPageResponse, error) {
	var response pb.IndexPageResponse
	var userId string
	md, _ := metadata.FromIncomingContext(ctx)
	userId, ok := GetUserId(ctx)
	if !ok {
		md.Append("UserId", userId)
		err := grpc.SetHeader(ctx, md)
		if err != nil {
			return nil, err
		}
	}
	if !handlers.IsValidURL(in.OriginalURL) {
		response.Error = "Invalid URL"
		return &response, nil
	}
	code := handlers.GenerateRandomString(urlLen)
	err := store.AddURL(userId, code, in.OriginalURL)
	if err != nil {
		response.Error = fmt.Sprintf("Err %s occured adding url", err)
		return &response, nil
	}
	response.ShortURL = code
	return &response, nil
}

// RedirectTo - returns original url by its shorted url if exists
func (s *ShortenerServer) RedirectTo(ctx context.Context, in *pb.RedirectToRequest) (*pb.RedirectToResponse, error) {
	var response pb.RedirectToResponse
	originalURL, err := store.GetURL(in.ShortURL)
	if err != nil {
		response.Error = fmt.Sprintf("Err %s occured getting url", err)
		return &response, nil
	}
	response.OriginalURL = originalURL
	return &response, nil
}

// ListURLs - returns all shorted urls by user
func (s *ShortenerServer) ListURLs(ctx context.Context, in *pb.ListURLsRequest) (*pb.ListURLsResponse, error) {
	var response pb.ListURLsResponse
	var userId string
	md, _ := metadata.FromIncomingContext(ctx)
	userId, ok := GetUserId(ctx)
	if !ok {
		md.Append("UserId", userId)
		err := grpc.SetHeader(ctx, md)
		if err != nil {
			return nil, err
		}
	}
	resp, err := store.GetURLByID(userId)
	if err != nil {
		response.Error = fmt.Sprintf("Err %s occured getting url", err)
		return &response, nil
	}
	var originalShorts []*pb.OriginalShort

	for k, v := range resp {
		originalShorts = append(originalShorts, &pb.OriginalShort{
			OriginalURL: v,
			ShortURL:    config.GetBaseURL() + "/" + k,
		})
	}
	response.Urls = originalShorts
	return &response, nil
}

// DeleteUrl - delete urls if it was shorted by user
func (s *ShortenerServer) DeleteUrl(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	var response pb.DeleteResponse
	var userId string
	md, _ := metadata.FromIncomingContext(ctx)
	userId, ok := GetUserId(ctx)
	if !ok {
		md.Append("UserId", userId)
		err := grpc.SetHeader(ctx, md)
		if err != nil {
			return nil, err
		}
	}
	urls := in.URLStoDelete
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
		if !handlers.IsValidURL(req[k].OriginalURL) {
			response.Error = "Invalid URL"
			return &response, nil
		}
		code := handlers.GenerateRandomString(urlLen)
		originalShorts = append(originalShorts, &pb.OriginalShort{
			OriginalURL: req[k].CorrelationId,
			ShortURL:    config.GetBaseURL() + "/" + code})
		ok := store.AddURL(req[k].CorrelationId, code, req[k].OriginalURL)
		if ok != nil {
			response.Error = fmt.Sprintf("Err %s occured adding url", ok)
			return &response, nil
		}
	}
	response.Resp = originalShorts
	return &response, nil
}

// Stats - returns amount of shorted URLS and users
func (s *ShortenerServer) Stats(ctx context.Context, in *pb.StatsRequest) (*pb.StatsResponse, error) {
	var response pb.StatsResponse
	var err error
	subnet := config.GetTrustedSubnet()
	if subnet == "" {
		response.Error = "flag is not provided"
		return &response, nil
	}
	_, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		response.Error = fmt.Sprintf("json err %s", err)
		return &response, nil
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
			response.Error = "untrusted IP"
			return &response, nil
		}
	}
	domainRange := strings.Split(ipSplit[len(ipSplit)-1], "/")
	if realIPSplit[len(realIPSplit)-1] <= domainRange[0] || domainRange[1] <= realIPSplit[len(realIPSplit)-1] {
		response.Error = "untrusted IP"
		return &response, nil
	}
	response.URLs, response.Users, err = store.GetStats()
	if err != nil {
		response.Error = fmt.Sprintf("Err %s occured collecting stats", err)
		return &response, nil
	}
	return &response, nil
}
