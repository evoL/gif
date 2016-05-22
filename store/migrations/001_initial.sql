-- +migrate Up
CREATE TABLE images (
  id VARCHAR(40) PRIMARY KEY,
  url TEXT,
  added_at DATETIME NOT NULL
);

CREATE TABLE image_tags (
  image_id VARCHAR(40) NOT NULL,
  tag VARCHAR(255) NOT NULL
);

CREATE INDEX image_tags_index ON image_tags (tag);
CREATE UNIQUE INDEX image_tags_unique ON image_tags (image_id, tag);

-- +migrate Down
DROP INDEX image_tags_index;
DROP INDEX image_tags_unique;
DROP TABLE images;
DROP TABLE image_tags;
