// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package autoscaling

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/service"
	"github.com/aws/aws-sdk-go/aws/service/serviceinfo"
	"github.com/aws/aws-sdk-go/internal/protocol/query"
	"github.com/aws/aws-sdk-go/internal/signer/v4"
)

// Auto Scaling is designed to automatically launch or terminate EC2 instances
// based on user-defined policies, schedules, and health checks. Use this service
// in conjunction with the Amazon CloudWatch and Elastic Load Balancing services.
type AutoScaling struct {
	*service.Service
}

// Used for custom service initialization logic
var initService func(*service.Service)

// Used for custom request initialization logic
var initRequest func(*request.Request)

// New returns a new AutoScaling client.
func New(config *aws.Config) *AutoScaling {
	service := &service.Service{
		ServiceInfo: serviceinfo.ServiceInfo{
			Config:      defaults.DefaultConfig.Merge(config),
			ServiceName: "autoscaling",
			APIVersion:  "2011-01-01",
		},
	}
	service.Initialize()

	// Handlers
	service.Handlers.Sign.PushBack(v4.Sign)
	service.Handlers.Build.PushBack(query.Build)
	service.Handlers.Unmarshal.PushBack(query.Unmarshal)
	service.Handlers.UnmarshalMeta.PushBack(query.UnmarshalMeta)
	service.Handlers.UnmarshalError.PushBack(query.UnmarshalError)

	// Run custom service initialization if present
	if initService != nil {
		initService(service)
	}

	return &AutoScaling{service}
}

// newRequest creates a new request for a AutoScaling operation and runs any
// custom request initialization.
func (c *AutoScaling) newRequest(op *request.Operation, params, data interface{}) *request.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}
