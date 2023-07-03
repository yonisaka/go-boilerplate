package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yonisaka/go-boilerplate/internal/entities/repository"
)

var errInternal = errors.New("error")

func TestHealthUC_Liveness(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	type test struct {
		fields  fields
		args    args
		want    string
		wantErr error
	}

	tests := map[string]func(t *testing.T, ctrl *gomock.Controller) test{
		"Given valid request of Liveness, When repository executed successfully, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx: ctx,
			}

			want := "Life is not about waiting for the storms to pass; " +
				"it's about learning to dance in the rain and embracing every challenge as an opportunity for growth."

			mockJourneyRepo := repository.NewGoMockHealthRepo(ctrl)
			mockJourneyRepo.EXPECT().GetLiveness(args.ctx).Return(nil)

			return test{
				fields: fields{
					healthRepo: mockJourneyRepo,
				},
				args:    args,
				want:    want,
				wantErr: nil,
			}
		},
		"Given valid request of Liveness, When repository executed failed, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx: ctx,
			}

			mockJourneyRepo := repository.NewGoMockHealthRepo(ctrl)
			mockJourneyRepo.EXPECT().GetLiveness(args.ctx).Return(errInternal)

			return test{
				fields: fields{
					healthRepo: mockJourneyRepo,
				},
				args:    args,
				want:    "",
				wantErr: errInternal,
			}
		},
	}

	for name, testFn := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt := testFn(t, ctrl)

			sut := sut(tt.fields)

			got, err := sut.Liveness(tt.args.ctx)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
