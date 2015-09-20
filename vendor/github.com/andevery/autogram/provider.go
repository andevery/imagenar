package autogram

import (
	"errors"
	"github.com/andevery/instaw"
	"github.com/andevery/instax"
)

type Source int

const (
	TAGS Source = iota
	MEDIA
	FOLLOWERS
	NONE
)

var (
	UknownSource   = errors.New("Unknown provider source")
	NotImplemented = errors.New("Source not implemented")
	EOF            = errors.New("End of feed")
)

type Provider struct {
	Limiter *Limiter

	source  Source
	queries []string

	state struct {
		queryNum  int
		step      int
		subStep   int
		mediaFeed *instax.MediaFeed
		media     []instax.Media
	}
}

func NewProvider(src Source, q []string, l *Limiter) (*Provider, error) {
	p := &Provider{
		source:  src,
		queries: q,
		Limiter: l,
	}

	switch src {
	case TAGS, MEDIA:
		p.setNextMediaFeed()
	case FOLLOWERS:
		return nil, NotImplemented
	case NONE:
	default:
		return nil, UknownSource
	}

	return p, nil
}

func EmptyProvider(l *Limiter) (*Provider, error) {
	return NewProvider(NONE, []string{}, l)
}

func (p *Provider) ApiClient() *instax.Client {
	return p.Limiter.ApiClient()
}

func (p *Provider) WebClient() *instaw.Client {
	return p.Limiter.WebClient()
}

func (p *Provider) TotalAmount() uint32 {
	return p.Limiter.TotalAmount()
}

func (p *Provider) setNextMediaFeed() (err error) {
	if p.state.queryNum >= len(p.queries) {
		err = EOF
		return
	}

	switch p.source {
	case TAGS:
		p.state.mediaFeed = p.ApiClient().MediaByTag(p.queries[p.state.queryNum])
	case MEDIA:
		p.state.mediaFeed = p.ApiClient().MediaByUser(p.queries[p.state.queryNum])
	}

	if err != nil {
		return
	}

	p.state.queryNum++
	p.state.step = 0
	return
}

func (p *Provider) NextUsers() ([]instax.UserShort, error) {
	var err error
	if p.state.subStep >= len(p.state.media) {
		p.state.subStep = 0
		p.state.media, err = p.NextMedia()
		if err != nil {
			return nil, err
		}
	}

	users, err := p.ApiClient().Likes(p.state.media[p.state.subStep].ID)
	if err != nil {
		return nil, err
	}
	p.state.subStep++

	return users, nil
}

func (p *Provider) NextMedia() ([]instax.Media, error) {
	var media []instax.Media
	var err error

	if p.state.step == 0 {
		media, err = p.state.mediaFeed.Get()
	} else if p.state.mediaFeed.CanNext() {
		media, err = p.state.mediaFeed.Next()
	} else {
		if err = p.setNextMediaFeed(); err == nil {
			media, err = p.state.mediaFeed.Get()
		}
	}

	if err != nil {
		return nil, err
	}

	p.state.step++
	return media, err
}
