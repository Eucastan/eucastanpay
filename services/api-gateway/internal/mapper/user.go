package mapper

import (
	userpb "github.com/Eucastan/eucastanpay/common/proto/user"

	userReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/user"

	userResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/user"
)

func ToProtoRegister(req *userReq.RegisterRequest) *userpb.RegisterRequest {
	return &userpb.RegisterRequest{
		Email:       req.Email,
		Phone:       req.Phone,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Password:    req.Password,
		DateOfBirth: req.DateOfBirth,
	}
}

func ToProtoLogin(req *userReq.LoginRequest) *userpb.LoginRequest {
	return &userpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}
}

func ToProtoListUsers() *userpb.ListUsersRequest {
	return &userpb.ListUsersRequest{}
}

func ToProtoGetUser(userID string) *userpb.GetUserByIDRequest {
	return &userpb.GetUserByIDRequest{
		UserId: userID,
	}
}

func ToProtoUpdateUser(userID string, req *userReq.UpdateRequest) *userpb.UpdateRequest {
	return &userpb.UpdateRequest{
		UserId:        userID,
		Password:      req.Password,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Status:        req.Status,
		EmailVerified: req.EmailVerified,
	}
}

func ToProtoDeleteUser(userID string) *userpb.GetUserByIDRequest {
	return &userpb.GetUserByIDRequest{
		UserId: userID,
	}
}

func ToUpdateResponse(resp *userpb.ActionResponse) *userResp.MessageResponse {
	return &userResp.MessageResponse{
		Message: resp.Message,
	}
}

func ToDeleteResponse(resp *userpb.ActionResponse) *userResp.MessageResponse {
	return &userResp.MessageResponse{
		Message: resp.Message,
	}
}

func ToRegisterResponse(resp *userpb.UserResponse) userResp.RegisterResponse {

	return userResp.RegisterResponse{
		Message: "registration successful",
		User: userResp.UserResponse{
			ID:        resp.UserId,
			Email:     resp.Email,
			FirstName: resp.FirstName,
			LastName:  resp.LastName,
		},
	}
}

func ToLoginResponse(r *userpb.LoginResponse) userResp.AuthResponse {

	return userResp.AuthResponse{
		Message: r.Message,
		User: userResp.UserResponse{
			ID:            r.Resp.UserId,
			Email:         r.Resp.Email,
			Phone:         r.Resp.Phone,
			FirstName:     r.Resp.FirstName,
			LastName:      r.Resp.LastName,
			Status:        r.Resp.Status,
			EmailVerified: r.Resp.EmailVerified,
			CreatedAt:     r.Resp.CreatedAt.AsTime(),
		},
		AccessToken:  r.AccessToken,
		RefreshToken: r.RefreshToken,
	}
}

func ToGetUserResponse(r *userpb.UserResponse) *userResp.UserResponse {
	return &userResp.UserResponse{
		ID:            r.UserId,
		Email:         r.Email,
		Phone:         r.Phone,
		FirstName:     r.FirstName,
		LastName:      r.LastName,
		Status:        r.Status,
		EmailVerified: r.EmailVerified,
		CreatedAt:     r.CreatedAt.AsTime(),
	}
}

func ToListUsersResponse(req *userpb.ListUsersResponse) *userResp.ListUsersResponse {
	data := make([]*userResp.UserResponse, 0, len(req.Users))
	for _, r := range req.Users {
		data = append(data, &userResp.UserResponse{
			ID:            r.UserId,
			Email:         r.Email,
			Phone:         r.Phone,
			FirstName:     r.FirstName,
			LastName:      r.LastName,
			Status:        r.Status,
			EmailVerified: r.EmailVerified,
			CreatedAt:     r.CreatedAt.AsTime(),
		})
	}

	return &userResp.ListUsersResponse{
		User: data,
	}
}
