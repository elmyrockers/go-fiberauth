package xvelope

// import "github.com/davecgh/go-spew/spew"

type HttpContext interface {
	SetContext( ctx any )
}

type FastHttpContext struct {
	context any
}

func (f *FastHttpContext) SetContext( ctx any ){
	// spew.Dump( ctx )
	f.context = ctx
}

type NetHttpContext struct {
	context any
}

func (n *NetHttpContext) SetContext( ctx any ){
	n.context = ctx
}