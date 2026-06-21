package grpcserver

import (
	"context"

	"github.com/Eucastan/eucastanpay/common/proto/user"
	"github.com/Eucastan/eucastanpay/services/user/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserServer struct {
	user.UnimplementedUserServiceServer
	User usecase.UserUseCaseInterface
	Kyc  usecase.KYCUseCase
}

func NewUserServiceServer(user usecase.UserUseCaseInterface, kyc usecase.KYCUseCase) *UserServer {
	return &UserServer{
		User: user,
		Kyc:  kyc,
	}
}

func (s *UserServer) GetUserByID(ctx context.Context, req *user.GetUserByIDRequest) (*user.UserResponse, error) {
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
		return nil, status.Error(codes.Internal, err.Error())
	}

	kyc, err := s.Kyc.GetKYC(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &user.UserResponse{
		UserId:    userInfo.ID,
		Email:     userInfo.Email,
		FirstName: userInfo.FirstName,
		LastName:  userInfo.LastName,
		KycStatus: kyc.Status,
	}, nil
}

func (s *UserServer) GetAllUsers(
	ctx context.Context,
	req *user.ListUsersRequest,
) (*user.ListUsersResponse, error) {
	users, err := s.User.GetAllUsers(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get all users")
	}

	resp := make([]*user.User, 0, len(users))
	for _, v := range users {
		resp = append(resp, &user.User{
			UserId:        v.ID,
			Email:         v.Email,
			Phone:         v.Phone,
			FirstName:     v.FirstName,
			LastName:      v.LastName,
			DateOfBirth:   "",
			Status:        v.Status,
			EmailVerified: v.EmailVerified,
			CreatedAt:     timestamppb.New(v.CreatedAt),
		})
	}

	return &user.ListUsersResponse{
		Users: resp,
	}, nil
}

func (s *UserServer) ActionOnUser(ctx context.Context, req *user.ActionRequest) (*user.ActionResponse, error) {
	msg, err := s.User.UserCurrentStatus(ctx, req.UserId, req.Status)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &user.ActionResponse{
		Message: msg,
	}, nil
}
