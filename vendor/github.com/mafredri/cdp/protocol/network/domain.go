// Code generated by cdpgen. DO NOT EDIT.

// Package network implements the Network domain. Network domain allows tracking network activities of the page. It exposes information about http, file, data and other requests and responses, their headers, bodies, timing, etc.
package network

import (
	"context"

	"github.com/mafredri/cdp/protocol/internal"
	"github.com/mafredri/cdp/rpcc"
)

// domainClient is a client for the Network domain. Network domain allows tracking network activities of the page. It exposes information about http, file, data and other requests and responses, their headers, bodies, timing, etc.
type domainClient struct{ conn *rpcc.Conn }

// NewClient returns a client for the Network domain with the connection set to conn.
func NewClient(conn *rpcc.Conn) *domainClient {
	return &domainClient{conn: conn}
}

// Enable invokes the Network method. Enables network tracking, network events will now be delivered to the client.
func (d *domainClient) Enable(ctx context.Context, args *EnableArgs) (err error) {
	if args != nil {
		err = rpcc.Invoke(ctx, "Network.enable", args, nil, d.conn)
	} else {
		err = rpcc.Invoke(ctx, "Network.enable", nil, nil, d.conn)
	}
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "Enable", Err: err}
	}
	return
}

// Disable invokes the Network method. Disables network tracking, prevents network events from being sent to the client.
func (d *domainClient) Disable(ctx context.Context) (err error) {
	err = rpcc.Invoke(ctx, "Network.disable", nil, nil, d.conn)
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "Disable", Err: err}
	}
	return
}

// SetUserAgentOverride invokes the Network method. Allows overriding user agent with the given string.
func (d *domainClient) SetUserAgentOverride(ctx context.Context, args *SetUserAgentOverrideArgs) (err error) {
	if args != nil {
		err = rpcc.Invoke(ctx, "Network.setUserAgentOverride", args, nil, d.conn)
	} else {
		err = rpcc.Invoke(ctx, "Network.setUserAgentOverride", nil, nil, d.conn)
	}
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "SetUserAgentOverride", Err: err}
	}
	return
}

// SetExtraHTTPHeaders invokes the Network method. Specifies whether to always send extra HTTP headers with the requests from this page.
func (d *domainClient) SetExtraHTTPHeaders(ctx context.Context, args *SetExtraHTTPHeadersArgs) (err error) {
	if args != nil {
		err = rpcc.Invoke(ctx, "Network.setExtraHTTPHeaders", args, nil, d.conn)
	} else {
		err = rpcc.Invoke(ctx, "Network.setExtraHTTPHeaders", nil, nil, d.conn)
	}
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "SetExtraHTTPHeaders", Err: err}
	}
	return
}

// GetResponseBody invokes the Network method. Returns content served for the given request.
func (d *domainClient) GetResponseBody(ctx context.Context, args *GetResponseBodyArgs) (reply *GetResponseBodyReply, err error) {
	reply = new(GetResponseBodyReply)
	if args != nil {
		err = rpcc.Invoke(ctx, "Network.getResponseBody", args, reply, d.conn)
	} else {
		err = rpcc.Invoke(ctx, "Network.getResponseBody", nil, reply, d.conn)
	}
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "GetResponseBody", Err: err}
	}
	return
}

// SetBlockedURLs invokes the Network method. Blocks URLs from loading.
func (d *domainClient) SetBlockedURLs(ctx context.Context, args *SetBlockedURLsArgs) (err error) {
	if args != nil {
		err = rpcc.Invoke(ctx, "Network.setBlockedURLs", args, nil, d.conn)
	} else {
		err = rpcc.Invoke(ctx, "Network.setBlockedURLs", nil, nil, d.conn)
	}
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "SetBlockedURLs", Err: err}
	}
	return
}

// ReplayXHR invokes the Network method. This method sends a new XMLHttpRequest which is identical to the original one. The following parameters should be identical: method, url, async, request body, extra headers, withCredentials attribute, user, password.
func (d *domainClient) ReplayXHR(ctx context.Context, args *ReplayXHRArgs) (err error) {
	if args != nil {
		err = rpcc.Invoke(ctx, "Network.replayXHR", args, nil, d.conn)
	} else {
		err = rpcc.Invoke(ctx, "Network.replayXHR", nil, nil, d.conn)
	}
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "ReplayXHR", Err: err}
	}
	return
}

