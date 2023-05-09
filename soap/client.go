package soap

import "context"

type SOAPClient interface {
	Call(ctx context.Context, soapAction string, input, output, faultDetail any) error
}
