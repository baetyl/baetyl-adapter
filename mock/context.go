// Code generated by MockGen. DO NOT EDIT.
// Source: ./dmcontext/context.go

// Package dm is a generated GoMock package.
package dm

import (
	os "os"
	reflect "reflect"

	context "github.com/baetyl/baetyl-go/v2/context"
	dmcontext "github.com/baetyl/baetyl-go/v2/dmcontext"
	http "github.com/baetyl/baetyl-go/v2/http"
	log "github.com/baetyl/baetyl-go/v2/log"
	mqtt "github.com/baetyl/baetyl-go/v2/mqtt"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	gomock "github.com/golang/mock/gomock"
)

// MockContext is a mock of Context interface.
type MockContext struct {
	ctrl     *gomock.Controller
	recorder *MockContextMockRecorder
}

// MockContextMockRecorder is the mock recorder for MockContext.
type MockContextMockRecorder struct {
	mock *MockContext
}

// NewMockContext creates a new mock instance.
func NewMockContext(ctrl *gomock.Controller) *MockContext {
	mock := &MockContext{ctrl: ctrl}
	mock.recorder = &MockContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContext) EXPECT() *MockContextMockRecorder {
	return m.recorder
}

// AppName mocks base method.
func (m *MockContext) AppName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppName")
	ret0, _ := ret[0].(string)
	return ret0
}

// AppName indicates an expected call of AppName.
func (mr *MockContextMockRecorder) AppName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppName", reflect.TypeOf((*MockContext)(nil).AppName))
}

// AppVersion mocks base method.
func (m *MockContext) AppVersion() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppVersion")
	ret0, _ := ret[0].(string)
	return ret0
}

// AppVersion indicates an expected call of AppVersion.
func (mr *MockContextMockRecorder) AppVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppVersion", reflect.TypeOf((*MockContext)(nil).AppVersion))
}

// CheckSystemCert mocks base method.
func (m *MockContext) CheckSystemCert() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckSystemCert")
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckSystemCert indicates an expected call of CheckSystemCert.
func (mr *MockContextMockRecorder) CheckSystemCert() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckSystemCert", reflect.TypeOf((*MockContext)(nil).CheckSystemCert))
}

// Close mocks base method.
func (m *MockContext) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockContextMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockContext)(nil).Close))
}

// ConfFile mocks base method.
func (m *MockContext) ConfFile() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConfFile")
	ret0, _ := ret[0].(string)
	return ret0
}

// ConfFile indicates an expected call of ConfFile.
func (mr *MockContextMockRecorder) ConfFile() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfFile", reflect.TypeOf((*MockContext)(nil).ConfFile))
}

// Delete mocks base method.
func (m *MockContext) Delete(key interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Delete", key)
}

// Delete indicates an expected call of Delete.
func (mr *MockContextMockRecorder) Delete(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockContext)(nil).Delete), key)
}

// GetAccessConfig mocks base method.
func (m *MockContext) GetAccessConfig() map[string]dmcontext.AccessConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccessConfig")
	ret0, _ := ret[0].(map[string]dmcontext.AccessConfig)
	return ret0
}

// GetAccessConfig indicates an expected call of GetAccessConfig.
func (mr *MockContextMockRecorder) GetAccessConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccessConfig", reflect.TypeOf((*MockContext)(nil).GetAccessConfig))
}

// GetAccessTemplates mocks base method.
func (m *MockContext) GetAccessTemplates(device *dmcontext.DeviceInfo) (*dmcontext.AccessTemplate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccessTemplates", device)
	ret0, _ := ret[0].(*dmcontext.AccessTemplate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccessTemplates indicates an expected call of GetAccessTemplates.
func (mr *MockContextMockRecorder) GetAccessTemplates(device interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccessTemplates", reflect.TypeOf((*MockContext)(nil).GetAccessTemplates), device)
}

// GetAllAccessTemplates mocks base method.
func (m *MockContext) GetAllAccessTemplates() map[string]dmcontext.AccessTemplate {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllAccessTemplates")
	ret0, _ := ret[0].(map[string]dmcontext.AccessTemplate)
	return ret0
}

// GetAllAccessTemplates indicates an expected call of GetAllAccessTemplates.
func (mr *MockContextMockRecorder) GetAllAccessTemplates() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllAccessTemplates", reflect.TypeOf((*MockContext)(nil).GetAllAccessTemplates))
}

// GetAllDeviceModels mocks base method.
func (m *MockContext) GetAllDeviceModels() map[string][]dmcontext.DeviceProperty {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllDeviceModels")
	ret0, _ := ret[0].(map[string][]dmcontext.DeviceProperty)
	return ret0
}

// GetAllDeviceModels indicates an expected call of GetAllDeviceModels.
func (mr *MockContextMockRecorder) GetAllDeviceModels() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllDeviceModels", reflect.TypeOf((*MockContext)(nil).GetAllDeviceModels))
}

