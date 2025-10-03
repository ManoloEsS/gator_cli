package cli

import (
	"testing"

	"github.com/ManoloEsS/gator_cli/test"
)

func TestHandlerReset(t *testing.T) {
	tests := []struct {
		name        string
		cmd         Command
		setupDB     func(*test.MockDb)
		expectError bool
	}{
		{
			name: "reset users successfully",
			cmd: Command{
				Name:      "reset",
				Arguments: []string{},
			},
			setupDB:     func(db *test.MockDb) {},
			expectError: false,
		},
		// Note: We can't test the error case because HandlerReset uses log.Fatalf
		// which would terminate the test process. In a real refactor, we would
		// change HandlerReset to return errors instead of calling log.Fatalf
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := test.NewMockDb()
			mockCfg := &test.MockCfg{}
			tt.setupDB(mockDb)

			state := &State{
				Db:  mockDb,
				Cfg: mockCfg,
			}

			err := HandlerReset(state, tt.cmd)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}
