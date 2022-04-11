// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package optimus

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SecretServiceClient is the client API for SecretService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SecretServiceClient interface {
	// RegisterSecret creates a new secret of a project
	RegisterSecret(ctx context.Context, in *RegisterSecretRequest, opts ...grpc.CallOption) (*RegisterSecretResponse, error)
	// UpdateSecret updates secret at project level
	UpdateSecret(ctx context.Context, in *UpdateSecretRequest, opts ...grpc.CallOption) (*UpdateSecretResponse, error)
	// ListSecrets shows the secrets registered for a project
	ListSecrets(ctx context.Context, in *ListSecretsRequest, opts ...grpc.CallOption) (*ListSecretsResponse, error)
	// DeleteSecret deletes a secret for a project
	DeleteSecret(ctx context.Context, in *DeleteSecretRequest, opts ...grpc.CallOption) (*DeleteSecretResponse, error)
}

type secretServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSecretServiceClient(cc grpc.ClientConnInterface) SecretServiceClient {
	return &secretServiceClient{cc}
}

func (c *secretServiceClient) RegisterSecret(ctx context.Context, in *RegisterSecretRequest, opts ...grpc.CallOption) (*RegisterSecretResponse, error) {
	out := new(RegisterSecretResponse)
	err := c.cc.Invoke(ctx, "/odpf.optimus.core.v1beta1.SecretService/RegisterSecret", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *secretServiceClient) UpdateSecret(ctx context.Context, in *UpdateSecretRequest, opts ...grpc.CallOption) (*UpdateSecretResponse, error) {
	out := new(UpdateSecretResponse)
	err := c.cc.Invoke(ctx, "/odpf.optimus.core.v1beta1.SecretService/UpdateSecret", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *secretServiceClient) ListSecrets(ctx context.Context, in *ListSecretsRequest, opts ...grpc.CallOption) (*ListSecretsResponse, error) {
	out := new(ListSecretsResponse)
	err := c.cc.Invoke(ctx, "/odpf.optimus.core.v1beta1.SecretService/ListSecrets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *secretServiceClient) DeleteSecret(ctx context.Context, in *DeleteSecretRequest, opts ...grpc.CallOption) (*DeleteSecretResponse, error) {
	out := new(DeleteSecretResponse)
	err := c.cc.Invoke(ctx, "/odpf.optimus.core.v1beta1.SecretService/DeleteSecret", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SecretServiceServer is the server API for SecretService service.
// All implementations must embed UnimplementedSecretServiceServer
// for forward compatibility
type SecretServiceServer interface {
	// RegisterSecret creates a new secret of a project
	RegisterSecret(context.Context, *RegisterSecretRequest) (*RegisterSecretResponse, error)
	// UpdateSecret updates secret at project level
	UpdateSecret(context.Context, *UpdateSecretRequest) (*UpdateSecretResponse, error)
	// ListSecrets shows the secrets registered for a project
	ListSecrets(context.Context, *ListSecretsRequest) (*ListSecretsResponse, error)
	// DeleteSecret deletes a secret for a project
	DeleteSecret(context.Context, *DeleteSecretRequest) (*DeleteSecretResponse, error)
	mustEmbedUnimplementedSecretServiceServer()
}

// UnimplementedSecretServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSecretServiceServer struct {
}

func (UnimplementedSecretServiceServer) RegisterSecret(context.Context, *RegisterSecretRequest) (*RegisterSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterSecret not implemented")
}
func (UnimplementedSecretServiceServer) UpdateSecret(context.Context, *UpdateSecretRequest) (*UpdateSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSecret not implemented")
}
func (UnimplementedSecretServiceServer) ListSecrets(context.Context, *ListSecretsRequest) (*ListSecretsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSecrets not implemented")
}
func (UnimplementedSecretServiceServer) DeleteSecret(context.Context, *DeleteSecretRequest) (*DeleteSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSecret not implemented")
}
func (UnimplementedSecretServiceServer) mustEmbedUnimplementedSecretServiceServer() {}

// UnsafeSecretServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SecretServiceServer will
// result in compilation errors.
type UnsafeSecretServiceServer interface {
	mustEmbedUnimplementedSecretServiceServer()
}

func RegisterSecretServiceServer(s grpc.ServiceRegistrar, srv SecretServiceServer) {
	s.RegisterService(&SecretService_ServiceDesc, srv)
}

func _SecretService_RegisterSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretServiceServer).RegisterSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.optimus.core.v1beta1.SecretService/RegisterSecret",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretServiceServer).RegisterSecret(ctx, req.(*RegisterSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SecretService_UpdateSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretServiceServer).UpdateSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.optimus.core.v1beta1.SecretService/UpdateSecret",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretServiceServer).UpdateSecret(ctx, req.(*UpdateSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SecretService_ListSecrets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSecretsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretServiceServer).ListSecrets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.optimus.core.v1beta1.SecretService/ListSecrets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretServiceServer).ListSecrets(ctx, req.(*ListSecretsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SecretService_DeleteSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretServiceServer).DeleteSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.optimus.core.v1beta1.SecretService/DeleteSecret",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretServiceServer).DeleteSecret(ctx, req.(*DeleteSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SecretService_ServiceDesc is the grpc.ServiceDesc for SecretService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SecretService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "odpf.optimus.core.v1beta1.SecretService",
	HandlerType: (*SecretServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterSecret",
			Handler:    _SecretService_RegisterSecret_Handler,
		},
		{
			MethodName: "UpdateSecret",
			Handler:    _SecretService_UpdateSecret_Handler,
		},
		{
			MethodName: "ListSecrets",
			Handler:    _SecretService_ListSecrets_Handler,
		},
		{
			MethodName: "DeleteSecret",
			Handler:    _SecretService_DeleteSecret_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "odpf/optimus/core/v1beta1/secret.proto",
}