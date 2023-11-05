package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gudimz/urlShortener/internal/app/repository/psql/models"
	"github.com/gudimz/urlShortener/internal/pkg/ds"
)

type testSuite struct {
	mockRepo *MockRepository
	service  *Service
}

const (
	testShortURL = "shortURL"
	testURL      = "https://youtube.com/"
)

func initTestSuite(t *testing.T) *testSuite {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	service := New(mockRepo)

	return &testSuite{
		mockRepo: mockRepo,
		service:  service,
	}
}

func TestService_CreateShorten(t *testing.T) {
	s := initTestSuite(t)

	errCommon := errors.New("some error")
	now := time.Now().UTC()
	expectedDBShorten := &models.DBShorten{
		ShortURL:    testShortURL,
		OriginURL:   testURL,
		Visits:      0,
		DateCreated: now,
		DateUpdated: now,
	}

	expectedShorten := &ds.Shorten{
		ShortURL:    testShortURL,
		OriginURL:   testURL,
		Visits:      0,
		DateCreated: now,
		DateUpdated: now,
	}

	testCases := []struct {
		name          string
		mock          func()
		input         ds.InputShorten
		want          *ds.Shorten
		expectedError error
	}{
		{
			name: "success: create shorten with given short URL",
			mock: func() {
				s.mockRepo.EXPECT().
					CreateShorten(gomock.Any(), gomock.Any()).
					Return(expectedDBShorten, nil).Times(1)
			},
			input: ds.InputShorten{
				ShortenURL: mo.Some(testShortURL),
				OriginURL:  testURL,
			},
			want:          expectedShorten,
			expectedError: nil,
		},
		{
			name: "success: create shorten with generate short URL",
			mock: func() {
				s.mockRepo.EXPECT().
					CreateShorten(gomock.Any(), gomock.Any()).
					Return(expectedDBShorten, nil).Times(1)
			},
			input: ds.InputShorten{
				ShortenURL: mo.None[string](),
				OriginURL:  testURL,
			},
			want:          expectedShorten,
			expectedError: nil,
		},
		{
			name: "failed: short url already exist",
			mock: func() {
				s.mockRepo.EXPECT().
					CreateShorten(gomock.Any(), gomock.Any()).
					Return(nil, &pgconn.PgError{Code: "23505"}).Times(1)
			},
			input: ds.InputShorten{
				ShortenURL: mo.None[string](),
				OriginURL:  testURL,
			},
			want:          nil,
			expectedError: ds.ErrShortURLAlreadyExists,
		},
		{
			name: "failed: common error",
			mock: func() {
				s.mockRepo.EXPECT().
					CreateShorten(gomock.Any(), gomock.Any()).
					Return(nil, errCommon)
			},
			input: ds.InputShorten{
				ShortenURL: mo.None[string](),
				OriginURL:  testURL,
			},
			want:          nil,
			expectedError: errCommon,
		},
	}

	for _, tt := range testCases {
		if tt.mock != nil {
			tt.mock()
		}

		got, err := s.service.CreateShorten(context.Background(), tt.input)
		if tt.expectedError != nil {
			require.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)
		} else {
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}
	}
}

func TestService_GetShorten(t *testing.T) {
	s := initTestSuite(t)

	errCommon := errors.New("some error")
	now := time.Now().UTC()
	expectedDBShorten := &models.DBShorten{
		ShortURL:    testShortURL,
		OriginURL:   testURL,
		Visits:      0,
		DateCreated: now,
		DateUpdated: now,
	}

	expectedShorten := &ds.Shorten{
		ShortURL:    testShortURL,
		OriginURL:   testURL,
		Visits:      0,
		DateCreated: now,
		DateUpdated: now,
	}

	testCases := []struct {
		name          string
		mock          func()
		input         string
		want          *ds.Shorten
		expectedError error
	}{
		{
			name: "success: url shorten exist",
			mock: func() {
				s.mockRepo.EXPECT().
					GetShorten(gomock.Any(), gomock.Any()).
					Return(expectedDBShorten, nil).Times(1)
			},
			input:         testShortURL,
			want:          expectedShorten,
			expectedError: nil,
		},
		{
			name: "failed: url shorten not found",
			mock: func() {
				s.mockRepo.EXPECT().
					GetShorten(gomock.Any(), gomock.Any()).
					Return(nil, pgx.ErrNoRows).Times(1)
			},
			input:         testShortURL,
			want:          nil,
			expectedError: ds.ErrShortURLNotFound,
		},
		{
			name: "failed: common error",
			mock: func() {
				s.mockRepo.EXPECT().
					GetShorten(gomock.Any(), gomock.Any()).
					Return(nil, errCommon).Times(1)
			},
			input:         testShortURL,
			want:          nil,
			expectedError: errCommon,
		},
	}

	for _, tt := range testCases {
		if tt.mock != nil {
			tt.mock()
		}

		got, err := s.service.GetShorten(context.Background(), tt.input)
		if tt.expectedError != nil {
			require.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)
		} else {
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}
	}
}

