package flags

import "testing"

func TestPostgres_ConnectionString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		pg   Postgres
		want string
	}{
		{
			name: "typical options",
			pg: Postgres{
				PostgresHost: "localhost",
				PostgresPort: "5432",
				PostgresUser: "user",
				PostgresPass: "pass",
				PostgresDB:   "radar",
			},
			want: "postgres://user:pass@localhost:5432/radar",
		},
		{
			name: "special symbols",
			pg: Postgres{
				PostgresHost: "localhost",
				PostgresPort: "5432",
				PostgresUser: "user",
				PostgresPass: "p@$$word!",
				PostgresDB:   "radar",
			},
			want: "postgres://user:p%40%24%24word%21@localhost:5432/radar",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.pg.PostgresConnectionString(); got != tt.want {
				t.Errorf("ConnectionString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