// CanClearBrowserCache invokes the Network method. Tells whether clearing browser cache is supported.
func (d *domainClient) CanClearBrowserCache(ctx context.Context) (reply *CanClearBrowserCacheReply, err error) {
	reply = new(CanClearBrowserCacheReply)
	err = rpcc.Invoke(ctx, "Network.canClearBrowserCache", nil, reply, d.conn)
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "CanClearBrowserCache", Err: err}
	}
	return
}

// ClearBrowserCache invokes the Network method. Clears browser cache.
func (d *domainClient) ClearBrowserCache(ctx context.Context) (err error) {
	err = rpcc.Invoke(ctx, "Network.clearBrowserCache", nil, nil, d.conn)
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "ClearBrowserCache", Err: err}
	}
	return
}

// CanClearBrowserCookies invokes the Network method. Tells whether clearing browser cookies is supported.
func (d *domainClient) CanClearBrowserCookies(ctx context.Context) (reply *CanClearBrowserCookiesReply, err error) {
	reply = new(CanClearBrowserCookiesReply)
	err = rpcc.Invoke(ctx, "Network.canClearBrowserCookies", nil, reply, d.conn)
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "CanClearBrowserCookies", Err: err}
	}
	return
}

// ClearBrowserCookies invokes the Network method. Clears browser cookies.
func (d *domainClient) ClearBrowserCookies(ctx context.Context) (err error) {
	err = rpcc.Invoke(ctx, "Network.clearBrowserCookies", nil, nil, d.conn)
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "ClearBrowserCookies", Err: err}
	}
	return
}

// GetCookies invokes the Network method. Returns all browser cookies for the current URL. Depending on the backend support, will return detailed cookie information in the cookies field.
func (d *domainClient) GetCookies(ctx context.Context, args *GetCookiesArgs) (reply *GetCookiesReply, err error) {
	reply = new(GetCookiesReply)
	if args != nil {
		err = rpcc.Invoke(ctx, "Network.getCookies", args, reply, d.conn)
	} else {
		err = rpcc.Invoke(ctx, "Network.getCookies", nil, reply, d.conn)
	}
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "GetCookies", Err: err}
	}
	return
}

// GetAllCookies invokes the Network method. Returns all browser cookies. Depending on the backend support, will return detailed cookie information in the cookies field.
func (d *domainClient) GetAllCookies(ctx context.Context) (reply *GetAllCookiesReply, err error) {
	reply = new(GetAllCookiesReply)
	err = rpcc.Invoke(ctx, "Network.getAllCookies", nil, reply, d.conn)
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "GetAllCookies", Err: err}
	}
	return
}

// DeleteCookie invokes the Network method. Deletes browser cookie with given name, domain and path.
func (d *domainClient) DeleteCookie(ctx context.Context, args *DeleteCookieArgs) (err error) {
	if args != nil {
		err = rpcc.Invoke(ctx, "Network.deleteCookie", args, nil, d.conn)
	} else {
		err = rpcc.Invoke(ctx, "Network.deleteCookie", nil, nil, d.conn)
	}
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "DeleteCookie", Err: err}
	}
	return
}

// SetCookie invokes the Network method. Sets a cookie with the given cookie data; may overwrite equivalent cookies if they exist.
func (d *domainClient) SetCookie(ctx context.Context, args *SetCookieArgs) (reply *SetCookieReply, err error) {
	reply = new(SetCookieReply)
	if args != nil {
		err = rpcc.Invoke(ctx, "Network.setCookie", args, reply, d.conn)
	} else {
		err = rpcc.Invoke(ctx, "Network.setCookie", nil, reply, d.conn)
	}
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "SetCookie", Err: err}
	}
	return
}

// SetCookies invokes the Network method. Sets given cookies.
func (d *domainClient) SetCookies(ctx context.Context, args *SetCookiesArgs) (err error) {
	if args != nil {
		err = rpcc.Invoke(ctx, "Network.setCookies", args, nil, d.conn)
	} else {
		err = rpcc.Invoke(ctx, "Network.setCookies", nil, nil, d.conn)
	}
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "SetCookies", Err: err}
	}
	return
}

