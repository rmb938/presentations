---
# You can also start simply with 'default'
theme: seriph
# random image from a curated Unsplash collection by Anthony
# like them? see https://unsplash.com/collections/94734566/slidev
background: /images/mario-verduzco-xSdFf1Lcx6o-unsplash.jpg
# some information about your slides (markdown enabled)
title: From Ticker Tape to Trendlines
info: |
  A Stream Processing Journey into Market Dynamics
# apply unocss classes to the current slide
class: text-center
# https://sli.dev/features/drawing
drawings:
  persist: false
# slide transition: https://sli.dev/guide/animations.html#slide-transitions
transition: slide-left
# enable MDC Syntax: https://sli.dev/features/mdc
mdc: true
# open graph
# seoMeta:
#  ogImage: https://cover.sli.dev

fonts:
  sans: Plus Jakarta
  weights: 400,500,600
---
# From Ticker Tape to Trendlines

<div class="mt-12 py-1">
  A Stream Processing Journey into Market Dynamics
</div>

<!--
The last comment block of each slide will be treated as slide notes. It will be visible and editable in Presenter Mode along with the slide. [Read more in the docs](https://sli.dev/guide/syntax.html#notes)
-->

---
layout: profile
image: /images/saKWC52.jpg
---

# Ryan Belgrave
Staff Engineer @ Confluent - WarpStream

~20 Years of technical experience, self-taught with Java in 2005

Ordered my first Ubuntu liveCD soon after and forever got stuck in vim

Previously helped run the infrastructure automation for a few popular Minecraft servers

Deployed the largest deployment of baremetal K8s in the mid-west, probably the US, for a few months in early 2019.

<br />

<carbon-logo-github /> [rmb938](https://github.com/rmb938)
<carbon-logo-linkedin /> [rbelgrave](https://www.linkedin.com/in/rbelgrave/)

---
layout: two-cols
---
# Trading Terminals at Home
If you don't mind the delay

Analyze over 4 years of trading data

Build our own Indexes

Injest "real-time" data from freely available APIs

Do our own calculations on economy health

All using streaming processing technologies, Apache Iceberg, Clickhouse, and Grafana

::right::

<img src="/images/arthur-a-9rz5x8LGBb8-unsplash.jpg" class="w-100"/>

---
layout: two-cols
---

# Stream Processing the Events
Using Bento

https://warpstreamlabs.github.io/bento/

<style>
.image-container {
  /* Enables Flexbox layout for direct children */
  display: flex;

  /* Centers the images horizontally within the container (optional) */
  justify-content: center;

  /* Adds some space around the images (optional) */
  gap: 300px;

  /* Sets a maximum width for the container (optional) */
  max-width: 800px;
  margin: 0 auto;
}

.image-container img {
  /* Makes the images take up equal space */
  /* flex: 1; */

  /* Ensures the images don't exceed the container's width */
  max-width: 100%;

  /* Ensures images maintain their aspect ratio */
  height: auto;
}
</style>

<div>
  <img src="/images/bento_boringly_easy.png" width="320"/>
  <img src="/images/qr-code-bento.svg" width="220"/>
</div>

::right::

````md magic-move
```yaml
input:
  gcp_pubsub:
    project: foo
    subscription: bar

# Mapping Example
pipeline:
  processors:
    - mapping: |
        root.message = this
        root.meta.link_count = this.links.length()
        root.user.age = this.user.age.number()

output:
  redis_streams:
    url: tcp://TODO:6379
    stream: baz
    max_in_flight: 20
```
```yaml
...
# Multiplexing Outputs Example
output:
  switch:
    cases:
      - check: doc.tags.contains("AWS")
        output:
          aws_sqs:
            url: https://sqs.us-west-2.amazonaws.com/TODO/TODO
            max_in_flight: 20

      - output:
          redis_pubsub:
            url: tcp://TODO:6379
            channel: baz
            max_in_flight: 20
```
```yaml
...
# Windowing Example
buffer:
  system_window:
    timestamp_mapping: root = this.created_at
    size: 1h

pipeline:
  processors:
    - group_by_value:
        value: '${! json("traffic_light_id") }'
    - mapping: |
        root = if batch_index() == 0 {
          {
            "traffic_light_id": this.traffic_light_id,
            "created_at": @window_end_timestamp,
            "total_cars": json("registration_plate").from_all().unique().length(),
            "passengers": json("passengers").from_all().sum(),
          }
        } else { deleted() }

...
```
```yaml
...
# External Enrichments Example
pipeline:
  processors:
    - branch:
        request_map: |
          root.id = this.doc.id
          root.content = this.doc.body
        processors:
          - aws_lambda:
              function: sentiment_analysis
        result_map: root.results.sentiment = this

...
```
```yaml
# Easily Extenable with Custom Plugins
input:
  my_custom_input_plugin: {}

pipeline:
  processors:
    - my_super_secret_processor: {}

output:
  my_amazing_output_plugin: {}
...
```
````

---
layout: quote
---
# <logos-kafka-icon /> Kafka
If we are stream processing we need a stream!

Kafka is a distributed event streaming platform used for high-performance data pipelines, streaming analytics, data integration, and mission-critical applications.

https://kafka.apache.org/

---
transition: fade
---
# WarpStream
Zero-Disk Kafka with no inter-zone networking

![](/images/without_warpstream.png)

---
---
# WarpStream
Zero-Disk Kafka with no inter-zone networking

![](/images/with_warpstream.png)

---
---
# WarpStream Agent Groups
Connect to Kafka easily Across Network Boundries

<img src="/images/warpstream-agent-groups.svg" width="720"/>

---
---
# Apache Iceberg
Spark and Compaction

![](/images/iceberg_without_warpstream.png)

---
---
# WarpStream Tableflow
Materialize Kafka Topics as Iceberg Tables

![](/images/iceberg_with_warpstream.png)

---
layout: two-cols
---
# WarpStream Tableflow
Any Kafka-compatible Source

WarpStream Tableflow works with **any** Kafka-compatible source (**Open-Source Kafka**, 
**MSK**, **Confluent Cloud**, **WarpStream**, etc), and can run in any cloud or even 
on-premise. 

<br />

Ingest simultaneously from multiple different Kafka clusters to centralize your data in 
a single lake.

<br />

Kafka in and Apache Iceberg out means there is no vendor lock-in.

::right::

![](/images/ws_tableflow_any_kafka.svg)
<img src="/images/ws_tableflow_any_kafka_sticker.svg" width="250"/>

---
layout: two-cols
---
# Fully Bring your Own Cloud
Your VPC, Your Bucket, No Remote Access

WarpStream has **zero** access to your data or your network

WarpStream requires **no** remote access into your VPCs

<img src="/images/qr-code-warpstream-security-considerations.svg" width="200"/>

::right::

**Impossible** for WarpStream to access your data under any circumstances

Your data stays in **your cloud**, **your vpc**, and **your bucket**

WarpStream Control Plane only stores **metadata** like topic names, topic configs, partition 
counts.

http://console.warpstream.com/security-considerations

---
layout: full
---

<img src="/images/warpstream_logo.svg" width="1000"/>

## [warpstream.com](https://www.warpstream.com/)
## [warpstream.com/tableflow](https://www.warpstream.com/tableflow)

<br />
<br />

<style>
.image-container {
  /* Enables Flexbox layout for direct children */
  display: flex;

  /* Centers the images horizontally within the container (optional) */
  justify-content: center;

  /* Adds some space around the images (optional) */
  gap: 300px;

  /* Sets a maximum width for the container (optional) */
  max-width: 800px;
  margin: 0 auto;
}

.image-container img {
  /* Makes the images take up equal space */
  /* flex: 1; */

  /* Ensures the images don't exceed the container's width */
  max-width: 100%;

  /* Ensures images maintain their aspect ratio */
  height: auto;
}
</style>

<div class="image-container">
  <img src="/images/qr-code-warpstream.svg" width="200"/>
  <img src="/images/qr-code-warpstream-tableflow.svg" width="200"/>
</div>

---
---
# Loading Data With Bento

````md magic-move
```yaml {*|3|5}
generate:
    interval: 5m
    batch_size: 1
    auto_replay_nacks: false
    count: 0
    mapping: |
      root = {}
```
```yaml {*|5}
- http:
    # Query 5m with timestamp
    # Timestamp is rounded to nearest 5 minute interval then goes back in time by 10 mins so we have data.
    # If we don't go back in time by 10 mins, we query the current 5 min interval which may not be populated at query time.
    url: https://example.com/api/v1/trades/5m?timestamp=${! ((((timestamp_unix().number() / 300).round()) * 300) - 600) }
    verb: GET
    retries: 100
    headers:
      # Custom User agent with discord username as recommended by API devs
      User-Agent: WarpStream - Bento - @rmb938 in Discord
    parallel: false
```
````

---
---
# Deploying with Docker Compose

```yaml {*|1-9|10-17}{maxHeight:'400px'}
services:
  warpstream-agent-kafka:
    image: public.ecr.aws/warpstream-labs/warpstream_agent:latest
    restart: unless-stopped
    command:
      - agent
    env_file: .env/warpstream-agent-kafka
    ports:
      - 9092:9092
  bento:
    image: ghcr.io/warpstreamlabs/bento
    restart: unless-stopped
    command:
      - --log.level
      - debug
    volumes:
      - ./bento.yaml:/bento.yaml
```

---
---
# Loading All the Data into a Kafka Topic
Lots of Records

![](/images/Screenshot_20251029_140910.png)

---
---
# Tableflow Configuration
Like Bento it's YAML and Simple

```yaml {*|4-5|6|8-9|10|12|13-22}
source_clusters:
  - name: "warpstream"
    bootstrap_brokers:
      - hostname: "warpstream-agent-kafka"
        port: 9092
destination_bucket_url: "gs://rmb-lab-vancouver_kafka_meetup_2025_nov_19"
tables:
  - source_cluster_name: "warpstream"
    source_topic: "prices_5m"
    source_format: "json"
    schema_mode: "inline"
    partitioning_scheme: "hour"
    schema:
      fields:
        - { name: id, type: int, id: 1 }
        - { name: timestamp, type: timestamp, id: 2 }
        - { name: name, type: string, id: 3 }
        - { name: icon, type: string, id: 4 }
        - { name: avgHighPrice, type: int, id: 5 }
        - { name: highPriceVolume, type: int, id: 6 }
        - { name: avgLowPrice, type: int, id: 7 }
        - { name: lowPriceVolume, type: int, id: 8 }
```

---
layout: center
---

# So where did the data come from?
We all were thinking stocks... it's not really stocks

---
layout: center
---

![](/images/Old_School_RuneScape_logo.png)

---
layout: center
---
# Why not use real stock market data?
The APIs are expensive and the free ones have limits on the number of requests per day


---
layout: two-cols
---
# RuneScape really has a giant economy
But it's just a game!?

Billions of items traded daily through the ingame trading hub called the Grand Exchange

It's like the Stock Market where players can create orders to buy and sell assets at different prices

Orders are automatically fulfilled based on cost and order placement

Large organizations of "Day Traders" buying and selling items to make small amounts of profit

::right::

![](/images/Grand_Exchange_logo.png)

---
---
# Grand Exchange Pricing Site
Automated Pricing Site ran by the Wiki

Items with highest daily volume - https://prices.runescape.wiki/

![](/images/Screenshot_20251001_171409.png)

---
---
# Real time and historical Prices API
Yes there's an API for a game

API endpoint: `prices.runescape.wiki/api/v1/osrs`

`/latest`
> Get the latest high and low prices for the items that we have data for, and the Unix timestamp when that transaction took place.
> Map from itemId (see here for a reference) to an object of (high, highTime, low, lowTime). If we've never seen an item traded, it won't be in the response. If we've never seen an instant-buy price, then high and highTime will be null (and similarly for low and lowTime if we've never seen an instant-sell). 


`/5m`
> Gives 5-minute average of item high and low prices as well as the number traded for the items that we have data on. Comes with a Unix timestamp indicating the 5 minute block the data is from. 

---
---
# Using the API

```yaml {*|4}
- branch:
  processors:
    - http:
        url: https://prices.runescape.wiki/api/v1/osrs/5m?timestamp=${! this.timestamp }
        verb: GET
        retries: 100
        headers:
          # Custom User agent with discord username as recommended by API devs
          User-Agent: WarpStream - Bento - @rmb938 in Discord
        parallel: false
  result_map: |
    root.data = this.data
```

---
---
# Full Bento Config
A bit complex but still simple

```yaml {*|1-8|26-28|30-35|37-47|65-71|71-82|92-105}{maxHeight:'400px'}
input:
  generate:
    interval: 1s
    batch_size: 1
    auto_replay_nacks: true
    count: 0
    mapping: |
      root = {}

pipeline:
  processors:
    # Load the current timestamp offset from cache
    # Try to pull timestamp from cache
    - cache:
        resource: timestamp_offset
        key: timestamp_offset
        operator: get

    # Otherwise set it to 0
    - catch:
        - mapping: |
            root = {}
            root.timestamp_offset = 0

    # Increment the offset
    - mapping: |
        root = this
        root.timestamp_offset = this.timestamp_offset + 300

    # Set the timestamp from the offset
    - mapping: |
        root = this

        # Offset from Monday, March 8, 2021 8:00:00 AM
        root.timestamp = 1615190400 + this.timestamp_offset

    # If timestamp > (now minus 10 minutes) delete the messages
    # We will retry the same offset over and over until the message isn't deleted
    - mapping: |
        root = this

        # Round to the nearest 5 minute interval and subtract 10 minutes
        let cuttoff_ts = ((((timestamp_unix().number() / 300).round()) * 300) - 600)

        if root.timestamp > $cuttoff_ts {
          root = deleted()
        }

    # Finally save the content to the cache so we save the timestamp offset
    - cache:
        resource: timestamp_offset
        key: timestamp_offset
        operator: set
        value: "${! content() }"

    # Remove timestamp offset
    - mapping: |
        root = this
        root.timestamp_offset = deleted()

    # Query for the item mapping
    - branch:
        request_map: 'root = ""'
        processors:
          # Try to pull mapping from cache
          - cache:
              resource: item_mapping
              key: osrs_item_mapping
              operator: get
          # Otherwise store it in the cache
          - catch:
              - http:
                  url: https://prices.runescape.wiki/api/v1/osrs/mapping
                  verb: GET
                  retries: 100
                  headers:
                    # Custom User agent with discord username as recommended by API devs
                    User-Agent: WarpStream - Bento - @rmb938 in Discord
              - mapping: |
                  root = this.map_each(item -> {
                    (item.id.string()): item
                  }).squash()
              - cache:
                  resource: item_mapping
                  key: osrs_item_mapping
                  operator: set
                  ttl: 1h # Using 1h here cause mapping data could change when new items are added
                  value: "${! content() }"
        result_map: |
          root.item_mapping = this

    # Make the http call on 5m endpoint, using branch so we don't overwrite root and clear the root.item_mapping
    - branch:
        processors:
          - http:
              # Query 5m with timestamp
              # Timestamp is rounded to nearest 5 minute interval then goes back in time by 10 mins so we have data.
              # If we don't go back in time by 10 mins, we query the current 5 min interval which may not be populated at query time.
              url: https://prices.runescape.wiki/api/v1/osrs/5m?timestamp=${! this.timestamp }
              verb: GET
              retries: 100
              headers:
                # Custom User agent with discord username as recommended by API devs
                User-Agent: WarpStream - Bento - @rmb938 in Discord
              parallel: false
        result_map: |
          root.data = this.data

    # If there's no data for the 5m period, delete
    # This could be when prices api or osrs is down
    # No point processing data when it's down since all volume will be 0
    # Better to just have missing data points probably
    - mapping: |
        root = this
        if root.data.keys().length() == 0 {
          root = deleted()
        }

    # Copy pricing data and timestamp into the mapping bits
    - mapping: |
        root = {}
        root.data = this.item_mapping.map_each(item -> item.value.merge({"timestamp": this.timestamp, "avgHighPrice": this.data.get(item.key).avgHighPrice, "avgLowPrice": this.data.get(item.key).avgLowPrice, "highPriceVolume": this.data.get(item.key).highPriceVolume, "lowPriceVolume": this.data.get(item.key).lowPriceVolume}))

    # Move all the data items into the root
    - mapping: |
        root = this.data

    # Split the dict into their own messages
    - unarchive:
        format: json_map

    # Fill in volume nulls, if null it means there was no low or high trades in this 5m window
    - mapping: |
        root = this

        if root.lowPriceVolume == null {
          root.lowPriceVolume = 0
        }

        if root.highPriceVolume == null {
          root.highPriceVolume = 0
        }

    ### Start avgLowPrice fixing

    # If avgLowPrice is null try and find last price from cache
    - branch:
        request_map: |
          root = {}
          root.id = id
          if this.avgLowPrice != null {
            root = deleted()
          }
        processors:
          - cache:
              resource: item_price
              key: '${! json("id") }'
              operator: get
          # If not in cache set avgLowPrice to null
          - catch:
              - mapping: |
                  root = {}
                  root.avgLowPrice = null
        result_map: |
          root.avgLowPrice = this.avgLowPrice

    # If avgLowPrice is still null delete it.
    # This means we haven't seen a low trade for this item yet
    - mapping: |
        root = this

        if root.avgLowPrice == null {
          root = deleted()
        }

    ### End avgLowPrice fixing
    ### Start avgHighPrice fixing

    # If avgHighPrice is null try and find last price from cache
    - branch:
        request_map: |
          root = {}
          root.id = id
          if this.avgHighPrice != null {
            root = deleted()
          }
        processors:
          - cache:
              resource: item_price
              key: '${! json("id") }'
              operator: get
          # If not in cache set avgHighPrice to null
          - catch:
              - mapping: |
                  root = {}
                  root.avgHighPrice = null
        result_map: |
          root.avgHighPrice = this.avgHighPrice

    # If avgHighPrice is still null delete it.
    # This means we haven't seen a high trade for this item yet
    - mapping: |
        root = this

        if root.avgHighPrice == null {
          root = deleted()
        }

    ### End avgHighPrice fixing

    # Cache each item info
    - cache:
        resource: item_price
        key: '${! json("id") }'
        operator: set
        value: "${! content() }"

    # Convert unix time to timestamp
    - mapping: |
        root = this
        root.timestamp = this.timestamp.ts_format("2006-01-02T15:04:05Z07:00", "UTC")

output:
  broker:
    pattern: fan_out
    outputs:
      - kafka_franz:
          seed_brokers:
            - warpstream-agent-kafka:9092
          topic: osrs_prices_5m
          key: ${! this.id }
          partitioner: murmur2_hash
          compression: zstd

          # Recommended WarpStream config
          metadata_max_age: 60s
          max_buffered_records: 1000000
          max_message_bytes: 16000000
      #- stdout:
      #    codec: lines

cache_resources:
  - label: timestamp_offset
    redis:
      url: redis://${REDIS_HOST:127.0.0.1}:6379
      prefix: timestamp_offset
  - label: item_mapping
    redis:
      url: redis://${REDIS_HOST:127.0.0.1}:6379
      prefix: item_mapping
  - label: item_price
    redis:
      url: redis://${REDIS_HOST:127.0.0.1}:6379
      prefix: item_price
```

---
layout: center
---
# Terminology
Before we look at the data let's define a few terms

---
layout: center
---
# Volume

Amount of assets traded during a specific period

---
layout: center
---
# Basket of Goods

A fixed representative sample of goods and services that consumers buy

<br />

# CPI - Consumer Price Index

The average change in price of the items in the basket

$\left(\frac{\text{Current Year Basket Cost}}{\text{Base Year Basket Cost}}\right) \times 100$

---
layout: center
---
# Inflation
A sustained increase in the general price of goods and services in an economy, leading to a decline in the purchasing power of money

$\left(\frac{\text{Current Year CPI}-\text{Previous Year CPI}}{\text{CPI Previous Year}}\right) \times 100$

---
layout: center
---
# Market Index
A group of assets that track the performance of a specific market segment

---
---
# Clickhouse and Tableflow

https://clickhouse.com/cloud

All our data is in GCP Object Storage so we are just paying for compute

````md magic-move
```sql
CREATE TABLE 
  osrs_prices_5m 
  ENGINE=IcebergS3(
    'https://storage.googleapis.com/rmb-lab-vancouver_kafka_meetup_2025_nov_19/
      warpstream/_tableflow/warpstream
      +osrs_prices_5m-1d3f32f1-dc10-4fc1-a6f1-a8b136e15582',
    's3api access key', 
    's3api secret key'
  )
```

```sql
SELECT 
    item_id,
    name,
    avg(avgHighPrice),
    sum(highPriceVolume) as highPriceVolume,
    avg(avgLowPrice),
    sum(lowPriceVolume)
FROM 
  osrs_prices_5m
WHERE 
  warpstream.timestamp >= subtractHours(now(), 1) and warpstream.timestamp <= now()
GROUP BY 
  item_id, name
ORDER BY 
  highPriceVolume desc
```
````

---
---
# Clickhouse
We need a query engine for Iceberg

![](/images/Screenshot_20251001_174058.png)

---
---
# WarpStream Tableflow Early Access
It works well in EA but just not for our queries

Primary Keys, Custom Partitioning, and sort order makes all the difference when querying our dataset.

WarpStream Tableflow at this moment only has support for time based partitioning, no sorting, and no primary keys.

Querying our dataset which has 2 billion rows filtering by time and id takes over an hour without custom partitioning and sorting.

Custom partitioning and sorting will be coming sometime during Q4 2025.

---
---
# Clickhouse as Query and Storage Engine
Same Data, Same Process, just not Tableflow for now :)

````md magic-move
```sql {*|1-2|4-16|14-16}
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
```

```sql
SELECT 
    id,
    avg(avgHighPrice),
    sum(highPriceVolume) as highPriceVolume,
    avg(avgLowPrice),
    sum(lowPriceVolume)
FROM 
  osrs_prices_5m
WHERE 
  timestamp >= subtractHours(now(), 1) and timestamp <= now()
GROUP BY 
  id
ORDER BY 
  highPriceVolume desc
```
````

---
---
# Clickhouse as Query and Storage Engine
Same Data, Same Process, just not Tableflow for now :)

