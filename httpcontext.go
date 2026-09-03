package xvelope

import (
	"encoding/json"
	"time"
	"net/http"

	"github.com/valyala/fasthttp"
	// "github.com/davecgh/go-spew/spew"
)



type SameSite int

const (
	SameSiteDefault SameSite = iota
	SameSiteLax
	SameSiteStrict
	SameSiteNone
)

type Cookie struct {
	Name     string
	Value    string
	Path     string
	Domain   string
	MaxAge   int
	Expires  time.Time
	Secure   bool
	HTTPOnly bool
	SameSite SameSite
}

type HttpContext interface {
	// --- Reading the incoming credential ---
	Header(key string) string           // "Authorization: Bearer <token>"
	Cookie(name string) string 			// "auth_payload" cookie
	Query(key string) string            // "?useCookies=true" - detect cookie or token based

	// --- Writing the response ---
	SetStatus(code int)
	SetCookie(cookie *Cookie)
	SendJSON(v any) error

	SetContext( ctx any )
}

// Compile-time interface checks
var (
	_ HttpContext = (*FastHttpContext)(nil)
	_ HttpContext = (*NetHttpContext)(nil)
)

// ------------------------------------------ FastHttpContext
type FastHttpContext struct {
	context *fasthttp.RequestCtx
}

func (f *FastHttpContext) SetContext( ctx any ){
	f.context = ctx.(*fasthttp.RequestCtx)
}

func (f *FastHttpContext) Header(key string) string {
	return string(f.context.Request.Header.Peek(key))
}

func (f *FastHttpContext) Cookie(name string) string {
	return string(f.context.Request.Header.Cookie(name))
}

func (f *FastHttpContext) Query(key string) string {
	return string(f.context.QueryArgs().Peek(key))
}

func (f *FastHttpContext) SetStatus(code int) {
	f.context.SetStatusCode(code)
}

func (f *FastHttpContext) SetCookie(cookie *Cookie) {
	if cookie == nil { return }

	// Get cookie instance
		fastCookie := fasthttp.AcquireCookie()
		defer fasthttp.ReleaseCookie(fastCookie)

	// Set cookie
		fastCookie.SetKey(cookie.Name)
		fastCookie.SetValue(cookie.Value)
		fastCookie.SetPath(cookie.Path)
		fastCookie.SetDomain(cookie.Domain)
		fastCookie.SetMaxAge(cookie.MaxAge)
		fastCookie.SetExpire(cookie.Expires)
		fastCookie.SetSecure(cookie.Secure)
		fastCookie.SetHTTPOnly(cookie.HTTPOnly)

		switch cookie.SameSite {
		case SameSiteLax:
			fastCookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)
		case SameSiteStrict:
			fastCookie.SetSameSite(fasthttp.CookieSameSiteStrictMode)
		case SameSiteNone:
			fastCookie.SetSameSite(fasthttp.CookieSameSiteNoneMode)
		default:
			fastCookie.SetSameSite(fasthttp.CookieSameSiteDefaultMode)
		}

		f.context.Response.Header.SetCookie(fastCookie)
}

func (f *FastHttpContext) SendJSON(v any) error {
	f.context.SetContentType("application/json")
	return json.NewEncoder(f.context).Encode(v)
}

// ------------------------------------------ NetHttpContext
type NetHttpContextPair struct {
	Request *http.Request
	Response http.ResponseWriter
}

type NetHttpContext struct {
	context NetHttpContextPair
}

func (n *NetHttpContext) SetContext( ctx any ){
	n.context = ctx.(NetHttpContextPair)
}

func (n *NetHttpContext) Header(key string) string {
	return n.context.Request.Header.Get(key)
}

func (n *NetHttpContext) Cookie(name string) string {
	c, err := n.context.Request.Cookie(name)
	if err != nil { return "" }
	return c.Value
}

func (n *NetHttpContext) Query(key string) string {
	return n.context.Request.URL.Query().Get(key)
}

func (n *NetHttpContext) SetStatus(code int) {
	n.context.Response.WriteHeader(code)
}

func (n *NetHttpContext) SetCookie(cookie *Cookie) {
	if cookie == nil { return }
	
	var sameSite http.SameSite
	switch cookie.SameSite {
	case SameSiteStrict:
		sameSite = http.SameSiteStrictMode
	case SameSiteNone:
		sameSite = http.SameSiteNoneMode
	case SameSiteLax:
		sameSite = http.SameSiteLaxMode
	default:
		sameSite = http.SameSiteDefaultMode
	}

	http.SetCookie(n.context.Response, &http.Cookie{
		Name:     cookie.Name,
		Value:    cookie.Value,
		Path:     cookie.Path,
		Domain:   cookie.Domain,
		MaxAge:   cookie.MaxAge,
		Expires:  cookie.Expires,
		Secure:   cookie.Secure,
		HttpOnly: cookie.HTTPOnly,
		SameSite: sameSite,
	})
}

func (n *NetHttpContext) SendJSON(v any) error {
	n.context.Response.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(n.context.Response).Encode(v)
}