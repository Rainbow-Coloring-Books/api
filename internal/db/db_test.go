package db_test

import (
	"context"
	"testing"
	"errors"

	"github.com/golang/mock/gomock"
	"rainbowcoloringbooks/internal/db"
)

func TestConnect(t *testing.T) {
	testCases := []struct {
		name    string
		setup   func(*gomock.Controller) db.Database
		wantErr bool
	}{
		{
			name: "valid connection",
			setup: func(ctrl *gomock.Controller) db.Database {
				mockDB := db.NewMockDatabase(ctrl)
				mockDB.EXPECT().Connect(gomock.Any()).Return(nil)
				return mockDB
			},
			wantErr: false,
		},
		{
			name: "invalid credentials",
			setup: func(ctrl *gomock.Controller) db.Database {
				mockDB := db.NewMockDatabase(ctrl)
				mockDB.EXPECT().Connect(gomock.Any()).Return(errors.New("invalid credentials"))
				return mockDB
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := tc.setup(ctrl)
			ctx := context.Background()
			err := mockDB.Connect(ctx)

			if tc.wantErr {
				if err == nil {
					t.Errorf("expected an error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
