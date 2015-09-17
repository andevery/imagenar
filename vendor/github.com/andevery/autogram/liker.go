package autogram

import (
	"github.com/andevery/instaw"
	"github.com/andevery/instax"
)

type Liker struct {
	Limiter *Limiter

	WebClient *instaw.Client
	ApiClient *instax.Client
}

func (l *Liker) LikeAFew(media []instax.Media, count int) {
	for _, i := range randomIndexes(len(media), count) {
		if l.isMediaMatch(&media[i]) {
			<-l.Limiter.Timer
			l.WebClient.Like(&media[i])
		}
	}
}

func (l *Liker) LikeABatch(media []instax.Media) {
	l.LikeAFew(media, len(media))
}

func (l *Liker) isMediaMatch(media *instax.Media) bool {
	return !media.UserHasLiked
}