```yaml {*|2-21|31-53}{maxHeight:'400px'}
# Common config fields, showing default values
input:
  kafka_franz:
    seed_brokers:
      - warpstream-agent-kafka:9092
      # - localhost:9092
    # client_id: "bento,ws_host_override=localhost"
    topics:
      - osrs_prices_5m
    consumer_group: bento_clickhouse
    auto_replay_nacks: true
    auto_offset_reset: earliest

    # Recommended WarpStream config
    metadata_max_age: 60s
    fetch_max_bytes: 40MiB
    fetch_max_partition_bytes: 40MiB
    fetch_max_wait: "10s"
    batching:
      count: 10000
      period: 15s

pipeline: {}

output:
  broker:
    pattern: fan_out
    outputs:
      # - stdout:
      #     codec: lines
      - sql_insert:
          driver: clickhouse
          dsn: ${CLICKHOUSE_DSN}
          init_verify_conn: true
          table: osrs_prices_5m
          columns:
            [
              id,
              timestamp,
              avgHighPrice,
              highPriceVolume,
              avgLowPrice,
              lowPriceVolume,
            ]
          args_mapping: |
            root = [
              this.id.uint16(),
              this.timestamp.ts_format("2006-01-02 15:04:05"),
              this.avgHighPrice.uint64(),
              this.highPriceVolume.uint64(),
              this.avgLowPrice.uint64(),
              this.lowPriceVolume.uint64()
            ]

```

