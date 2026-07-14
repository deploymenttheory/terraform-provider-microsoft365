package crud

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"
)

// FatalPollError marks an error returned by a PollUntil check as non-retryable,
// aborting the poll loop immediately instead of waiting for the context deadline.
type FatalPollError struct {
	Err error
}

func (e *FatalPollError) Error() string {
	return e.Err.Error()
}

func (e *FatalPollError) Unwrap() error {
	return e.Err
}

// PollUntil repeatedly invokes check every interval until check reports done, check returns
// a *FatalPollError, or the context deadline is reached. It complements ReadWithRetry for
// verification loops that cannot be expressed as a resource Read (e.g. polling until an
// asynchronously-processed deletion is actually reflected by the API before removing the
// resource from state).
//
// The check contract: return (true, nil) when the awaited condition holds; return
// (false, err) with an err describing the still-pending condition (or the transient API
// error) otherwise. The last such error is wrapped into the deadline error so the caller
// can surface why the poll never completed. Wrap an error in *FatalPollError to abort
// immediately on permanent failures (e.g. authorization errors).
func PollUntil(ctx context.Context, interval time.Duration, check func(ctx context.Context) (bool, error)) error {
	var lastErr error
	for {
		done, err := check(ctx)
		if err != nil {
			var fatal *FatalPollError
			if stderrors.As(err, &fatal) {
				return fatal.Err
			}
			lastErr = err
		} else if done {
			return nil
		}

		timer := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			timer.Stop()
			if lastErr != nil {
				return fmt.Errorf("context deadline reached while polling: %w", lastErr)
			}
			return fmt.Errorf("context deadline reached while polling: %w", ctx.Err())
		case <-timer.C:
		}
	}
}