// CanEmulateNetworkConditions invokes the Network method. Tells whether emulation of network conditions is supported.
func (d *domainClient) CanEmulateNetworkConditions(ctx context.Context) (reply *CanEmulateNetworkConditionsReply, err error) {
	reply = new(CanEmulateNetworkConditionsReply)
	err = rpcc.Invoke(ctx, "Network.canEmulateNetworkConditions", nil, reply, d.conn)
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "CanEmulateNetworkConditions", Err: err}
	}
	return
}

// EmulateNetworkConditions invokes the Network method. Activates emulation of network conditions.
func (d *domainClient) EmulateNetworkConditions(ctx context.Context, args *EmulateNetworkConditionsArgs) (err error) {
	if args != nil {
		err = rpcc.Invoke(ctx, "Network.emulateNetworkConditions", args, nil, d.conn)
	} else {
		err = rpcc.Invoke(ctx, "Network.emulateNetworkConditions", nil, nil, d.conn)
	}
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "EmulateNetworkConditions", Err: err}
	}
	return
}

// SetCacheDisabled invokes the Network method. Toggles ignoring cache for each request. If true, cache will not be used.
func (d *domainClient) SetCacheDisabled(ctx context.Context, args *SetCacheDisabledArgs) (err error) {
	if args != nil {
		err = rpcc.Invoke(ctx, "Network.setCacheDisabled", args, nil, d.conn)
	} else {
		err = rpcc.Invoke(ctx, "Network.setCacheDisabled", nil, nil, d.conn)
	}
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "SetCacheDisabled", Err: err}
	}
	return
}

// SetBypassServiceWorker invokes the Network method. Toggles ignoring of service worker for each request.
func (d *domainClient) SetBypassServiceWorker(ctx context.Context, args *SetBypassServiceWorkerArgs) (err error) {
	if args != nil {
		err = rpcc.Invoke(ctx, "Network.setBypassServiceWorker", args, nil, d.conn)
	} else {
		err = rpcc.Invoke(ctx, "Network.setBypassServiceWorker", nil, nil, d.conn)
	}
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "SetBypassServiceWorker", Err: err}
	}
	return
}

// SetDataSizeLimitsForTest invokes the Network method. For testing.
func (d *domainClient) SetDataSizeLimitsForTest(ctx context.Context, args *SetDataSizeLimitsForTestArgs) (err error) {
	if args != nil {
		err = rpcc.Invoke(ctx, "Network.setDataSizeLimitsForTest", args, nil, d.conn)
	} else {
		err = rpcc.Invoke(ctx, "Network.setDataSizeLimitsForTest", nil, nil, d.conn)
	}
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "SetDataSizeLimitsForTest", Err: err}
	}
	return
}

// GetCertificate invokes the Network method. Returns the DER-encoded certificate.
func (d *domainClient) GetCertificate(ctx context.Context, args *GetCertificateArgs) (reply *GetCertificateReply, err error) {
	reply = new(GetCertificateReply)
	if args != nil {
		err = rpcc.Invoke(ctx, "Network.getCertificate", args, reply, d.conn)
	} else {
		err = rpcc.Invoke(ctx, "Network.getCertificate", nil, reply, d.conn)
	}
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "GetCertificate", Err: err}
	}
	return
}

// SetRequestInterceptionEnabled invokes the Network method.
func (d *domainClient) SetRequestInterceptionEnabled(ctx context.Context, args *SetRequestInterceptionEnabledArgs) (err error) {
	if args != nil {
		err = rpcc.Invoke(ctx, "Network.setRequestInterceptionEnabled", args, nil, d.conn)
	} else {
		err = rpcc.Invoke(ctx, "Network.setRequestInterceptionEnabled", nil, nil, d.conn)
	}
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "SetRequestInterceptionEnabled", Err: err}
	}
	return
}

