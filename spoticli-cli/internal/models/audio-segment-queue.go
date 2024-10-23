package models

import (
	"fmt"

	"github.com/gopxl/beep/v2"
)

type AudioSegmentQueue struct {
	segmentCtr int // number representing which segment number is being played
	streamers  []beep.Streamer
}

func (q *AudioSegmentQueue) Add(streamers ...beep.Streamer) {
	q.streamers = append(q.streamers, streamers...)
	fmt.Printf("%d segments in queue\n", len(q.streamers))
}

func (q *AudioSegmentQueue) Stream(samples [][2]float64) (n int, ok bool) {
	// We use the filled variable to track how many samples we've
	// successfully filled already. We loop until all samples are filled.
	filled := 0
	for filled < len(samples) {
		// There are no streamers in the queue, so we stream silence.
		if len(q.streamers) == 0 {
			for i := range samples[filled:] {
				samples[i][0] = 0
				samples[i][1] = 0
			}
			break
		}

		// We stream from the first streamer in the queue.
		n, ok := q.streamers[0].Stream(samples[filled:])
		// If it's drained, we pop it from the queue, thus continuing with
		// the next streamer.
		if !ok {
			q.streamers = q.streamers[1:]
			q.segmentCtr += 1
			fmt.Printf("segment %d has been dequeued with %d samples left in current segment and %d segments in queue\n", q.segmentCtr, n, len(q.streamers))
		}
		// We update the number of filled samples.
		filled += n
	}
	return len(samples), true
}

func (q *AudioSegmentQueue) Err() error {
	return nil
}