---
---
# Clickhouse as Query and Storage Engine
Same Data, Same Process, just not Tableflow for now :)

![](/images/Screenshot_20251110_132214.png)

---
layout: two-cols
---
# Grafana
Tables are nice but I want time series

Using Grafana Cloud

https://grafana.com/products/cloud/

<br />

There's an official ClickHouse Data Source Plugin

https://grafana.com/grafana/plugins/grafana-clickhouse-datasource/

::right::

![](/images/Screenshot_20251001_174207.png)

---
---
# Looking at Specific Items
<p />

<v-switch>
  <template #1>
    <h3>Toadflax</h3>
    <img src="/images/Screenshot_20251110_132523.png" alt="Image 1">
  </template>
  <template #2>
    <h3>Sunfire Splinters</h3>
    <img src="/images/Screenshot_20251110_132600.png" alt="Image 2">
  </template>

  <template #3>

  ### Low and High Price Change

  ```sql
  SELECT
      (
          last_value(avgLowPrice) OVER (ORDER BY timestamp ASC)
          - first_value(avgLowPrice) OVER (ORDER BY timestamp ASC)
      ) AS lowPriceChange,
      (
          last_value(avgHighPrice) OVER (ORDER BY timestamp ASC)
          - first_value(avgHighPrice) OVER (ORDER BY timestamp ASC)
      ) AS highPriceChange
  FROM osrs_prices_5m 
  where 
      id = ${item_id} and $__timeFilter(timestamp)
  ORDER by timestamp DESC
  limit 1
  ```

  </template #3>

  <template #4>

  ### Item Volume

  ```sql
  SELECT $__timeInterval(timestamp) as time,
      sum(lowPriceVolume) as averageInstaSellVolume,
      sum(highPriceVolume) as averageInstaBuyVolume
  FROM osrs_prices_5m 
  where 
      id = ${item_id} and $__timeFilter(timestamp)
  group by 
      id, time
  order by time ASC
  LIMIT 10000
  ```

  </template #4>

  <template #5>

  ### Item Price

  ```sql
  SELECT $__timeInterval(timestamp) as time,
      last_value(avgLowPrice) as averageInstaSellPrice, 
      last_value(avgHighPrice) as averageInstaBuyPrice 
  FROM osrs_prices_5m 
  where 
      id = ${item_id} and $__timeFilter(timestamp)
  group by 
      id, time
  order by time ASC
  LIMIT 10000
  ```

  </template #5>

  <template #6>

  ### Item Icon

  ```sql
  SELECT 
    id, 
    concat('https://oldschool.runescape.wiki/images/', replaceAll(icon, ' ', '_')) as iconURL 
  from osrs_items_metadata 
  where id = ${item_id} 
  limit 1;
  ```

  </template #6>
