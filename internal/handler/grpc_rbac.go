package handler

import (
	"context"

	"github.com/UniqueStudio/UniqueSSOBackend/internal/model"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/proto/sso"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/tracer"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SSOServer struct {
	sso.UnimplementedSSOServiceServer
}

var _ sso.SSOServiceServer = &SSOServer{}

func NewSSOServer() *SSOServer { return &SSOServer{} }

func (*SSOServer) CheckPermission(ctx context.Context, req *sso.CheckPermissionRequest) (*sso.CheckPermissionResponse, error) {
	return nil, nil
}

func (*SSOServer) GetUserByUID(ctx context.Context, req *sso.GetUserByUIDRequest) (*sso.GetUserByUIDResponse, error) {
	apmCtx, span := tracer.Tracer.Start(ctx, "GetUserByUID")
	defer span.End()
	user, err := model.GetUserByUID(apmCtx, req.GetUid())
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &sso.GetUserByUIDResponse{
		Uid:   user.UID,
		Phone: user.Phone,
		Email: user.Email,
		Name:  user.Name,
		//JoinTime:    timestamppb.New(user.JoinTime),
		JoinTime:  user.JoinTime,
		AvatarUrl: user.AvatarURL,
		Gender:    sso.Gender(user.Gender),
		Groups:    user.Groups,
	}, nil
}

func (*SSOServer) GetRolesByUID(ctx context.Context, req *sso.GetRolesByUIDRequest) (*sso.GetRolesByUIDResponse, error) {
	apmCtx, span := tracer.Tracer.Start(ctx, "GetRolesByUID")
	defer span.End()

	roles, err := model.GetRolesByUID(apmCtx, req.GetUid())
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &sso.GetRolesByUIDResponse{
		Roles: roles,
	}, nil
}

func (*SSOServer) GetUsers(ctx context.Context, req *sso.GetUsersRequest) (*sso.GetUsersResponse, error) {
	apmCtx, span := tracer.Tracer.Start(ctx, "GetUsers")
	defer span.End()

	users, err := model.GetUsersByUids(apmCtx, req.GetUid())
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	resp := &sso.GetUsersResponse{}
	for _, user := range users {
		resp.Users = append(resp.Users, &sso.GetUserByUIDResponse{
			Uid:       user.UID,
			Phone:     user.Phone,
			Email:     user.Email,
			Name:      user.Name,
			JoinTime:  user.JoinTime,
			AvatarUrl: user.AvatarURL,
			Gender:    sso.Gender(user.Gender),
			Groups:    user.Groups,
		})
	}

	return resp, nil
}

func (*SSOServer) GetGroupsDetail(ctx context.Context, empty *emptypb.Empty) (*sso.GetGroupsDetailResponse, error) {
	apmCtx, span := tracer.Tracer.Start(ctx, "GetGroupsDetail")
	defer span.End()

	groupsDetail, err := model.GetGroupsDetail(apmCtx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	groupsVal := make(map[string]*structpb.Value)
	for group, count := range groupsDetail {
		groupsVal[group] = &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: float64(count)}}
	}

	resp := &sso.GetGroupsDetailResponse{Groups: &structpb.Struct{Fields: groupsVal}}
	return resp, nil
}
