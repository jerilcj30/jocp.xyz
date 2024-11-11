-- +goose Up
-- +goose StatementBegin
CREATE TABLE leads (
    id SERIAL PRIMARY KEY,
    lead JSON    
); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE leads;
-- +goose StatementEnd