</v-switch>

---
---
# Graphs of CPI and Inflation
Querying this was hard

<v-switch>
  <template #1><img src="/images/Screenshot_20251110_133457.png" alt="Image 1"></template>
  <template #2><img src="/images/Screenshot_20251110_133533.png" alt="Image 2"></template>
  <template #3><img src="/images/Screenshot_20251110_133728.png" alt="Image 2"></template>
  <template #4><img src="/images/Screenshot_20251110_133813.png" alt="Image 2"></template>
</v-switch>

---
---
# Graphs of CPI and Inflation
Querying this was hard

````md magic-move
```sql
-- Base basket of goods price per item
SELECT
  osrs_items_metadata.name,
  avg(avgHighPrice) AS baseAverageHighPrice
FROM osrs_prices_5m 
JOIN osrs_items_metadata on osrs_prices_5m.id = osrs_items_metadata.id
WHERE 
  id IN (${basketOfGoods})
  AND timestamp >= toDateTime(1614556800) AND timestamp < toDateTime(1617235199)
GROUP BY
  osrs_items_metadata.name
ORDER BY
  baseAverageHighPrice DESC
```
```sql
-- Base basket of goods price total
select
    sum(baseAverageHighPrice)
FROM
    (
      SELECT 
        id,
        avg(avgHighPrice) as baseAverageHighPrice
      FROM osrs_prices_5m 
      WHERE 
        id IN (${basketOfGoods})
      AND timestamp >= toDateTime(1614556800) AND timestamp < toDateTime(1617235199)
      GROUP BY
        id
    )
```
```sql
-- CPI Calculation
SELECT 
  time,
  ((sum(avgHighPrice1)/${basketOfGoodsTotalBase})*100) as CPI
FROM
  (SELECT
    $__timeInterval(timestamp) as time,
    id,
    avg(avgHighPrice) as avgHighPrice1
  FROM osrs_prices_5m 
  WHERE 
    id IN (${basketOfGoods})
    AND $__timeFilter(timestamp)
  GROUP BY
    id, time
  LIMIT 100000
  )
group by time
order by time asc
```
```sql
-- Inflation Percent Calculation
SELECT 
  time,
  (((sum(avgHighPrice1)/${basketOfGoodsTotalBase})*100)-100) as InflationPercent
FROM
  (SELECT
    $__timeInterval(timestamp) as time,
    id,
    avg(avgHighPrice) as avgHighPrice1
  FROM osrs_prices_5m 
  WHERE 
    id IN (${basketOfGoods})
    AND $__timeFilter(timestamp)
  GROUP BY
    id, time
  LIMIT 100000
  )
group by time
order by time asc
```
````

