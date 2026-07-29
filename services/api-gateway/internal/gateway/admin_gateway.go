package gateway

import (
	"context"

	adminpb "github.com/Eucastan/eucastanpay/common/proto/admin"
	adminReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/admin"
	adminResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/admin"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/mapper"
	"github.com/sirupsen/logrus"
)

type AdminGateway struct {
	client adminpb.AdminServiceClient
	logger *logrus.Logger
}

func NewAdminGateway(client adminpb.AdminServiceClient, logger *logrus.Logger) *AdminGateway {
	return &AdminGateway{
		client: client,
		logger: logger,
	}
}

func (g *AdminGateway) CreateBootstrapAdmin(ctx context.Context, req *adminReq.CreateAdminRequest) (*adminResp.AdminResponse, error) {

	grpcResp, err := g.client.BootstrapAdmin(
		ctx,
		mapper.ToProtoCreateAdmin(req),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	resp := mapper.ToCreateAdminResponse(grpcResp)
	return &resp, nil
}

func (g *AdminGateway) Register(ctx context.Context, req *adminReq.CreateAdminRequest) (*adminResp.AdminResponse, error) {

	grpcResp, err := g.client.Register(
		ctx,
		mapper.ToProtoCreateAdmin(req),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	resp := mapper.ToCreateAdminResponse(grpcResp)
	return &resp, nil
}

func (g *AdminGateway) Login(ctx context.Context, req *adminReq.AdminLoginRequest) (*adminResp.AdminLoginResponse, error) {

	grpcResp, err := g.client.Login(
		ctx,
		mapper.ToProtoAdminLogin(req),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	resp := mapper.ToLoginAdminResponse(grpcResp)
	return &resp, nil
}

func (g *AdminGateway) GetAllAdmins(ctx context.Context, limit, page int) (*adminResp.ListAdminsResponse, error) {

	grpcResp, err := g.client.GetAllAdmins(
		ctx,
		mapper.ToProtoListAdmins(limit, page),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	resp := mapper.ToListAdminsResponse(grpcResp)
	return resp, nil
}

func (g *AdminGateway) GetAdmin(ctx context.Context, adminID string) (*adminResp.AdminResponse, error) {

	grpcResp, err := g.client.GetAdminByID(
		ctx,
		mapper.ToProtoGetAdmin(adminID),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	resp := mapper.ToGetAdminResponse(grpcResp)
	return resp, nil
}

func (g *AdminGateway) UpdateAdmin(ctx context.Context, adminID string, req *adminReq.UpdateAdminRequest) (*adminResp.MessageResponse, error) {

	grpcResp, err := g.client.Update(
		ctx,
		mapper.ToProtoUpdateAdmin(adminID, req),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	resp := mapper.ToUpdateAdminResponse(grpcResp)
	return resp, nil
}

func (g *AdminGateway) DeleteAdmin(ctx context.Context, adminID string) (*adminResp.MessageResponse, error) {

	grpcResp, err := g.client.Delete(
		ctx,
		mapper.ToProtoDeleteAdmin(adminID),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	resp := mapper.ToDeleteAdminResponse(grpcResp)
	return resp, nil
}
