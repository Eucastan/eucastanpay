package gateway

import (
	"context"

	adminpb "github.com/Eucastan/eucastanpay/common/proto/admin"
	adminReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/admin"
	adminResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/admin"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/mapper"
)

type AdminGateway struct {
	client adminpb.AdminServiceClient
}

func NewAdminGateway(client adminpb.AdminServiceClient) *AdminGateway {
	return &AdminGateway{
		client: client,
	}
}

func (s *AdminGateway) Register(ctx context.Context, req *adminReq.CreateAdminRequest) (*adminResp.AdminResponse, error) {

	grpcResp, err := s.client.Register(
		ctx,
		mapper.ToProtoCreateAdmin(req),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToCreateAdminResponse(grpcResp)

	return &resp, nil
}

func (s *AdminGateway) Login(ctx context.Context, req *adminReq.AdminLoginRequest) (*adminResp.AdminLoginResponse, error) {

	grpcResp, err := s.client.Login(
		ctx,
		mapper.ToProtoAdminLogin(req),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToLoginAdminResponse(grpcResp)

	return &resp, nil
}

func (s *AdminGateway) GetAllAdmins(ctx context.Context, limit, page int) (*adminResp.ListAdminsResponse, error) {

	grpcResp, err := s.client.GetAllAdmins(
		ctx,
		mapper.ToProtoListAdmins(limit, page),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToListAdminsResponse(grpcResp)

	return resp, nil
}

func (s *AdminGateway) GetAdmin(ctx context.Context, adminID string) (*adminResp.AdminResponse, error) {

	grpcResp, err := s.client.GetAdminByID(
		ctx,
		mapper.ToProtoGetAdmin(adminID),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToGetAdminResponse(grpcResp)

	return resp, nil
}

func (s *AdminGateway) UpdateAdmin(ctx context.Context, adminID string, req *adminReq.UpdateAdminRequest) (*adminResp.MessageResponse, error) {

	grpcResp, err := s.client.Update(
		ctx,
		mapper.ToProtoUpdateAdmin(adminID, req),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToUpdateAdminResponse(grpcResp)

	return resp, nil
}

func (s *AdminGateway) DeleteAdmin(ctx context.Context, adminID string) (*adminResp.MessageResponse, error) {

	grpcResp, err := s.client.Delete(
		ctx,
		mapper.ToProtoDeleteAdmin(adminID),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToDeleteAdminResponse(grpcResp)

	return resp, nil
}
