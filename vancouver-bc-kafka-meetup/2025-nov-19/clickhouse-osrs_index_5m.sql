CREATE TABLE osrs_index_5m
(
    index_name String,
    divisor UInt64,

    id UInt16,
    timestamp DateTime('UTC'),
    name String,
    icon String,
    avgHighPrice UInt64,
    highPriceVolume UInt64,
    avgLowPrice UInt64,
    lowPriceVolume UInt64,
)
ENGINE = MergeTree
PRIMARY KEY (index_name, timestamp)
ORDER BY (index_name, timestamp)
PARTITION BY toStartOfDay(timestamp)