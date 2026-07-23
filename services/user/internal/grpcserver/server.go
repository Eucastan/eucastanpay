package grpcserver

import (
	"context"
	"fmt"

	"github.com/Eucastan/eucastanpay/common/pkg/grpcstatus"
	userpb "github.com/Eucastan/eucastanpay/common/proto/user"
	"github.com/Eucastan/eucastanpay/services/user/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/user/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
	User usecase.UserUseCaseInterface
	Kyc  usecase.KYCUseCase
}

func NewUserServiceServer(user usecase.UserUseCaseInterface, kyc usecase.KYCUseCase) *UserServer {
	return &UserServer{
		User: user,
		Kyc:  kyc,
	}
}

func (s *UserServer) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.UserResponse, error) {
	input := &request.RegisterRequest{
		Email:       req.Email,
		Phone:       req.Phone,
		Password:    req.Password,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		DateOfBirth: req.DateOfBirth,
	}

	user, err := s.User.Register(ctx, input)
	if err != nil {
		return nil, grpcstatus.ToUserStatus(err)
	}

	return &userpb.UserResponse{
		UserId:        user.ID,
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Status:        user.Status,
		EmailVerified: user.EmailVerified,
		CreatedAt:     timestamppb.New(user.CreatedAt),
	}, nil
}

func (s *UserServer) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	input := &request.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := s.User.Login(ctx, input)
	if err != nil {
		return nil, grpcstatus.ToUserStatus(err)
	}

	data := userpb.User{
		UserId:        user.User.ID,
		Email:         user.User.Email,
		Phone:         user.User.Phone,
		FirstName:     user.User.FirstName,
		LastName:      user.User.LastName,
		Status:        user.User.Status,
		EmailVerified: user.User.EmailVerified,
		CreatedAt:     timestamppb.New(user.User.CreatedAt),
	}

	return &userpb.LoginResponse{
		Message:      "User login successful",
		Resp:         &data,
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
	}, nil
}

func (s *UserServer) GetUserByID(ctx context.Context, req *userpb.GetUserByIDRequest) (*userpb.UserResponse, error) {
	userID := ctx.Value("user_id")
	idStr, ok := userID.(string)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "user ID not found in context")
	}

	if req.UserId != idStr {
		return nil, status.Error(codes.PermissionDenied, "unauthorized access")
	}

	userInfo, err := s.User.GetUserByID(ctx, req.UserId)
	if err != nil {
		return nil, grpcstatus.ToUserStatus(err)
	}

	return &userpb.UserResponse{
		UserId:        userInfo.ID,
		Email:         userInfo.Email,
		FirstName:     userInfo.FirstName,
		LastName:      userInfo.LastName,
		Status:        userInfo.Status,
		EmailVerified: userInfo.EmailVerified,
		CreatedAt:     timestamppb.New(userInfo.CreatedAt),
	}, nil
}

func (s *UserServer) GetAllUsers(
	ctx context.Context,
	req *userpb.ListUsersRequest,
) (*userpb.ListUsersResponse, error) {
	users, err := s.User.GetAllUsers(ctx)
	if err != nil {
		return nil, grpcstatus.ToUserStatus(err)
	}

	resp := make([]*userpb.User, 0, len(users))
	for _, v := range users {
		resp = append(resp, &userpb.User{
			UserId:        v.ID,
			Email:         v.Email,
			Phone:         v.Phone,
			FirstName:     v.FirstName,
			LastName:      v.LastName,
			Status:        v.Status,
			EmailVerified: v.EmailVerified,
			CreatedAt:     timestamppb.New(v.CreatedAt),
		})
	}

	return &userpb.ListUsersResponse{
		Users: resp,
	}, nil
}

func (s *UserServer) LogoutAllUser(ctx context.Context, req *userpb.GetUserByIDRequest) (*userpb.ActionResponse, error) {
	if err := s.User.LogoutAllUsers(ctx, req.UserId); err != nil {
		return nil, grpcstatus.ToUserStatus(err)
	}

	return &userpb.ActionResponse{
		Message: "User logged out from every device",
	}, nil
}

func (s *UserServer) ActionOnUser(ctx context.Context, req *userpb.ActionRequest) (*userpb.ActionResponse, error) {
	msg, err := s.User.UserCurrentStatus(ctx, req.UserId, req.Status)
	if err != nil {
		return nil, grpcstatus.ToUserStatus(err)
	}

	return &userpb.ActionResponse{
		Message: msg,
	}, nil
}

func (s *UserServer) Update(ctx context.Context, req *userpb.UpdateRequest) (*userpb.ActionResponse, error) {
	input := request.UpdateRequest{
		Password:      req.Password,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Status:        req.Status,
		EmailVerified: req.EmailVerified,
	}

	if err := s.User.Update(ctx, req.UserId, &input); err != nil {
		return nil, grpcstatus.ToUserStatus(err)
	}

	return &userpb.ActionResponse{
		Message: "User deleted successfully",
	}, nil
}

func (s *UserServer) Delete(ctx context.Context, req *userpb.GetUserByIDRequest) (*userpb.ActionResponse, error) {
	if err := s.User.DeleteUser(ctx, req.UserId); err != nil {
		return nil, grpcstatus.ToUserStatus(err)
	}

	return &userpb.ActionResponse{
		Message: "User deleted successfully",
	}, nil
}

func (s *UserServer) CreateKYC(ctx context.Context, req *userpb.KycRequest) (*userpb.KycResponse, error) {
	input := &request.KYCRequest{
		IDType:   req.IdType,
		IDNumber: req.IdNumber,
	}

	userId := ctx.Value("user_id")
	if userId == "" {
		return nil, status.Error(codes.Unauthenticated, "user id not found in context")
	}

	userID, ok := userId.(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid user id format")
	}

	err := s.Kyc.CreateKYC(ctx, userID, input)
	if err != nil {
		return nil, grpcstatus.ToUserStatus(err)
	}

	return &userpb.KycResponse{
		Message: "Created successfully",
	}, nil
}

func (s *UserServer) ApproveKYC(ctx context.Context, req *userpb.KycIdRequest) (*userpb.KycResponse, error) {

	err := s.Kyc.ApproveKYC(ctx, req.UserId)
	if err != nil {
		return nil, grpcstatus.ToUserStatus(err)
	}

	return &userpb.KycResponse{
		Message: "Your KYC is approved",
	}, nil
}

func (s *UserServer) GetKYC(ctx context.Context, req *userpb.KycIdRequest) (*userpb.KycResponse, error) {
	kyc, err := s.Kyc.GetKYC(ctx, req.UserId)
	if err != nil {
		return nil, grpcstatus.ToUserStatus(err)
	}

	return &userpb.KycResponse{
		Message: fmt.Sprintf("%s:%s", kyc.Status, kyc.Message),
	}, nil
}
