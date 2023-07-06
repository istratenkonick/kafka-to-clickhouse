CREATE TABLE IF NOT EXISTS logs
(
    name String,
    body String,
    timestamp String
)
    ENGINE = MergeTree()
    ORDER BY timestamp;