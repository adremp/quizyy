-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS quiz (
	Id SERIAL PRIMARY KEY,
	Question TEXT NOT NULL CHECK(Question != ''),
	Answer TEXT NOT NULL CHECK(Answer != ''),
	Variants TEXT[] NOT NULL DEFAULT(ARRAY[]::TEXT[])
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE quiz;
-- +goose StatementEnd
