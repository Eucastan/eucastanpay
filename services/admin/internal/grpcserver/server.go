package grpcserver

import (
	"context"

	"github.com/Eucastan/eucastanpay/common/pkg/grpcstatus"
	adminpb "github.com/Eucastan/eucastanpay/common/proto/admin"
	"github.com/Eucastan/eucastanpay/services/admin/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/admin/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AdminServer struct {
	adminpb.UnimplementedAdminServiceServer
	admin usecase.AdminUseCase
}

func NewAdminServiceServer(admin usecase.AdminUseCase) *AdminServer {
	return &AdminServer{
		admin: admin,
	}
}

func (s *AdminServer) Register(ctx context.Context, req *adminpb.RegisterRequest) (*adminpb.AdminResponse, error) {
	input := &request.CreateAdminRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
	}

	admin, err := s.admin.CreateAdmin(ctx, input)
	if err != nil {
		return nil, grpcstatus.ToAdminStatus(err)
	}

	return &adminpb.AdminResponse{
		AdminId:   admin.ID,
		Email:     admin.Email,
		FirstName: admin.FirstName,
		LastName:  admin.LastName,
		Status:    admin.Status,
		CreatedAt: timestamppb.New(admin.CreatedAt),
		UpdatedAt: timestamppb.New(admin.UpdatedAt),
	}, nil
}

func (s *AdminServer) Login(ctx context.Context, req *adminpb.LoginRequest) (*adminpb.LoginResponse, error) {
	input := &request.AdminLoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	admin, err := s.admin.Login(ctx, input)
	if err != nil {
		return nil, grpcstatus.ToAdminStatus(err)
	}

	data := adminpb.Admin{
		AdminId:   admin.Data.ID,
		Email:     admin.Data.Email,
		FirstName: admin.Data.FirstName,
		LastName:  admin.Data.LastName,
		Status:    admin.Data.Status,
		Role:      admin.Data.Role,
		CreatedAt: timestamppb.New(admin.Data.CreatedAt),
		UpdatedAt: timestamppb.New(admin.Data.UpdatedAt),
	}

	return &adminpb.LoginResponse{
		Message:      "User login successful",
		Resp:         &data,
		AccessToken:  admin.AccessToken,
		RefreshToken: admin.RefreshToken,
	}, nil
}

func (s *AdminServer) GetAdminByID(ctx context.Context, req *adminpb.GetAdminByIDRequest) (*adminpb.AdminResponse, error) {
	adminID := ctx.Value("user_id")
	if adminID == "" {
		return nil, status.Error(codes.Unauthenticated, "id not found in context")
	}

	idStr, ok := adminID.(string)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid user id format")
	}

	if req.AdminId != idStr {
		return nil, status.Error(codes.PermissionDenied, "unauthorized access")
	}

	admin, err := s.admin.GetAdminByID(ctx, req.AdminId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &adminpb.AdminResponse{
		AdminId:   admin.ID,
		Email:     admin.Email,
		FirstName: admin.FirstName,
		LastName:  admin.LastName,
		Status:    admin.Status,
		Role:      admin.Role,
		CreatedAt: timestamppb.New(admin.CreatedAt),
		UpdatedAt: timestamppb.New(admin.UpdatedAt),
	}, nil
}

func (s *AdminServer) GetAllUsers(
	ctx context.Context,
	req *adminpb.ListAdminsRequest,
) (*adminpb.ListAdminsResponse, error) {
	admins, err := s.admin.ListAdmins(ctx, int(req.Limit), int(req.Page))
	if err != nil {
		return nil, grpcstatus.ToAdminStatus(err)
	}

	resp := make([]*adminpb.Admin, 0, len(admins))
	for _, v := range admins {
		resp = append(resp, &adminpb.Admin{
			AdminId:   v.ID,
			Email:     v.Email,
			FirstName: v.FirstName,
			LastName:  v.LastName,
			Status:    v.Status,
			Role:      v.Role,
			CreatedAt: timestamppb.New(v.CreatedAt),
			UpdatedAt: timestamppb.New(v.UpdatedAt),
		})
	}

	return &adminpb.ListAdminsResponse{
		Admins: resp,
	}, nil
}

func (s *AdminServer) Logout(ctx context.Context, req *adminpb.GetAdminByIDRequest) (*adminpb.ActionResponse, error) {
	if err := s.admin.LogoutByAdminID(ctx, req.AdminId); err != nil {
		return nil, grpcstatus.ToAdminStatus(err)
	}

	return &adminpb.ActionResponse{
		Message: "User logged out from every device",
	}, nil
}

func (s *AdminServer) Update(ctx context.Context, req *adminpb.UpdateRequest) (*adminpb.ActionResponse, error) {
	input := request.UpdateAdminRequest{
		Status: &req.Status,
		Role:   &req.Role,
	}

	_, err := s.admin.UpdateAdmin(ctx, req.AdminId, &input)
	if err != nil {
		return nil, grpcstatus.ToAdminStatus(err)
	}

	return &adminpb.ActionResponse{
		Message: "User deleted successfully",
	}, nil
}

func (s *AdminServer) Delete(ctx context.Context, req *adminpb.GetAdminByIDRequest) (*adminpb.ActionResponse, error) {
	if err := s.admin.DeleteAdmin(ctx, req.AdminId); err != nil {
		return nil, grpcstatus.ToAdminStatus(err)
	}

	return &adminpb.ActionResponse{
		Message: "User deleted successfully",
	}, nil
}