// GetAllDevices mocks base method.
func (m *MockContext) GetAllDevices() []dmcontext.DeviceInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllDevices")
	ret0, _ := ret[0].([]dmcontext.DeviceInfo)
	return ret0
}

// GetAllDevices indicates an expected call of GetAllDevices.
func (mr *MockContextMockRecorder) GetAllDevices() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllDevices", reflect.TypeOf((*MockContext)(nil).GetAllDevices))
}

// GetDevice mocks base method.
func (m *MockContext) GetDevice(device string) (*dmcontext.DeviceInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDevice", device)
	ret0, _ := ret[0].(*dmcontext.DeviceInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDevice indicates an expected call of GetDevice.
func (mr *MockContextMockRecorder) GetDevice(device interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDevice", reflect.TypeOf((*MockContext)(nil).GetDevice), device)
}

// GetDeviceAccessConfig mocks base method.
func (m *MockContext) GetDeviceAccessConfig(device *dmcontext.DeviceInfo) (*dmcontext.AccessConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceAccessConfig", device)
	ret0, _ := ret[0].(*dmcontext.AccessConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceAccessConfig indicates an expected call of GetDeviceAccessConfig.
func (mr *MockContextMockRecorder) GetDeviceAccessConfig(device interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceAccessConfig", reflect.TypeOf((*MockContext)(nil).GetDeviceAccessConfig), device)
}

// GetDeviceModel mocks base method.
func (m *MockContext) GetDeviceModel(device *dmcontext.DeviceInfo) ([]dmcontext.DeviceProperty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceModel", device)
	ret0, _ := ret[0].([]dmcontext.DeviceProperty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceModel indicates an expected call of GetDeviceModel.
func (mr *MockContextMockRecorder) GetDeviceModel(device interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceModel", reflect.TypeOf((*MockContext)(nil).GetDeviceModel), device)
}

// GetDeviceProperties mocks base method.
func (m *MockContext) GetDeviceProperties(device *dmcontext.DeviceInfo) (*dmcontext.DeviceShadow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceProperties", device)
	ret0, _ := ret[0].(*dmcontext.DeviceShadow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceProperties indicates an expected call of GetDeviceProperties.
func (mr *MockContextMockRecorder) GetDeviceProperties(device interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceProperties", reflect.TypeOf((*MockContext)(nil).GetDeviceProperties), device)
}

// GetDevicePropertiesConfig mocks base method.
func (m *MockContext) GetDevicePropertiesConfig(device *dmcontext.DeviceInfo) ([]dmcontext.DeviceProperty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDevicePropertiesConfig", device)
	ret0, _ := ret[0].([]dmcontext.DeviceProperty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDevicePropertiesConfig indicates an expected call of GetDevicePropertiesConfig.
func (mr *MockContextMockRecorder) GetDevicePropertiesConfig(device interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDevicePropertiesConfig", reflect.TypeOf((*MockContext)(nil).GetDevicePropertiesConfig), device)
}

// GetDriverConfig mocks base method.
func (m *MockContext) GetDriverConfig() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDriverConfig")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetDriverConfig indicates an expected call of GetDriverConfig.
func (mr *MockContextMockRecorder) GetDriverConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDriverConfig", reflect.TypeOf((*MockContext)(nil).GetDriverConfig))
}

// GetPropertiesConfig mocks base method.
func (m *MockContext) GetPropertiesConfig() map[string][]dmcontext.DeviceProperty {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPropertiesConfig")
	ret0, _ := ret[0].(map[string][]dmcontext.DeviceProperty)
	return ret0
}

// GetPropertiesConfig indicates an expected call of GetPropertiesConfig.
func (mr *MockContextMockRecorder) GetPropertiesConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPropertiesConfig", reflect.TypeOf((*MockContext)(nil).GetPropertiesConfig))
}

// Load mocks base method.
func (m *MockContext) Load(key interface{}) (interface{}, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", key)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Load indicates an expected call of Load.
func (mr *MockContextMockRecorder) Load(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockContext)(nil).Load), key)
}

// LoadCustomConfig mocks base method.
func (m *MockContext) LoadCustomConfig(cfg interface{}, files ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{cfg}
	for _, a := range files {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "LoadCustomConfig", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// LoadCustomConfig indicates an expected call of LoadCustomConfig.
func (mr *MockContextMockRecorder) LoadCustomConfig(cfg interface{}, files ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{cfg}, files...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadCustomConfig", reflect.TypeOf((*MockContext)(nil).LoadCustomConfig), varargs...)
}

// LoadOrStore mocks base method.
func (m *MockContext) LoadOrStore(key, value interface{}) (interface{}, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadOrStore", key, value)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// LoadOrStore indicates an expected call of LoadOrStore.
func (mr *MockContextMockRecorder) LoadOrStore(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadOrStore", reflect.TypeOf((*MockContext)(nil).LoadOrStore), key, value)
}

// Log mocks base method.
func (m *MockContext) Log() *log.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Log")
	ret0, _ := ret[0].(*log.Logger)
	return ret0
}

// Log indicates an expected call of Log.
func (mr *MockContextMockRecorder) Log() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Log", reflect.TypeOf((*MockContext)(nil).Log))
}

// NewBrokerClient mocks base method.
func (m *MockContext) NewBrokerClient(arg0 mqtt.ClientConfig) (*mqtt.Client, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewBrokerClient", arg0)
	ret0, _ := ret[0].(*mqtt.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewBrokerClient indicates an expected call of NewBrokerClient.
func (mr *MockContextMockRecorder) NewBrokerClient(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewBrokerClient", reflect.TypeOf((*MockContext)(nil).NewBrokerClient), arg0)
}

// NewCoreHttpClient mocks base method.
func (m *MockContext) NewCoreHttpClient() (*http.Client, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewCoreHttpClient")
	ret0, _ := ret[0].(*http.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewCoreHttpClient indicates an expected call of NewCoreHttpClient.
func (mr *MockContextMockRecorder) NewCoreHttpClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewCoreHttpClient", reflect.TypeOf((*MockContext)(nil).NewCoreHttpClient))
}

// NewFunctionHttpClient mocks base method.
func (m *MockContext) NewFunctionHttpClient() (*http.Client, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewFunctionHttpClient")
	ret0, _ := ret[0].(*http.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewFunctionHttpClient indicates an expected call of NewFunctionHttpClient.
func (mr *MockContextMockRecorder) NewFunctionHttpClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewFunctionHttpClient", reflect.TypeOf((*MockContext)(nil).NewFunctionHttpClient))
}

// NewSystemBrokerClient mocks base method.
func (m *MockContext) NewSystemBrokerClient(arg0 []mqtt.QOSTopic) (*mqtt.Client, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewSystemBrokerClient", arg0)
	ret0, _ := ret[0].(*mqtt.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewSystemBrokerClient indicates an expected call of NewSystemBrokerClient.
func (mr *MockContextMockRecorder) NewSystemBrokerClient(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewSystemBrokerClient", reflect.TypeOf((*MockContext)(nil).NewSystemBrokerClient), arg0)
}

// NewSystemBrokerClientConfig mocks base method.
func (m *MockContext) NewSystemBrokerClientConfig() (mqtt.ClientConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewSystemBrokerClientConfig")
	ret0, _ := ret[0].(mqtt.ClientConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewSystemBrokerClientConfig indicates an expected call of NewSystemBrokerClientConfig.
func (mr *MockContextMockRecorder) NewSystemBrokerClientConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewSystemBrokerClientConfig", reflect.TypeOf((*MockContext)(nil).NewSystemBrokerClientConfig))
}

// NodeName mocks base method.
func (m *MockContext) NodeName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NodeName")
	ret0, _ := ret[0].(string)
	return ret0
}

// NodeName indicates an expected call of NodeName.
func (mr *MockContextMockRecorder) NodeName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NodeName", reflect.TypeOf((*MockContext)(nil).NodeName))
}

// Offline mocks base method.
func (m *MockContext) Offline(device *dmcontext.DeviceInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Offline", device)
	ret0, _ := ret[0].(error)
	return ret0
}

// Offline indicates an expected call of Offline.
func (mr *MockContextMockRecorder) Offline(device interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Offline", reflect.TypeOf((*MockContext)(nil).Offline), device)
}

// Online mocks base method.
func (m *MockContext) Online(device *dmcontext.DeviceInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Online", device)
	ret0, _ := ret[0].(error)
	return ret0
}

// Online indicates an expected call of Online.
func (mr *MockContextMockRecorder) Online(device interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Online", reflect.TypeOf((*MockContext)(nil).Online), device)
}

// RegisterDeltaCallback mocks base method.
func (m *MockContext) RegisterDeltaCallback(cb dmcontext.DeltaCallback) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterDeltaCallback", cb)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterDeltaCallback indicates an expected call of RegisterDeltaCallback.
func (mr *MockContextMockRecorder) RegisterDeltaCallback(cb interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterDeltaCallback", reflect.TypeOf((*MockContext)(nil).RegisterDeltaCallback), cb)
}

// RegisterEventCallback mocks base method.
func (m *MockContext) RegisterEventCallback(cb dmcontext.EventCallback) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterEventCallback", cb)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterEventCallback indicates an expected call of RegisterEventCallback.
func (mr *MockContextMockRecorder) RegisterEventCallback(cb interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterEventCallback", reflect.TypeOf((*MockContext)(nil).RegisterEventCallback), cb)
}

// ReportDeviceProperties mocks base method.
func (m *MockContext) ReportDeviceProperties(arg0 *dmcontext.DeviceInfo, arg1 v1.Report) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReportDeviceProperties", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReportDeviceProperties indicates an expected call of ReportDeviceProperties.
func (mr *MockContextMockRecorder) ReportDeviceProperties(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReportDeviceProperties", reflect.TypeOf((*MockContext)(nil).ReportDeviceProperties), arg0, arg1)
}

// ReportDevicePropertiesWithFilter mocks base method.
func (m *MockContext) ReportDevicePropertiesWithFilter(arg0 *dmcontext.DeviceInfo, arg1 v1.Report) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReportDevicePropertiesWithFilter", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReportDevicePropertiesWithFilter indicates an expected call of ReportDevicePropertiesWithFilter.
func (mr *MockContextMockRecorder) ReportDevicePropertiesWithFilter(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReportDevicePropertiesWithFilter", reflect.TypeOf((*MockContext)(nil).ReportDevicePropertiesWithFilter), arg0, arg1)
}

// ServiceName mocks base method.
func (m *MockContext) ServiceName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ServiceName")
	ret0, _ := ret[0].(string)
	return ret0
}

// ServiceName indicates an expected call of ServiceName.
func (mr *MockContextMockRecorder) ServiceName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServiceName", reflect.TypeOf((*MockContext)(nil).ServiceName))
}

// Start mocks base method.
func (m *MockContext) Start() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start")
}

// Start indicates an expected call of Start.
func (mr *MockContextMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockContext)(nil).Start))
}

// Store mocks base method.
func (m *MockContext) Store(key, value interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Store", key, value)
}

// Store indicates an expected call of Store.
func (mr *MockContextMockRecorder) Store(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockContext)(nil).Store), key, value)
}

// SystemConfig mocks base method.
func (m *MockContext) SystemConfig() *context.SystemConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SystemConfig")
	ret0, _ := ret[0].(*context.SystemConfig)
	return ret0
}

// SystemConfig indicates an expected call of SystemConfig.
func (mr *MockContextMockRecorder) SystemConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SystemConfig", reflect.TypeOf((*MockContext)(nil).SystemConfig))
}

// Wait mocks base method.
func (m *MockContext) Wait() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Wait")
}