func TestService_Redirect(t *testing.T) {
	s := initTestSuite(t)

	errCommon := errors.New("some error")
	now := time.Now().UTC()
	expectedDBShorten := &models.DBShorten{
		ShortURL:    testShortURL,
		OriginURL:   testURL,
		Visits:      0,
		DateCreated: now,
		DateUpdated: now,
	}

	testCases := []struct {
		name          string
		mock          func()
		input         string
		want          string
		expectedError error
	}{
		{
			name: "success: url shorten exist",
			mock: func() {
				s.mockRepo.EXPECT().
					GetShorten(gomock.Any(), gomock.Any()).
					Return(expectedDBShorten, nil).Times(1)
				s.mockRepo.EXPECT().
					UpdateShorten(gomock.Any(), gomock.Any()).
					Return(nil).Times(1)
			},
			input:         testShortURL,
			want:          testURL,
			expectedError: nil,
		},
		{
			name: "failed: url shorten not found",
			mock: func() {
				s.mockRepo.EXPECT().
					GetShorten(gomock.Any(), gomock.Any()).
					Return(nil, pgx.ErrNoRows).Times(1)
			},
			input:         testShortURL,
			want:          "",
			expectedError: ds.ErrShortURLNotFound,
		},
		{
			name: "failed: GetShorten common error",
			mock: func() {
				s.mockRepo.EXPECT().
					GetShorten(gomock.Any(), gomock.Any()).
					Return(nil, errCommon).Times(1)
			},
			input:         testShortURL,
			want:          "",
			expectedError: errCommon,
		},
		{
			name: "failed: UpdateShorten common error",
			mock: func() {
				s.mockRepo.EXPECT().
					GetShorten(gomock.Any(), gomock.Any()).
					Return(expectedDBShorten, nil).Times(1)
				s.mockRepo.EXPECT().
					UpdateShorten(gomock.Any(), gomock.Any()).
					Return(errCommon).Times(1)
			},
			input:         testShortURL,
			want:          "",
			expectedError: errCommon,
		},
	}

	for _, tt := range testCases {
		if tt.mock != nil {
			tt.mock()
		}

		got, err := s.service.Redirect(context.Background(), tt.input)
		if tt.expectedError != nil {
			require.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)
		} else {
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}
	}
}

func TestService_DeleteShorten(t *testing.T) {
	s := initTestSuite(t)

	errCommon := errors.New("some error")

	testCases := []struct {
		name          string
		mock          func()
		input         string
		expectedError error
	}{
		{
			name: "success: url shorten delete",
			mock: func() {
				s.mockRepo.EXPECT().
					DeleteShorten(gomock.Any(), gomock.Any()).
					Return(int64(1), nil).Times(1)
			},
			input:         testShortURL,
			expectedError: nil,
		},
		{
			name: "failed: url shorten not found",
			mock: func() {
				s.mockRepo.EXPECT().
					DeleteShorten(gomock.Any(), gomock.Any()).
					Return(int64(0), nil).Times(1)
			},
			input:         testShortURL,
			expectedError: ds.ErrShortURLNotFound,
		},
		{
			name: "failed: common error",
			mock: func() {
				s.mockRepo.EXPECT().
					DeleteShorten(gomock.Any(), gomock.Any()).
					Return(int64(0), errCommon).Times(1)
			},
			input:         testShortURL,
			expectedError: errCommon,
		},
	}

	for _, tt := range testCases {
		if tt.mock != nil {
			tt.mock()
		}

		err := s.service.DeleteShorten(context.Background(), tt.input)
		if tt.expectedError != nil {
			require.Equal(t, tt.expectedError, err)
		} else {
			require.NoError(t, err)
		}
	}
}