---
---
# Graphs of Indexes
<p />

<v-switch>
  <template #1><img src="/images/Screenshot_20251110_134703.png" alt="Image 1"></template>
  <template #2><img src="/images/Screenshot_20251110_134711.png" alt="Image 2"></template>
</v-switch>

---
---
# Index Queries
<p />

````md magic-move
```sql
-- Index Price
SELECT
    time,
    sum(averageInstaBuyPrice)/${herbIndexDivisor} as index
FROM
    (SELECT $__timeInterval(timestamp) as time,
        id,
        avg(avgHighPrice) as averageInstaBuyPrice
    FROM osrs_prices_5m 
    where 
        id IN [${herbIndexItems}]
        and $__timeFilter(timestamp)
    group by 
        id, time
    order by time ASC
    LIMIT 100000)
group by time
```
```sql
-- Index Individual Item Price
SELECT $__timeInterval(timestamp) as time,
    osrs_items_metadata.name,
    avg(avgHighPrice) as averageInstaBuyPrice
FROM osrs_prices_5m 
JOIN osrs_items_metadata on osrs_prices_5m.id = osrs_items_metadata.id
where 
    id IN [${herbIndexItems}]
    and $__timeFilter(timestamp)
group by 
    osrs_items_metadata.name, time
order by time ASC
LIMIT 100000
```
```sql
-- Index Individual Item Price Table
SELECT
  name,
  first_value(averageInstaBuyPrice) as Last,
  min(averageInstaBuyPrice) as Min,
  max(averageInstaBuyPrice) as Max,
  avg(averageInstaBuyPrice) as Mean,
  first_value(avgHighPriceVolume) as Volume
FROM
  (SELECT $__timeInterval(timestamp) as time,
      osrs_items_metadata.name as name,
      avg(avgHighPrice) as averageInstaBuyPrice,
      avg(highPriceVolume) as avgHighPriceVolume
  FROM osrs_prices_5m 
  JOIN osrs_items_metadata on osrs_prices_5m.id = osrs_items_metadata.id
  where 
      id IN [${herbIndexItems}]
      and timestamp >= toDateTime(subtractHours(now(), 3)) AND timestamp <= toDateTime(now())
  group by 
      osrs_items_metadata.name, time
  order by time DESC
  LIMIT 100000)
group by name
```
````

