
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `ACCOUNT` (
	`UUID`					      VARCHAR(36)		NOT NULL,
	`GIVEN_NAME`			    VARCHAR(128)	NOT NULL,
	`SURNAME`				      VARCHAR(128)	NOT NULL,
	`PRIMARY_CREDENTIAL`	VARCHAR(128)	NOT NULL,
	`CREATED_AT`			    DATETIME		  NOT NULL,
	`UPDATED_AT`			    DATETIME		  NOT NULL,
	`DELETED_AT`			    DATETIME		  NULL,
	
	PRIMARY KEY (`UUID`)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `ACCOUNT`;
