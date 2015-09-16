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

}

func (l *Liker) LikeABatch(media []instax.Media) {
	for i, _ := range media {
		if l.isMediaMatch(&media[i]) {
			<-l.Limiter.Timer
			l.WebClient.Like(&media[i])
		}
	}
}

func (l *Liker) isMediaMatch(media *instax.Media) bool {
	return true
}
