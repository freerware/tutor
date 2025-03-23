-- +goose Up
-- +goose StatementBegin
CREATE TABLE `POST` (
  `UUID`					  VARCHAR(36)		  NOT NULL,
  `TITLE`					  VARCHAR(255)		NOT NULL,
  `CONTENT`				  TEXT			      NOT NULL,
  `AUTHOR_UUID`			VARCHAR(36)		  NOT NULL,
  `DRAFT`           BOOLEAN         NOT NULL,
  `LIKE_COUNT`      INT             NOT NULL,
  `CREATED_AT`			DATETIME		    NOT NULL,
  `UPDATED_AT`			DATETIME		    NOT NULL,
	`DELETED_AT`			DATETIME		    NULL,
  
  PRIMARY KEY (`UUID`),
  FOREIGN KEY (`AUTHOR_UUID`) REFERENCES `ACCOUNT`(`UUID`)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `POST`;
-- +goose StatementEnd
