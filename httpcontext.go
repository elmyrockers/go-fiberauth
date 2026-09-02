package xvelope

import (
	"encoding/json"
	"time"

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
	Cookie(name string) (string, error) // "auth_payload" cookie
	Query(key string) string            // "?useCookies=true" - detect cookie or token based

	// --- Writing the response ---
	SetStatus(code int)
	SetCookie(cookie *Cookie)
	SendJSON(v any) error

	SetContext( ctx any )
}

// ------------------------------------------ FastHttpContext
type FastHttpContext struct {
	context *fasthttp.RequestCtx
}

func (f *FastHttpContext) SetContext( ctx any ){
	f.context = ctx.(*fasthttp.RequestCtx)
}

func (f *FastHttpContext) Header(key string) string {
	if f.context == nil { return "" }
	return string(f.context.Request.Header.Peek(key))
}

func (f *FastHttpContext) Cookie(name string) (string, error) {
	if f.context == nil { return "", fasthttp.ErrNoCookieValue }

	val := f.context.Request.Header.Cookie(name)
	if len(val) == 0 { return "", fasthttp.ErrNoCookieValue }

	return string(val), nil
}

func (f *FastHttpContext) Query(key string) string {
	if f.context == nil { return "" }
	return string(f.context.QueryArgs().Peek(key))
}

func (f *FastHttpContext) SetStatus(code int) {
	if f.context != nil {
		f.context.SetStatusCode(code)
	}
}

func (f *FastHttpContext) SetCookie(cookie *Cookie) {
	if f.context == nil || cookie == nil { return }

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
	if f.context == nil { return nil }
	
	f.context.SetContentType("application/json")
	return json.NewEncoder(f.context).Encode(v)
}

// ------------------------------------------ NetHttpContext
type NetHttpContext struct {
	context any
}

func (n *NetHttpContext) SetContext( ctx any ){
	n.context = ctx
}