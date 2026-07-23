package gateway

import (
	"context"

	userpb "github.com/Eucastan/eucastanpay/common/proto/user"
	userReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/user"
	userResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/user"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/mapper"
)

type UserGateway struct {
	client userpb.UserServiceClient
}

func NewUserGateway(client userpb.UserServiceClient) *UserGateway {
	return &UserGateway{
		client: client,
	}
}

func (s *UserGateway) Register(ctx context.Context, req *userReq.RegisterRequest) (*userResp.RegisterResponse, error) {

	grpcResp, err := s.client.Register(
		ctx,
		mapper.ToProtoRegister(req),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToRegisterResponse(grpcResp)

	return &resp, nil
}

func (s *UserGateway) Login(ctx context.Context, req *userReq.LoginRequest) (*userResp.AuthResponse, error) {

	grpcResp, err := s.client.Login(
		ctx,
		mapper.ToProtoLogin(req),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToLoginResponse(grpcResp)

	return &resp, nil
}

func (s *UserGateway) GetAllUsers(ctx context.Context) (*userResp.ListUsersResponse, error) {

	grpcResp, err := s.client.GetAllUsers(
		ctx,
		mapper.ToProtoListUsers(),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToListUsersResponse(grpcResp)

	return resp, nil
}

func (s *UserGateway) GetUserByID(ctx context.Context, userID string) (*userResp.UserResponse, error) {

	grpcResp, err := s.client.GetUserByID(
		ctx,
		mapper.ToProtoGetUser(userID),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToGetUserResponse(grpcResp)

	return resp, nil
}

func (s *UserGateway) UpdateUser(ctx context.Context, userID string, req *userReq.UpdateRequest) (*userResp.MessageResponse, error) {

	grpcResp, err := s.client.Update(
		ctx,
		mapper.ToProtoUpdateUser(userID, req),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToUpdateResponse(grpcResp)

	return resp, nil
}

func (s *UserGateway) DeleteUser(ctx context.Context, userID string) (*userResp.MessageResponse, error) {

	grpcResp, err := s.client.Delete(
		ctx,
		mapper.ToProtoDeleteUser(userID),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToDeleteResponse(grpcResp)

	return resp, nil
}

func (s *UserGateway) CreateKYC(ctx context.Context, idNumber, idType string) (*userResp.KYCResponse, error) {

	grpcResp, err := s.client.CreateKYC(
		ctx,
		mapper.ToProtoCreateKYC(idNumber, idType),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToKYCResponse(grpcResp)

	return resp, nil
}

func (s *UserGateway) GetKYC(ctx context.Context, userID string) (*userResp.KYCResponse, error) {

	grpcResp, err := s.client.GetKYC(
		ctx,
		mapper.ToProtoGetKYCByUserID(userID),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToKYCResponse(grpcResp)

	return resp, nil
}

func (s *UserGateway) ApproveKYC(ctx context.Context, userID string) (*userResp.KYCResponse, error) {

	grpcResp, err := s.client.ApproveKYC(
		ctx,
		mapper.ToProtoGetKYCByUserID(userID),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToKYCResponse(grpcResp)

	return resp, nil
}
