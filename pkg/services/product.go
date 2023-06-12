package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/uyabpras/go-grpc-product-svc/pkg/db"
	"github.com/uyabpras/go-grpc-product-svc/pkg/models"
	"github.com/uyabpras/go-grpc-product-svc/pkg/proto/pb"
)

type Server struct {
	H db.Handler
}

func (s *Server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	var product models.Product

	product.Name = req.Name
	product.Stock = req.Stock
	product.Price = req.Price

	if result := s.H.DB.Find(&product, req.Name); result != nil {
		if result := s.H.DB.Where(&models.Product{Name: req.Name}).First(&product); result.Error != nil {
			return &pb.CreateProductResponse{
				Status: http.StatusConflict,
				Error:  "Stock already updated",
				Id:     product.Id,
			}, nil
		}

		product.Stock = product.Stock + req.Stock
		s.H.DB.Save(&product)

		return &pb.CreateProductResponse{
			Status: http.StatusOK,
			Id:     product.Id,
		}, nil
	}

	if result := s.H.DB.Create(&product); result.Error != nil {
		return &pb.CreateProductResponse{
			Status: http.StatusConflict,
			Error:  result.Error.Error(),
		}, nil
	}

	return &pb.CreateProductResponse{
		Status: http.StatusCreated,
		Id:     product.Id,
	}, nil

}

func (s *Server) FindOne(ctx context.Context, req *pb.FindOneRequest) (*pb.FindOneResponse, error) {
	var product models.Product

	if result := s.H.DB.Find(&product, req.Id); result.Error != nil {
		return &pb.FindOneResponse{
			Status: http.StatusNotFound,
			Error:  result.Error.Error(),
		}, nil
	}
	data := &pb.FindOneData{
		Id:    product.Id,
		Name:  product.Name,
		Stock: product.Stock,
		Price: product.Price,
	}

	return &pb.FindOneResponse{
		Status: http.StatusOK,
		Data:   data,
	}, nil
}

func (s *Server) DecreaseStock(ctx context.Context, req *pb.DecreaseStockRequest) (*pb.DecreaseStockResponse, error) {
	var product models.Product

	if result := s.H.DB.First(&product, req.Id); result.Error != nil {
		return &pb.DecreaseStockResponse{
			Status: http.StatusNotFound,
			Error:  result.Error.Error(),
		}, nil
	}

	if product.Stock < req.Quantity {
		return &pb.DecreaseStockResponse{
			Status: http.StatusConflict,
			Error:  "Stock too low",
		}, nil
	}

	var log models.Stock_decrease

	if result := s.H.DB.Where(&models.Stock_decrease{OrderID: req.OrderId}).First(&log); result.Error == nil {
		return &pb.DecreaseStockResponse{
			Status: http.StatusConflict,
			Error:  "Stock already decreased",
		}, nil
	}

	product.Stock = product.Stock - req.Quantity

	s.H.DB.Save(&product)

	log.OrderID = req.OrderId
	log.ProductRefer = product.Id
	log.Quantity = req.Quantity

	s.H.DB.Create(&log)

	return &pb.DecreaseStockResponse{
		Status: http.StatusOK,
	}, nil
}

func (s *Server) ListProduk(ctx context.Context, req *pb.ListproductsRequest) (*pb.ListProductResponse, error) {
	var products []models.Product
	var count int64

	fmt.Println(req.Page)

	if req.Page == 0 {
		req.Page = 1
	}

	if req.Limit == 0 {
		req.Limit = 5
	}

	offset := (req.Page - 1) * req.Limit

	fmt.Println(req.Page, req.Limit, offset)
	if result := s.H.DB.Offset(int(offset)).Limit(int(req.Limit)).Find(&products); result.Error != nil {
		return &pb.ListProductResponse{
			Status: http.StatusNotFound,
			Error:  result.Error.Error(),
		}, nil
	}

	if result := s.H.DB.Model(&products).Count(&count); result.Error != nil {
		return &pb.ListProductResponse{
			Status: http.StatusConflict,
			Error:  result.Error.Error(),
		}, nil
	}
	totalPages := count / int64(req.Limit)

	var data []*pb.Product

	for _, product := range products {
		productMessage := &pb.Product{
			Id:    product.Id,
			Name:  product.Name,
			Stock: product.Stock,
			Price: product.Price,
		}
		data = append(data, productMessage)
	}

	response := &pb.ListProductResponse{
		Data:       data,
		TotalPages: totalPages,
	}
	return response, nil
}
