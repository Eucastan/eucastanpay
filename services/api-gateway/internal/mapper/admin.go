package mapper

import (
	adminpb "github.com/Eucastan/eucastanpay/common/proto/admin"

	adminReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/admin"

	adminResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/admin"
)

func ToProtoCreateAdmin(req *adminReq.CreateAdminRequest) *adminpb.RegisterRequest {
	return &adminpb.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
	}
}

func ToProtoAdminLogin(req *adminReq.AdminLoginRequest) *adminpb.LoginRequest {
	return &adminpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}
}

func ToProtoListAdmins(limit, page int) *adminpb.ListAdminsRequest {
	return &adminpb.ListAdminsRequest{
		Page:  int32(page),
		Limit: int32(limit),
	}
}

func ToProtoGetAdmin(adminID string) *adminpb.GetAdminByIDRequest {
	return &adminpb.GetAdminByIDRequest{
		AdminId: adminID,
	}
}

func ToProtoUpdateAdmin(adminID string, req *adminReq.UpdateAdminRequest) *adminpb.UpdateRequest {
	return &adminpb.UpdateRequest{
		AdminId:   adminID,
		Password:  *req.Password,
		FirstName: *req.FirstName,
		LastName:  *req.LastName,
		Status:    *req.Status,
		Role:      *req.Role,
	}
}

func ToProtoDeleteAdmin(adminID string) *adminpb.GetAdminByIDRequest {
	return &adminpb.GetAdminByIDRequest{
		AdminId: adminID,
	}
}

func ToUpdateAdminResponse(resp *adminpb.ActionResponse) *adminResp.MessageResponse {
	return &adminResp.MessageResponse{
		Message: resp.Message,
	}
}

func ToDeleteAdminResponse(resp *adminpb.ActionResponse) *adminResp.MessageResponse {
	return &adminResp.MessageResponse{
		Message: resp.Message,
	}
}

func ToCreateAdminResponse(resp *adminpb.AdminResponse) adminResp.AdminResponse {

	return adminResp.AdminResponse{
		ID:        resp.AdminId,
		Email:     resp.Email,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
	}
}

func ToLoginAdminResponse(r *adminpb.LoginResponse) adminResp.AdminLoginResponse {

	return adminResp.AdminLoginResponse{
		Message: r.Message,
		Data: adminResp.AdminResponse{
			ID:        r.Resp.AdminId,
			Email:     r.Resp.Email,
			FirstName: r.Resp.FirstName,
			LastName:  r.Resp.LastName,
			Status:    r.Resp.Status,
			Role:      r.Resp.Role,
			CreatedAt: r.Resp.CreatedAt.AsTime(),
			UpdatedAt: r.Resp.UpdatedAt.AsTime(),
		},
		AccessToken:  r.AccessToken,
		RefreshToken: r.RefreshToken,
	}
}

func ToGetAdminResponse(r *adminpb.AdminResponse) *adminResp.AdminResponse {
	return &adminResp.AdminResponse{
		ID:        r.AdminId,
		Email:     r.Email,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Status:    r.Status,
		Role:      r.Role,
		CreatedAt: r.CreatedAt.AsTime(),
		UpdatedAt: r.UpdatedAt.AsTime(),
	}
}

func ToListAdminsResponse(req *adminpb.ListAdminsResponse) *adminResp.ListAdminsResponse {
	data := make([]*adminResp.AdminResponse, 0, len(req.Admins))
	for _, r := range req.Admins {
		data = append(data, &adminResp.AdminResponse{
			ID:        r.AdminId,
			Email:     r.Email,
			FirstName: r.FirstName,
			LastName:  r.LastName,
			Status:    r.Status,
			Role:      r.Role,
			CreatedAt: r.CreatedAt.AsTime(),
			UpdatedAt: r.UpdatedAt.AsTime(),
		})
	}

	return &adminResp.ListAdminsResponse{
		Data: data,
	}
}
