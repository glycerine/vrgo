package primary

import (
	"strconv"

	"github.com/BoolLi/vrgo/globals"
	vrrpc "github.com/BoolLi/vrgo/rpc"
)

// VrgoRPC defines the user RPCs exported by server.
type VrgoRPC int

func (v *VrgoRPC) Execute(req *vrrpc.Request, resp *vrrpc.Response) error {
	// If mode is not primary, then tell client to retry or who the new primary is.
	if globals.Mode != "primary" {
		if globals.Mode == "viewchange" || globals.Mode == "viewchange-init" {
			globals.Log("Execute", "currently under view change. Try again later")
		} else {
			globals.Log("Execute", "I am not primary anymore; view num: %v", globals.ViewNum)
		}
		*resp = vrrpc.Response{
			ViewNum: globals.ViewNum,
		}
		return nil
	}

	k := strconv.Itoa(req.ClientId)
	res, ok := globals.ClientTable.Get(k)

	// If the client request is already executed before, resend the response.
	if ok && req.RequestNum <= res.(vrrpc.Response).RequestNum {
		globals.Log("Execute", "request %+v is already executed; returning previous result %+v directly", req, res)
		*resp = res.(vrrpc.Response)
		return nil
	}

	// First time receiving from this client.
	if !ok {
		globals.Log("Execute", "first time receiving request %v from client %v\n", req.RequestNum, req.ClientId)
	}

	ch := AddIncomingReq(req)
	select {
	case res := <-ch:
		globals.Log("Execute", "done processing request; got result %v\n", res.OpResult.Message)
		*resp = *res
	}

	return nil
}
