package client

import (
	"context"
	"testing"

	"github.com/cloudevents/sdk-go/pkg/event"

	"github.com/google/go-cmp/cmp"
)

func TestWithEventDefaulter(t *testing.T) {

	v2 := func(ctx context.Context, event event.Event) event.Event {
		event.Context = event.Context.AsV02()
		return event
	}

	v3 := func(ctx context.Context, event event.Event) event.Event {
		event.Context = event.Context.AsV03()
		return event
	}

	testCases := map[string]struct {
		c       *ceClient
		fns     []EventDefaulter
		want    int // number of defaulters
		wantErr string
	}{
		"none": {
			c:    &ceClient{},
			want: 0,
		},
		"one": {
			c:    &ceClient{},
			fns:  []EventDefaulter{v2},
			want: 1,
		},
		"two": {
			c:    &ceClient{},
			fns:  []EventDefaulter{v2, v3},
			want: 2,
		},
		"nil fn": {
			c:       &ceClient{},
			fns:     []EventDefaulter{nil},
			wantErr: "client option was given an nil event defaulter",
		},
	}
	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			var err error
			for _, fn := range tc.fns {
				err = tc.c.applyOptions(WithEventDefaulter(fn))
				if err != nil {
					break
				}
			}

			if tc.wantErr != "" || err != nil {
				var gotErr string
				if err != nil {
					gotErr = err.Error()
				}
				if diff := cmp.Diff(tc.wantErr, gotErr); diff != "" {
					t.Errorf("unexpected error (-want, +got) = %v", diff)
				}
				return
			}

			got := len(tc.c.eventDefaulterFns)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("unexpected (-want, +got) = %v", diff)
			}
		})
	}
}

func TestWith_Defaulters(t *testing.T) {

	testCases := map[string]struct {
		c       *ceClient
		opts    []Option
		want    int // number of defaulters
		wantErr string
	}{
		"none": {
			c:    &ceClient{},
			want: 0,
		},
		"uuid": {
			c:    &ceClient{},
			opts: []Option{WithUUIDs()},
			want: 1,
		},
		"time": {
			c:    &ceClient{},
			opts: []Option{WithTimeNow()},
			want: 1,
		},
		"uuid and time": {
			c:    &ceClient{},
			opts: []Option{WithUUIDs(), WithTimeNow()},
			want: 2,
		},
	}
	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			var err error
			if len(tc.opts) > 0 {
				err = tc.c.applyOptions(tc.opts...)
			}

			if tc.wantErr != "" || err != nil {
				var gotErr string
				if err != nil {
					gotErr = err.Error()
				}
				if diff := cmp.Diff(tc.wantErr, gotErr); diff != "" {
					t.Errorf("unexpected error (-want, +got) = %v", diff)
				}
				return
			}

			got := len(tc.c.eventDefaulterFns)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("unexpected (-want, +got) = %v", diff)
			}
		})
	}
}