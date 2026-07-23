package mapper

import (
	userpb "github.com/Eucastan/eucastanpay/common/proto/user"
	userResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/user"
)

func ToProtoCreateKYC(idNumber, idType string) *userpb.KycRequest {
	return &userpb.KycRequest{
		IdNumber: idNumber,
		IdType:   idType,
	}
}

func ToProtoGetKYCByUserID(userID string) *userpb.KycIdRequest {
	return &userpb.KycIdRequest{
		UserId: userID,
	}
}

func ToKYCResponse(resp *userpb.KycResponse) *userResp.KYCResponse {
	return &userResp.KYCResponse{
		Message: resp.Message,
	}
}
