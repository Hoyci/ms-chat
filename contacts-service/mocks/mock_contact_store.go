package mocks

import (
	"context"
	"reflect"

	"github.com/golang/mock/gomock"
	"github.com/hoyci/ms-chat/contacts-service/types"
)

type ContactStore struct {
	ctrl     *gomock.Controller
	recorder *ContactStoreMockRecorder
}

type ContactStoreMockRecorder struct {
	mock *ContactStore
}

func NewMockContactStore(ctrl *gomock.Controller) *ContactStore {
	mock := &ContactStore{ctrl: ctrl}
	mock.recorder = &ContactStoreMockRecorder{mock}
	return mock
}

func (m *ContactStore) EXPECT() *ContactStoreMockRecorder {
	return m.recorder
}

func (m *ContactStore) CreateContact(arg0 context.Context, arg1, arg2 string) (*types.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateContact", arg0, arg1, arg2)
	ret0, _ := ret[0].(*types.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *ContactStoreMockRecorder) CreateContact(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(
		mr.mock, "CreateContact", reflect.TypeOf((*ContactStore)(nil).CreateContact), arg0, arg1, arg2,
	)
}

func (m *ContactStore) GetAllContactsByOwnerID(arg0 context.Context, arg1 string) ([]*types.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllContactsByOwnerID", arg0, arg1)
	ret0, _ := ret[0].([]*types.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *ContactStoreMockRecorder) GetAllContactsByOwnerID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(
		mr.mock, "GetAllContactsByOwnerID", reflect.TypeOf((*ContactStore)(nil).GetAllContactsByOwnerID), arg0, arg1,
	)
}

func (m *ContactStore) GetContactByOwnerID(arg0 context.Context, arg1, arg2 string) (*types.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContactByOwnerID", arg0, arg1, arg2)
	ret0, _ := ret[0].(*types.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *ContactStoreMockRecorder) GetContactByOwnerID(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(
		mr.mock, "GetContactByOwnerID", reflect.TypeOf((*ContactStore)(nil).GetContactByOwnerID), arg0, arg1, arg2,
	)
}
