CREATE TABLE osrs_items_metadata
ENGINE=URL('https://prices.runescape.wiki/api/v1/osrs/mapping', JSON);

CREATE TABLE osrs_prices_5m
(
    id UInt16,
    timestamp DateTime('UTC'),
    avgHighPrice UInt64,
    highPriceVolume UInt64,
    avgLowPrice UInt64,
    lowPriceVolume UInt64,
)
ENGINE = MergeTree
PRIMARY KEY (id, timestamp)
ORDER BY (id, timestamp)
PARTITION BY toStartOfDay(timestamp)