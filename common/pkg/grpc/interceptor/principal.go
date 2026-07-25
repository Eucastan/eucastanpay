package interceptor

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Principal struct {
	IsUser  bool
	IsAdmin bool

	UserID string
	Role   string

	AdminID   string
	AdminRole string
}

func CurrentPrincipal(ctx context.Context) (*Principal, error) {

	p := &Principal{}

	if v, ok := ctx.Value(ContextUserID).(string); ok && v != "" {
		p.IsUser = true
		p.UserID = v

		if role, ok := ctx.Value(ContextRole).(string); ok {
			p.Role = role
		}
	}

	if v, ok := ctx.Value(ContextAdminID).(string); ok && v != "" {
		p.IsAdmin = true
		p.AdminID = v

		if role, ok := ctx.Value(ContextAdminRole).(string); ok {
			p.AdminRole = role
		}
	}

	if !p.IsUser && !p.IsAdmin {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	return p, nil
}

func RequireUser(ctx context.Context) (*Principal, error) {

	p, err := CurrentPrincipal(ctx)
	if err != nil {
		return nil, err
	}

	if !p.IsUser {
		return nil, status.Error(codes.PermissionDenied, "user authentication required")
	}

	return p, nil
}

func RequireAdmin(ctx context.Context) (*Principal, error) {

	p, err := CurrentPrincipal(ctx)
	if err != nil {
		return nil, err
	}

	if !p.IsAdmin {
		return nil, status.Error(codes.PermissionDenied, "admin authentication required")
	}

	return p, nil
}

func RequireUserOwner(
	ctx context.Context,
	resourceUserID string,
) (*Principal, error) {

	p, err := RequireUser(ctx)
	if err != nil {
		return nil, err
	}

	if p.UserID != resourceUserID {
		return nil, status.Error(
			codes.PermissionDenied,
			"permission denied",
		)
	}

	return p, nil
}

func RequireAdminRole(
	ctx context.Context,
	roles ...string,
) (*Principal, error) {

	p, err := RequireAdmin(ctx)
	if err != nil {
		return nil, err
	}

	for _, r := range roles {
		if p.AdminRole == r {
			return p, nil
		}
	}

	return nil, status.Error(
		codes.PermissionDenied,
		"insufficient privileges",
	)
}
