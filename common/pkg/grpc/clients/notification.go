package clients

import (
	commongrpc "github.com/Eucastan/eucastanpay/common/pkg/grpc"

	notifypb "github.com/Eucastan/eucastanpay/common/proto/notification"
)

func Notification(manager *commongrpc.Manager) notifypb.NotificationServiceClient {
	return notifypb.NewNotificationServiceClient(
		manager.Get("notification"),
	)
}