// ContinueInterceptedRequest invokes the Network method. Response to Network.requestIntercepted which either modifies the request to continue with any modifications, or blocks it, or completes it with the provided response bytes. If a network fetch occurs as a result which encounters a redirect an additional Network.requestIntercepted event will be sent with the same InterceptionId.
func (d *domainClient) ContinueInterceptedRequest(ctx context.Context, args *ContinueInterceptedRequestArgs) (err error) {
	if args != nil {
		err = rpcc.Invoke(ctx, "Network.continueInterceptedRequest", args, nil, d.conn)
	} else {
		err = rpcc.Invoke(ctx, "Network.continueInterceptedRequest", nil, nil, d.conn)
	}
	if err != nil {
		err = &internal.OpError{Domain: "Network", Op: "ContinueInterceptedRequest", Err: err}
	}
	return
}

func (d *domainClient) ResourceChangedPriority(ctx context.Context) (ResourceChangedPriorityClient, error) {
	s, err := rpcc.NewStream(ctx, "Network.resourceChangedPriority", d.conn)
	if err != nil {
		return nil, err
	}
	return &resourceChangedPriorityClient{Stream: s}, nil
}

type resourceChangedPriorityClient struct{ rpcc.Stream }

func (c *resourceChangedPriorityClient) Recv() (*ResourceChangedPriorityReply, error) {
	event := new(ResourceChangedPriorityReply)
	if err := c.RecvMsg(event); err != nil {
		return nil, &internal.OpError{Domain: "Network", Op: "ResourceChangedPriority Recv", Err: err}
	}
	return event, nil
}

func (d *domainClient) RequestWillBeSent(ctx context.Context) (RequestWillBeSentClient, error) {
	s, err := rpcc.NewStream(ctx, "Network.requestWillBeSent", d.conn)
	if err != nil {
		return nil, err
	}
	return &requestWillBeSentClient{Stream: s}, nil
}

type requestWillBeSentClient struct{ rpcc.Stream }

func (c *requestWillBeSentClient) Recv() (*RequestWillBeSentReply, error) {
	event := new(RequestWillBeSentReply)
	if err := c.RecvMsg(event); err != nil {
		return nil, &internal.OpError{Domain: "Network", Op: "RequestWillBeSent Recv", Err: err}
	}
	return event, nil
}

func (d *domainClient) RequestServedFromCache(ctx context.Context) (RequestServedFromCacheClient, error) {
	s, err := rpcc.NewStream(ctx, "Network.requestServedFromCache", d.conn)
	if err != nil {
		return nil, err
	}
	return &requestServedFromCacheClient{Stream: s}, nil
}

type requestServedFromCacheClient struct{ rpcc.Stream }

func (c *requestServedFromCacheClient) Recv() (*RequestServedFromCacheReply, error) {
	event := new(RequestServedFromCacheReply)
	if err := c.RecvMsg(event); err != nil {
		return nil, &internal.OpError{Domain: "Network", Op: "RequestServedFromCache Recv", Err: err}
	}
	return event, nil
}

func (d *domainClient) ResponseReceived(ctx context.Context) (ResponseReceivedClient, error) {
	s, err := rpcc.NewStream(ctx, "Network.responseReceived", d.conn)
	if err != nil {
		return nil, err
	}
	return &responseReceivedClient{Stream: s}, nil
}

type responseReceivedClient struct{ rpcc.Stream }

func (c *responseReceivedClient) Recv() (*ResponseReceivedReply, error) {
	event := new(ResponseReceivedReply)
	if err := c.RecvMsg(event); err != nil {
		return nil, &internal.OpError{Domain: "Network", Op: "ResponseReceived Recv", Err: err}
	}
	return event, nil
}

func (d *domainClient) DataReceived(ctx context.Context) (DataReceivedClient, error) {
	s, err := rpcc.NewStream(ctx, "Network.dataReceived", d.conn)
	if err != nil {
		return nil, err
	}
	return &dataReceivedClient{Stream: s}, nil
}

type dataReceivedClient struct{ rpcc.Stream }

func (c *dataReceivedClient) Recv() (*DataReceivedReply, error) {
	event := new(DataReceivedReply)
	if err := c.RecvMsg(event); err != nil {
		return nil, &internal.OpError{Domain: "Network", Op: "DataReceived Recv", Err: err}
	}
	return event, nil
}

