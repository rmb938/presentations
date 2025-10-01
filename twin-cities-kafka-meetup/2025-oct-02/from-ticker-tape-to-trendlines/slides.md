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

Deployed the largest deployment of baremetal K8s in the mid-west, probably the US for a few months in early 2019.

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

All using streaming processing technologies, Apache Iceberg, and Grafana

::right::

<img src="/images/arthur-a-9rz5x8LGBb8-unsplash.jpg" class="w-100"/>

---
layout: two-cols
---

# Stream Processing the Events
Using Bento

![](/images/bento_boringly_easy.png)

Fancy stream processing made operationally mundane

https://warpstreamlabs.github.io/bento/

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
# Apache Iceberg
Spark and Compaction

![](/images/iceberg_without_warpstream.png)

---
---
# Apache Iceberg
Tableflow

TODO: nice image with warpstream

---
layout: two-cols
---
# Other BYOC Products

* Requires access to your VPC to deploy clusters.

* Requires expansive cross-IAM roles to manage the cluster remotely.

* Requires expensive and complex VPC peering.

 * Remote access is required for support, and elevated permissions can raise security issues. Break glass enables root access.

::right::

# WarpStream's BYOC

* Data / Metadata split enables WarpStream's control plane to function with no access to your VPC or object storage buckets.
* Only metadata is transferred between your VPC and WarpStream's control plane.
* It is impossible for WarpStream to access your data under any circumstances.

---
layout: full
---

<img src="/images/warpstream_logo.svg" width="1000"/>

## [warpstream.com](https://www.warpstream.com/)

<br />
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
  <img src="/images/qr-code-warpstream.svg" width="220"/>
  TODO: QR code to tableflow docs
  <img src="" width="220"/>
</div>


---
---
# Loading Historical Data With Bento

````md magic-move
```yaml
input:
  generate:
    interval: 1s # Max 1s so API devs don't get mad
    batch_size: 1
    auto_replay_nacks: true
    # Set the timestamp in the message, starting at 1615190400 (Monday, March 8, 2021 8:00:00 AM GMT) 
    # and add 300 seconds each time
    count: 479378 # Maximum this many 5 min intervals since 1615190400, as of Sept 27th 2025
    mapping: |
      root = {}
```
```yaml {*|6|14}
# Load the current timestamp offset from cache
# Make the http call on 5m endpoint, using branch so we don't overwrite root and clear the root.item_mapping
- branch:
    processors:
      - http:
          url: https://example.com/api/v1/trades/5m?timestamp=${! this.timestamp }
          verb: GET
          retries: 100
          headers:
            # Custom User agent with discord username as recommended by API devs
            User-Agent: WarpStream - Bento - @rmb938 in Discord
          parallel: false
    result_map: |
      root.data = this.data
```

```yaml {*|1-6|8-10|11-14}
# Copy the timestamp into all the data items
- mapping: |
    root.data = this.data.map_each(item -> item.value.merge({
      "id": item.key.number(), 
      "timestamp": this.timestamp, 
    }))

# Move all the data items into the root
- mapping: |
    root = this.data

# Split the dict into their own messages
- unarchive:
    format: json_map
```
```yaml {*|5-16}
output:
  broker:
    pattern: fan_out
    outputs:
      - kafka_franz:
          seed_brokers:
            - warpstream-agent-kafka:9092
          topic: prices_5m
          key: ${! this.id }
          partitioner: murmur2_hash
          compression: zstd

          # Recommended WarpStream config
          metadata_max_age: 60s
          max_buffered_records: 1000000
          max_message_bytes: 16000000
      - stdout:
          codec: lines
```
````

---
---
# Loading Realtime Data With Bento

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
```yaml {*|7}
- branch:
    processors:
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
    result_map: |
      root.data = this.data
      root.timestamp = this.timestamp
```
````

---
---
# Loading All the Data
Lots of Records

![](/images/Screenshot_20250930_193810.png)

---
---
# Tableflow Configuration
Like Bento it's YAML

```yaml {*|4-5|6|8-9|10|12|13-22}
source_clusters:
  - name: "warpstream"
    bootstrap_brokers:
      - hostname: "warpstream-agent-kafka"
        port: 9092
destination_bucket_url: "gs://rmb-lab-twin_cities_kafka_meetup_2025_oct_02"
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
---
# Tableflow Deploy, Stats, & Partitions
Kafka -> Tableflow, no Java, no complex systems

```yaml
warpstream-agent-tableflow:
  image: public.ecr.aws/warpstream-labs/warpstream_agent:latest
  command:
    - agent
  env:
    WARPSTREAM_DEFAULT_VIRTUAL_CLUSTER_ID: "vci_dl_foo"
    WARPSTREAM_AGENT_KEY: "***********"
  secrets:
    - gcp_service_account_key
```

TODO: picture here of table flow throughput

then another pic of parittion file picker

then another pic of the partition stats

picture of bucket stats in gcs

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

todo: picture of https://prices.runescape.wiki/osrs/

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
layout: center
---
# Terminology
Before we look at the data let's define a few terms

---
layout: center
---
# Volume

stuff

---
layout: center
---
# Basket of Goods

stuff

<br />

# CPI - Consumer Price Index

stuff

---
layout: center
---
# Inflation

stuff

---
layout: center
---
# Market Index

stuff

---
layout: center
---
# Market Manipulation

stuff

---
---
# Clickhouse
We need a query engine for Iceberg

https://clickhouse.com/cloud

All our data is in Object Store so we are just paying for compute

```sql
table creation here
```

```sql
select * from table limiting now
```

nice picture of table of data

---
---
# Grafana
Tables are nice but I want time series

https://grafana.com/products/cloud/

https://grafana.com/grafana/plugins/grafana-clickhouse-datasource/

picture of datasource setup

---
---
# Graphs picking specific items

---
---
# Graphs of CPI and Inflation

---
---
# Graphs of Indexes

---
---
# Learnings and What's Next

talk about loading data is slow

talk about accidently putting historical and realtime in the same topic

talk about using bento to compute cpi, inflation and indexes in the stream proocessor

talk about wanting to add game updates, player events, ect.. to the data to see
what happens during interesting moments


---
layout: end
---

# Thanks You!
Q & A

Slides Available at 

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
