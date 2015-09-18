package autogram

import (
	"errors"
	"github.com/andevery/instax"
)

type Source int

const (
	TAGS Source = iota
	MEDIA
	FOLLOWERS
)

var (
	UknownSource   = errors.New("Unknown source")
	NotImplemented = errors.New("Not implemented")
	EOF            = errors.New("End of feed")
)

type Provider struct {
	source Source
	client *instax.Client

	queries []string

	state struct {
		queryNum  int
		step      int
		subStep   int
		mediaFeed *instax.MediaFeed
		media     []instax.Media
	}
}

func NewProvider(src Source, client *instax.Client, q []string) (*Provider, error) {
	p := &Provider{source: src, client: client, queries: q}
	switch src {
	case TAGS, MEDIA:
		p.setNextMediaFeed()
	case FOLLOWERS:
		return nil, NotImplemented
	default:
		return nil, UknownSource
	}

	return p, nil
}

func (p *Provider) setNextMediaFeed() (err error) {
	if p.state.queryNum >= len(p.queries) {
		err = EOF
		return
	}

	switch p.source {
	case TAGS:
		p.state.mediaFeed, err = p.client.MediaByFeed(p.queries[p.state.queryNum])
	case MEDIA:
		p.state.mediaFeed, err = p.client.MediaByUser(p.queries[p.state.queryNum])
	}

	if err != nil {
		return
	}

	p.state.queryNum++
	p.state.step = 0
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

	users, err := p.client.Likes(p.state.media[p.state.subStep].ID)
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

func (p *Provider) Next() *Data {
	return &Data{}
}

type Data struct{}

func (d *Data) Media() []instax.Media {
	return []instax.Media{}
}

func (d *Data) Users() []instax.UserShort {
	return []instax.UserShort{}
}
