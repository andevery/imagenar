package autogram

import (
	"github.com/andevery/instax"
)

type Liker struct {
	Provider *Provider
}

func DefaultLiker(p *Provider) *Liker {
	return &Liker{Provider: p}
}

func (l *Liker) LikeAFew(media []instax.Media, count int) {
	for _, i := range randomIndexes(len(media), count) {
		if l.isMediaMatch(&media[i]) {
			l.Provider.WebClient().Like(&media[i])
		}
	}
}

func (l *Liker) LikeABatch(media []instax.Media) {
	l.LikeAFew(media, len(media))
}

func (l *Liker) isMediaMatch(media *instax.Media) bool {
	return !media.UserHasLiked
}