---
---
# Computing Indexes in Bento
Changing Index Makeup and Divisor over time is better with stream processing

```yaml {*|7-11|14-20|22-32|34-41|43-56|58-63}{maxHeight:'400px'}
pipeline:
  processors:
    # Herb Index
    - mapping: |
        root = this

        # Monday, March 8, 2021 8:00:00 AM
        let start_time = 1615190400

        let message_time = this.timestamp.ts_unix()
        let message_offset = $message_time - $start_time

        # Lists are offsets after start_time
        let herb_index_items = {
          # Start
          0.string(): [249, 251, 253, 255, 257, 2998, 259, 261, 263, 3000, 265, 2481, 267, 269],

          # Huasca released September 25 2024
          # Offset is Monday, September 30, 2024 12:00:00 AM
          (1727654400 - $start_time).string(): [30097, 249, 251, 253, 255, 257, 2998, 259, 261, 263, 3000, 265, 2481, 267, 269],
        }

        # List of divisors, offsets don't need to line up with items map
        # These change whenever a new item is added or prices drastically change
        let herb_index_divisor = {
          # Start
          0.string(): 20,

          # Huasca released September 25 2024
          # Offset is Monday, September 30, 2024 12:00:00 AM
          (1727654400 - $start_time).string(): 28,
        }

        # Collect all the offsets and lists into a list
        # Filter for offsets that are <= message offset
        # Sort by increasing offset
        # Get the last one
        let matching_item_list = $herb_index_items.keys().map_each(key -> {
            "offset": key.number(),
            "indices": $herb_index_items.get(key)
        }).filter(item -> item.offset <= $message_offset).sort_by(ele -> ele.offset).index(-1).get("indices")

        if $matching_item_list.contains(this.id) {
          root.index_name = "herb"

          # Collect all the offsets and divisors into a list
          # Filter for offsets that are <= message offset
          # Sort by increasing offset
          # Get the last one
          let divisor = $herb_index_divisor.keys().map_each(key -> {
              "offset": key.number(),
              "divisor": $herb_index_divisor.get(key).number()
          }).filter(item -> item.offset <= $message_offset).sort_by(ele -> ele.offset).index(-1).get("divisor")

          root.divisor = $divisor
        }

    # Remove any items that never got assigned to an index
    - mapping: |
        root = this
        if root.index_name == null {
          root = deleted()
        }

output:
  broker:
    pattern: fan_out
    outputs:
      - stdout:
          codec: lines
```

