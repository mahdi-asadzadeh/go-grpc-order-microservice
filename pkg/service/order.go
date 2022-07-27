package service

import (
	"context"
	"github.com/mahdi-asadzadeh/go-grpc-order-microservice/pkg/db"
	"github.com/mahdi-asadzadeh/go-grpc-order-microservice/pkg/models"
	"github.com/mahdi-asadzadeh/go-grpc-order-microservice/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedOrderServiceServer
	H             db.Handler
	ProductClient pb.ProductServiceClient
}

func (s *Server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	productReq := pb.DetailProductRequest{Id: req.GetProductId()}
	productRes, err := s.ProductClient.DetailProduct(context.Background(), &productReq)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Product not found.")
	}
	newOrder := models.Order{
		Price:     float64(req.GetPrice() * float32(req.GetQuantity())),
		ProductID: uint(productRes.GetProduct().GetId()),
		UserID:    uint(req.GetUserId()),
		Quantity:  uint(req.GetQuantity()),
	}
	if err := s.H.DB.Create(&newOrder).Error; err != nil {
		return nil, status.Error(codes.Internal, "Order invalid data.")
	}

	reqOrder := pb.CreateOrderResponse{
		Order: &pb.Order{
			Id:        int64(newOrder.ID),
			Price:     float32(newOrder.Price),
			ProductId: productRes.GetProduct().GetId(),
			UserId:    req.GetUserId(),
			Quantity:  req.GetQuantity(),
		},
	}
	return &reqOrder, nil
}

func (s *Server) ListOrder(req *pb.ListOrderRequest, stream pb.OrderService_ListOrderServer) error {
	var orders []models.Order
	offSet := (req.GetPage() - 1) * req.GetPageSize()
	err := s.H.DB.Offset(int(offSet)).Limit(int(req.GetPageSize())).Find(&orders).Error
	if err != nil {
		return status.Error(codes.Internal, "Orders are empty.")
	}
	for _, or := range orders {
		res := pb.ListOrderResponse{
			Id:        int64(or.ID),
			Price:     float32(or.Price),
			ProductId: int64(or.ProductID),
			UserId:    int64(or.UserID),
			Quantity:  int32(or.Quantity),
			CreateAt:  or.CreatedAt.String(),
			UpdateAt:  or.UpdatedAt.String(),
		}
		stream.Send(&res)
	}
	return nil
}

func (s *Server) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	err := s.H.DB.Where("user_id = ? AND id = ?", uint(req.GetUserId()), uint(req.GetOrderId())).Delete(&models.Order{}).Error
	if err != nil {
		return nil, status.Error(codes.NotFound, "Not found order.")
	}
	res := pb.DeleteOrderResponse{}
	return &res, nil
}
