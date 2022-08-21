// Code generated by mockery v2.14.0. DO NOT EDIT.

package services

import (
	context "context"

	types "github.com/leveldorado/space-trouble/pkg/types"
	mock "github.com/stretchr/testify/mock"
)

// mockLaunchpadFirstDestinationRepo is an autogenerated mock type for the launchpadFirstDestinationRepo type
type mockLaunchpadFirstDestinationRepo struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, launchpad
func (_m *mockLaunchpadFirstDestinationRepo) Get(ctx context.Context, launchpad string) (types.LaunchpadFirstDestination, error) {
	ret := _m.Called(ctx, launchpad)

	var r0 types.LaunchpadFirstDestination
	if rf, ok := ret.Get(0).(func(context.Context, string) types.LaunchpadFirstDestination); ok {
		r0 = rf(ctx, launchpad)
	} else {
		r0 = ret.Get(0).(types.LaunchpadFirstDestination)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, launchpad)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTnewMockLaunchpadFirstDestinationRepo interface {
	mock.TestingT
	Cleanup(func())
}

// newMockLaunchpadFirstDestinationRepo creates a new instance of mockLaunchpadFirstDestinationRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockLaunchpadFirstDestinationRepo(t mockConstructorTestingTnewMockLaunchpadFirstDestinationRepo) *mockLaunchpadFirstDestinationRepo {
	mock := &mockLaunchpadFirstDestinationRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}