---
---
# Running this at home!
My own cloud instead of someone else's

<v-switch>
  <template #1><img src="/images/Screenshot_20251110_155229.png" alt="Image 1"></template>
  <template #2><img src="/images/Screenshot_20251110_155305.png" alt="Image 2"></template>
  <template #3>

  ```bash
  WARPSTREAM_DEFAULT_VIRTUAL_CLUSTER_ID=vci_...
  WARPSTREAM_AGENT_KEY=aks_...
  WARPSTREAM_BUCKET_URL=s3://vancouver-kafka-meetup-2025-nov-19-3 \
    ?region=us-homelab1&s3ForcePathStyle=true&endpoint=https://garage-s3.haproxy.us-homelab1.hl.rmb938.me
  AWS_ACCESS_KEY_ID=G...
  AWS_SECRET_ACCESS_KEY=0...
  WARPSTREAM_REGION=us-central1
  GOMAXPROCS=2
  WARPSTREAM_AVAILABILITY_ZONE=us-homelab1
  ```

  </template>
  <template #4>

  ```bash
  % podman ps
  CONTAINER ID  IMAGE                                                COMMAND               CREATED      STATUS                PORTS                   NAMES
  799178257cf1  localhost/2025-nov-19_warpstream-agent-kafka:latest  agent                 8 hours ago  Up 8 hours            0.0.0.0:9092->9092/tcp  2025-nov-19_warpstream-agent-kafka_1
  512814d164a0  localhost/2025-nov-19_bento-clickhouse:latest        --log.level debug     8 hours ago  Up 8 hours                                    2025-nov-19_bento-clickhouse_1
  1a69e7704e18  docker.io/library/redis:8                            redis-server --sa...  8 hours ago  Up 8 hours (healthy)  0.0.0.0:6379->6379/tcp  2025-nov-19_redis_1
  84a282e159c5  ghcr.io/warpstreamlabs/bento:latest                  --log.level info      8 hours ago  Up 8 hours                                    2025-nov-19_bento_1
  ```

  </template>
  <template #5>

  ```bash
  FROM smallstep/step-ca as CA

  RUN step ca root /tmp/smallstep-homelab-prod.crt --ca-url https://step-ca.us-homelab1.hl.rmb938.me --fingerprint 111301fb085dfc83f5390c0be68df3d68f5584853df0cb4c442383c33f2bd83a

  FROM public.ecr.aws/warpstream-labs/warpstream_agent:latest

  COPY --from=CA /tmp/smallstep-homelab-prod.crt /usr/local/share/ca-certificates/smallstep-homelab-prod.crt

  RUN cat /usr/local/share/ca-certificates/smallstep-homelab-prod.crt >> /etc/ssl/certs/ca-certificates.crt
  ```

  </template>
</v-switch>

---
---
# Learnings and What's Next
Graphing Data is hard

Source API has a low rate limit

Took about 8 days to load 4 years worth of data

Finishing proving out using Bento to compute indexes

Add game updates and player events to the graphs using Grafana Annotations

Data as a Service, figure out a cheap way to provide this data and queries to players.

---
layout: end
---

# Thanks You!
Q & A

Slides, Bento, Tableflow, & Clickhouse Configs Available at 

https://github.com/rmb938/presentations

OR

<style>
.image-container {
  /* Enables Flexbox layout for direct children */
  display: flex;

  /* Centers the images horizontally within the container (optional) */
  justify-content: center;

  /* Adds some space around the images (optional) */
  gap: 300px;

  /* Sets a maximum width for the container (optional) */
  max-width: 800px;
  margin: 0 auto;
}

.image-container img {
  /* Makes the images take up equal space */
  /* flex: 1; */

  /* Ensures the images don't exceed the container's width */
  max-width: 100%;

  /* Ensures images maintain their aspect ratio */
  height: auto;
}
</style>

<div class="image-container">
  <img src="/images/qr-code-slides.svg" width="220"/>
</div>