// Wait indicates an expected call of Wait.
func (mr *MockContextMockRecorder) Wait() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Wait", reflect.TypeOf((*MockContext)(nil).Wait))
}

// WaitChan mocks base method.
func (m *MockContext) WaitChan() <-chan os.Signal {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitChan")
	ret0, _ := ret[0].(<-chan os.Signal)
	return ret0
}

// WaitChan indicates an expected call of WaitChan.
func (mr *MockContextMockRecorder) WaitChan() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitChan", reflect.TypeOf((*MockContext)(nil).WaitChan))
}

// RegisterPropertyGetCallback mocks base method.
func (m *MockContext) RegisterPropertyGetCallback(arg0 dmcontext.PropertyGetCallback) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterPropertyGetCallback", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterPropertyGetCallback indicates an expected call of RegisterPropertyGetCallback.
func (mr *MockContextMockRecorder) RegisterPropertyGetCallback(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterPropertyGetCallback", reflect.TypeOf((*MockContext)(nil).RegisterPropertyGetCallback), arg0)
}

// ReportDeviceEvents mocks base method.
func (m *MockContext) ReportDeviceEvents(arg0 *dmcontext.DeviceInfo, arg1 v1.EventReport) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReportDeviceEvents", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReportDeviceEvents indicates an expected call of ReportDeviceEvents.
func (mr *MockContextMockRecorder) ReportDeviceEvents(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReportDeviceEvents", reflect.TypeOf((*MockContext)(nil).ReportDeviceEvents), arg0, arg1)
}
