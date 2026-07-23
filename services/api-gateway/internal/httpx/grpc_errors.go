package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGRPCError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	st, ok := status.FromError(err)
	if !ok {
		Error(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)

		return
	}

	switch st.Code() {

	case codes.InvalidArgument:

		Error(
			c,
			http.StatusBadRequest,
			st.Message(),
		)

	case codes.NotFound:

		Error(
			c,
			http.StatusNotFound,
			st.Message(),
		)

	case codes.AlreadyExists:

		Error(
			c,
			http.StatusConflict,
			st.Message(),
		)

	case codes.Unauthenticated:

		Error(
			c,
			http.StatusUnauthorized,
			st.Message(),
		)

	case codes.PermissionDenied:

		Error(
			c,
			http.StatusForbidden,
			st.Message(),
		)

	case codes.ResourceExhausted:

		Error(
			c,
			http.StatusTooManyRequests,
			st.Message(),
		)

	case codes.DeadlineExceeded:

		Error(
			c,
			http.StatusGatewayTimeout,
			"service timeout",
		)

	case codes.Unavailable:

		Error(
			c,
			http.StatusServiceUnavailable,
			"service unavailable",
		)

	case codes.FailedPrecondition:

		Error(
			c,
			http.StatusPreconditionFailed,
			st.Message(),
		)

	case codes.Aborted:

		Error(
			c,
			http.StatusConflict,
			st.Message(),
		)

	case codes.Internal:

		Error(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)

	default:

		Error(
			c,
			http.StatusInternalServerError,
			st.Message(),
		)
	}
}
