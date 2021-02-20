package util

import "time"

// SpaceTimeSlicer returns a channel that slices values from an incoming channel with max
// chunksize and TTL contraints.
//
// Gratuitously copied from Elliot Chance:
//   https://elliotchance.medium.com/batch-a-channel-by-size-or-time-in-go-92fa3098f65
func SpaceTimeSlicer(inbound <-chan string, maxItems int, ttl time.Duration) chan []string {
	outbound := make(chan []string)

	go func() {
		defer close(outbound)

		for stop := false; !stop; {

			batch := make([]string, 0, maxItems) // preallocate capacity
			expire := time.After(ttl)

			for {
				select {
				case value, ok := <-inbound:
					if !ok {
						stop = true
						goto done
					}

					batch = append(batch, value)
					if len(batch) == maxItems {
						goto done
					}

				case <-expire:
					goto done
				}
			}

		done:
			if len(batch) > 0 {
				outbound <- batch
			}
		}
	}()

	return outbound
}
