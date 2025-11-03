CREATE TABLE osrs_prices_5m
(
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
PRIMARY KEY (id, timestamp)
ORDER BY (id, timestamp)
PARTITION BY toStartOfDay(timestamp)