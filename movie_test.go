package radarr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	internal "github.com/SkYNewZ/radarr/internal/radarr"
)

func Test_newMovieService(t *testing.T) {
	s := &Service{client: http.DefaultClient, url: internal.DummyURL}

	tests := []struct {
		name    string
		service *Service
		want    *MovieService
	}{
		{
			name:    "Constructor",
			service: s,
			want:    &MovieService{s},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newMovieService(tt.service); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newMovieService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMovieService_Get(t *testing.T) {
	var expectedMovie *Movie
	err := json.Unmarshal([]byte(internal.DummyMovieResponse), &expectedMovie)
	if err != nil {
		t.Errorf("json.Unmarshal() error: %s", err.Error())
	}

	goodService := newMovieService(&Service{
		client: internal.DummyHTTPClient,
		url:    internal.DummyURL,
	})

	tests := []struct {
		name    string
		service *Service
		movieID int
		want    *Movie
		wantErr bool
	}{
		{
			name:    "Same response",
			movieID: 217,
			service: goodService.s,
			want:    expectedMovie,
			wantErr: false,
		},
	}

	var m *MovieService
	var got *Movie
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m = &MovieService{tt.service}
			got, err = m.Get(tt.movieID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MovieService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MovieService.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMovieService_List(t *testing.T) {
	var expectedMovies Movies
	err := json.Unmarshal([]byte(fmt.Sprintf("[%s, %s]", internal.DummyMovieResponse, internal.DummyMovieResponse)), &expectedMovies)
	if err != nil {
		t.Errorf("json.Unmarshal() error: %s", err.Error())
	}

	goodService := newMovieService(&Service{
		client: internal.DummyHTTPClient,
		url:    internal.DummyURL,
	})

	tests := []struct {
		name    string
		service *Service
		want    Movies
		wantErr bool
	}{
		{
			name:    "Same response",
			service: goodService.s,
			want:    expectedMovies,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MovieService{tt.service}
			got, err := m.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("MovieService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MovieService.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMovieService_Upcoming(t *testing.T) {
	var expectedMovies Movies
	err := json.Unmarshal([]byte(fmt.Sprintf("[%s, %s]", internal.DummyMovieResponse, internal.DummyMovieResponse)), &expectedMovies)
	if err != nil {
		t.Errorf("json.Unmarshal() error: %s", err.Error())
	}

	var expectedMovie *Movie
	err = json.Unmarshal([]byte(internal.DummyMovieResponse), &expectedMovie)
	if err != nil {
		t.Errorf("json.Unmarshal() error: %s", err.Error())
	}

	goodService := newMovieService(&Service{
		client: internal.DummyHTTPClient,
		url:    internal.DummyURL,
	})

	tests := []struct {
		name    string
		service *Service
		opts    []*UpcomingOptions
		want    Movies
		wantErr bool
	}{
		{
			name:    "Without filter",
			opts:    nil,
			service: goodService.s,
			want:    Movies{},
			wantErr: false,
		},
		{
			name:    "Dates with reverse order",
			service: goodService.s,
			wantErr: true,
			want:    nil,
			opts: func() []*UpcomingOptions {
				s := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.Local)
				e := time.Date(2010, time.January, 1, 0, 0, 0, 0, time.Local)
				return []*UpcomingOptions{{Start: &s, End: &e}}
			}(),
		},
		{
			name:    "Start filter",
			service: goodService.s,
			wantErr: false,
			want:    expectedMovies,
			opts: func() []*UpcomingOptions {
				s := time.Date(2019, time.November, 19, 23, 0, 0, 0, time.UTC)
				return []*UpcomingOptions{{Start: &s}}
			}(),
		},
		{
			name:    "End filter",
			service: goodService.s,
			want:    Movies{},
			wantErr: false,
			opts: func() []*UpcomingOptions {
				e := time.Date(2019, time.November, 20, 23, 0, 0, 0, time.UTC)
				return []*UpcomingOptions{{End: &e}}
			}(),
		},
		{
			name:    "Both filters",
			service: goodService.s,
			want:    Movies{expectedMovie},
			wantErr: false,
			opts: func() []*UpcomingOptions {
				start := time.Date(2019, time.November, 19, 23, 0, 0, 0, time.UTC)
				end := time.Date(2019, time.November, 20, 23, 0, 0, 0, time.UTC)
				return []*UpcomingOptions{{Start: &start, End: &end}}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MovieService{tt.service}
			got, err := m.Upcoming(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MovieService.Upcoming() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MovieService.Upcoming() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMovieService_Delete(t *testing.T) {
	type args struct {
		movie *Movie
		opts  []*DeleteMovieOptions
	}

	var dummyMovie *Movie
	err := json.Unmarshal([]byte(internal.DummyMovieResponse), &dummyMovie)
	if err != nil {
		t.Errorf("json.Unmarshal() error: %s", err.Error())
	}

	goodService := newMovieService(&Service{
		client: internal.DummyHTTPClient,
		url:    internal.DummyURL,
	})

	tests := []struct {
		name    string
		service *Service
		args    args
		wantErr bool
	}{
		{
			name:    "Delete without option",
			service: goodService.s,
			args:    args{movie: dummyMovie},
			wantErr: false,
		},
		{
			name:    "Delete with addExclusion option",
			args:    args{movie: dummyMovie, opts: []*DeleteMovieOptions{{AddExclusion: true}}},
			service: goodService.s,
			wantErr: false,
		},
		{
			name:    "Delete with deleteFiles option",
			args:    args{movie: dummyMovie, opts: []*DeleteMovieOptions{{DeleteFiles: true}}},
			service: goodService.s,
			wantErr: false,
		},
		{
			name:    "Delete with both options",
			args:    args{movie: dummyMovie, opts: []*DeleteMovieOptions{{DeleteFiles: true, AddExclusion: true}}},
			service: goodService.s,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MovieService{
				s: tt.service,
			}
			if err := m.Delete(tt.args.movie, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("MovieService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMovieService_Excluded(t *testing.T) {
	var expectedMovies ExcludedMovies
	err := json.Unmarshal([]byte(internal.DummyExcludedMovies), &expectedMovies)
	if err != nil {
		t.Errorf("json.Unmarshal() error: %s", err.Error())
	}

	tests := []struct {
		name    string
		s       *Service
		want    ExcludedMovies
		wantErr bool
	}{
		{
			name: "Expected response",
			s: newMovieService(&Service{
				client: internal.DummyHTTPClient,
				url:    internal.DummyURL,
			}).s,
			want:    expectedMovies,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MovieService{
				s: tt.s,
			}
			got, err := m.Excluded()
			if (err != nil) != tt.wantErr {
				t.Errorf("MovieService.Excluded() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MovieService.Excluded() = %v, want %v", got, tt.want)
			}
		})
	}
}
