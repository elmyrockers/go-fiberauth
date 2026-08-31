package xvelope


type HttpContext interface {
	SetHttpContext( ctx any )
}

type FastHttpContext struct {
	context any
}

func (f *FastHttpContext) SetHttpContext( ctx any ){
	f.context = ctx
}

type NetHttpContext struct {
	context any
}

func (n *NetHttpContext) SetHttpContext( ctx any ){
	n.context = ctx
}