func (d *domainClient) LoadingFinished(ctx context.Context) (LoadingFinishedClient, error) {
	s, err := rpcc.NewStream(ctx, "Network.loadingFinished", d.conn)
	if err != nil {
		return nil, err
	}
	return &loadingFinishedClient{Stream: s}, nil
}

type loadingFinishedClient struct{ rpcc.Stream }

func (c *loadingFinishedClient) Recv() (*LoadingFinishedReply, error) {
	event := new(LoadingFinishedReply)
	if err := c.RecvMsg(event); err != nil {
		return nil, &internal.OpError{Domain: "Network", Op: "LoadingFinished Recv", Err: err}
	}
	return event, nil
}

func (d *domainClient) LoadingFailed(ctx context.Context) (LoadingFailedClient, error) {
	s, err := rpcc.NewStream(ctx, "Network.loadingFailed", d.conn)
	if err != nil {
		return nil, err
	}
	return &loadingFailedClient{Stream: s}, nil
}

type loadingFailedClient struct{ rpcc.Stream }

func (c *loadingFailedClient) Recv() (*LoadingFailedReply, error) {
	event := new(LoadingFailedReply)
	if err := c.RecvMsg(event); err != nil {
		return nil, &internal.OpError{Domain: "Network", Op: "LoadingFailed Recv", Err: err}
	}
	return event, nil
}

func (d *domainClient) WebSocketWillSendHandshakeRequest(ctx context.Context) (WebSocketWillSendHandshakeRequestClient, error) {
	s, err := rpcc.NewStream(ctx, "Network.webSocketWillSendHandshakeRequest", d.conn)
	if err != nil {
		return nil, err
	}
	return &webSocketWillSendHandshakeRequestClient{Stream: s}, nil
}

type webSocketWillSendHandshakeRequestClient struct{ rpcc.Stream }

func (c *webSocketWillSendHandshakeRequestClient) Recv() (*WebSocketWillSendHandshakeRequestReply, error) {
	event := new(WebSocketWillSendHandshakeRequestReply)
	if err := c.RecvMsg(event); err != nil {
		return nil, &internal.OpError{Domain: "Network", Op: "WebSocketWillSendHandshakeRequest Recv", Err: err}
	}
	return event, nil
}

func (d *domainClient) WebSocketHandshakeResponseReceived(ctx context.Context) (WebSocketHandshakeResponseReceivedClient, error) {
	s, err := rpcc.NewStream(ctx, "Network.webSocketHandshakeResponseReceived", d.conn)
	if err != nil {
		return nil, err
	}
	return &webSocketHandshakeResponseReceivedClient{Stream: s}, nil
}

type webSocketHandshakeResponseReceivedClient struct{ rpcc.Stream }

func (c *webSocketHandshakeResponseReceivedClient) Recv() (*WebSocketHandshakeResponseReceivedReply, error) {
	event := new(WebSocketHandshakeResponseReceivedReply)
	if err := c.RecvMsg(event); err != nil {
		return nil, &internal.OpError{Domain: "Network", Op: "WebSocketHandshakeResponseReceived Recv", Err: err}
	}
	return event, nil
}

func (d *domainClient) WebSocketCreated(ctx context.Context) (WebSocketCreatedClient, error) {
	s, err := rpcc.NewStream(ctx, "Network.webSocketCreated", d.conn)
	if err != nil {
		return nil, err
	}
	return &webSocketCreatedClient{Stream: s}, nil
}

type webSocketCreatedClient struct{ rpcc.Stream }

func (c *webSocketCreatedClient) Recv() (*WebSocketCreatedReply, error) {
	event := new(WebSocketCreatedReply)
	if err := c.RecvMsg(event); err != nil {
		return nil, &internal.OpError{Domain: "Network", Op: "WebSocketCreated Recv", Err: err}
	}
	return event, nil
}

func (d *domainClient) WebSocketClosed(ctx context.Context) (WebSocketClosedClient, error) {
	s, err := rpcc.NewStream(ctx, "Network.webSocketClosed", d.conn)
	if err != nil {
		return nil, err
	}
	return &webSocketClosedClient{Stream: s}, nil
}

