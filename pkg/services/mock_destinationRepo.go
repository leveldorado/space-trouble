// Code generated by mockery v2.14.0. DO NOT EDIT.

package services

import (
	context "context"

	types "github.com/leveldorado/space-trouble/pkg/types"
	mock "github.com/stretchr/testify/mock"
)

// mockDestinationRepo is an autogenerated mock type for the destinationRepo type
type mockDestinationRepo struct {
	mock.Mock
}

// ListSorted provides a mock function with given fields: ctx
func (_m *mockDestinationRepo) ListSorted(ctx context.Context) ([]types.Destination, error) {
	ret := _m.Called(ctx)

	var r0 []types.Destination
	if rf, ok := ret.Get(0).(func(context.Context) []types.Destination); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.Destination)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTnewMockDestinationRepo interface {
	mock.TestingT
	Cleanup(func())
}

// newMockDestinationRepo creates a new instance of mockDestinationRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockDestinationRepo(t mockConstructorTestingTnewMockDestinationRepo) *mockDestinationRepo {
	mock := &mockDestinationRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