type webSocketClosedClient struct{ rpcc.Stream }

func (c *webSocketClosedClient) Recv() (*WebSocketClosedReply, error) {
	event := new(WebSocketClosedReply)
	if err := c.RecvMsg(event); err != nil {
		return nil, &internal.OpError{Domain: "Network", Op: "WebSocketClosed Recv", Err: err}
	}
	return event, nil
}

func (d *domainClient) WebSocketFrameReceived(ctx context.Context) (WebSocketFrameReceivedClient, error) {
	s, err := rpcc.NewStream(ctx, "Network.webSocketFrameReceived", d.conn)
	if err != nil {
		return nil, err
	}
	return &webSocketFrameReceivedClient{Stream: s}, nil
}

type webSocketFrameReceivedClient struct{ rpcc.Stream }

func (c *webSocketFrameReceivedClient) Recv() (*WebSocketFrameReceivedReply, error) {
	event := new(WebSocketFrameReceivedReply)
	if err := c.RecvMsg(event); err != nil {
		return nil, &internal.OpError{Domain: "Network", Op: "WebSocketFrameReceived Recv", Err: err}
	}
	return event, nil
}

func (d *domainClient) WebSocketFrameError(ctx context.Context) (WebSocketFrameErrorClient, error) {
	s, err := rpcc.NewStream(ctx, "Network.webSocketFrameError", d.conn)
	if err != nil {
		return nil, err
	}
	return &webSocketFrameErrorClient{Stream: s}, nil
}

type webSocketFrameErrorClient struct{ rpcc.Stream }

func (c *webSocketFrameErrorClient) Recv() (*WebSocketFrameErrorReply, error) {
	event := new(WebSocketFrameErrorReply)
	if err := c.RecvMsg(event); err != nil {
		return nil, &internal.OpError{Domain: "Network", Op: "WebSocketFrameError Recv", Err: err}
	}
	return event, nil
}

func (d *domainClient) WebSocketFrameSent(ctx context.Context) (WebSocketFrameSentClient, error) {
	s, err := rpcc.NewStream(ctx, "Network.webSocketFrameSent", d.conn)
	if err != nil {
		return nil, err
	}
	return &webSocketFrameSentClient{Stream: s}, nil
}

type webSocketFrameSentClient struct{ rpcc.Stream }

func (c *webSocketFrameSentClient) Recv() (*WebSocketFrameSentReply, error) {
	event := new(WebSocketFrameSentReply)
	if err := c.RecvMsg(event); err != nil {
		return nil, &internal.OpError{Domain: "Network", Op: "WebSocketFrameSent Recv", Err: err}
	}
	return event, nil
}

func (d *domainClient) EventSourceMessageReceived(ctx context.Context) (EventSourceMessageReceivedClient, error) {
	s, err := rpcc.NewStream(ctx, "Network.eventSourceMessageReceived", d.conn)
	if err != nil {
		return nil, err
	}
	return &eventSourceMessageReceivedClient{Stream: s}, nil
}

type eventSourceMessageReceivedClient struct{ rpcc.Stream }

func (c *eventSourceMessageReceivedClient) Recv() (*EventSourceMessageReceivedReply, error) {
	event := new(EventSourceMessageReceivedReply)
	if err := c.RecvMsg(event); err != nil {
		return nil, &internal.OpError{Domain: "Network", Op: "EventSourceMessageReceived Recv", Err: err}
	}
	return event, nil
}

func (d *domainClient) RequestIntercepted(ctx context.Context) (RequestInterceptedClient, error) {
	s, err := rpcc.NewStream(ctx, "Network.requestIntercepted", d.conn)
	if err != nil {
		return nil, err
	}
	return &requestInterceptedClient{Stream: s}, nil
}

type requestInterceptedClient struct{ rpcc.Stream }

func (c *requestInterceptedClient) Recv() (*RequestInterceptedReply, error) {
	event := new(RequestInterceptedReply)
	if err := c.RecvMsg(event); err != nil {
		return nil, &internal.OpError{Domain: "Network", Op: "RequestIntercepted Recv", Err: err}
	}
	return event, nil
}
