---
# You can also start simply with 'default'
theme: seriph
# random image from a curated Unsplash collection by Anthony
# like them? see https://unsplash.com/collections/94734566/slidev
background: /images/mario-verduzco-xSdFf1Lcx6o-unsplash.jpg
# some information about your slides (markdown enabled)
title: Embedding a Stream Processor Inside Your Go Application
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
# Embedding a Stream Processor Inside Your Go Application

<!--
The last comment block of each slide will be treated as slide notes. It will be visible and editable in Presenter Mode along with the slide. [Read more in the docs](https://sli.dev/guide/syntax.html#notes)
-->

---
layout: profile
image: /images/saKWC52.jpg
---

# Riley Belgrave
Staff Engineer @ Confluent - WarpStream

~20 Years of technical experience, self-taught with Java in 2005

Ordered my first Ubuntu liveCD soon after and forever got stuck in vim

Previously helped run the infrastructure automation for a few popular Minecraft servers

Deployed the largest deployment of baremetal K8s in the mid-west, probably the US, for a few months in early 2019.

<br />

<carbon-logo-github /> <a href="https://github.com/rmb938" style="color: #6b7280">rmb938</a>
<carbon-logo-linkedin /> <a href="https://www.linkedin.com/in/rbelgrave/" style="color: #6b7280">rbelgrave</a>
<circle-flags-lgbt-progress /> <span style="color: #6b7280">they/them</span>

---
---
# What is Stream Processing
<br />

Stream processing enables continuous data ingestion, streaming, filtering, and transformation as events happen in real time

Once processed, the data can be passed off to an application, data store, or another stream processing engine to provide actionable insights quickly

<br />

Typically used for real-time analytics, fruad detections, recommendation systems, and more

<br />

Various tools used include Kafka Streams, Apache Flink, and Spark Stream

---
layout: two-cols
---

# Stream Processing Made Easy
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
---
# Bento Custom Plugins
Faker Example

```go {*|18-24|26-46|50-58|70-72|74-93|93-102|105-107}{maxHeight:'400px'}
package plugins

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pioz/faker"
	"github.com/warpstreamlabs/bento/public/service"
)

const (
	fakerInputFieldUsername  = "username"
	fakerInputFieldAge       = "age"
	fakerInputFieldBatchSize = "batch_size"
)

func fakerInputSpec() *service.ConfigSpec {
	return service.NewConfigSpec().Fields(
		service.NewBoolField(fakerInputFieldUsername),
		service.NewBoolField(fakerInputFieldAge),
		service.NewIntField(fakerInputFieldBatchSize).Default(1),
	)
}

func newFakerInput(conf *service.ParsedConfig, opts []int) (service.BatchInput, error) {
	name, err := conf.FieldBool(fakerInputFieldUsername)
	if err != nil {
		return nil, err
	}

	age, err := conf.FieldBool(fakerInputFieldAge)
	if err != nil {
		return nil, err
	}

	batchSize, err := conf.FieldInt(fakerInputFieldBatchSize)
	if err != nil {
		return nil, err
	}

	if opts != nil {
		seed := opts[0]
		// AutoRetryNacksBatched will reattempt failed messages because this input won't be able to handle nacks.
		return service.AutoRetryNacksBatched(&fakerInput{name: name, age: age, batchSize: batchSize, seed: seed}), err
	}
	return service.AutoRetryNacksBatched(&fakerInput{name: name, age: age, batchSize: batchSize}), err
}

func init() {
	err := service.RegisterBatchInput("faker", fakerInputSpec(),
		func(pConf *service.ParsedConfig, res *service.Resources) (service.BatchInput, error) {
			return newFakerInput(pConf, nil)
		})
	if err != nil {
		panic(err)
	}
}

// ------------------------------------------------------------------------------------

type fakerInput struct {
	name      bool
	age       bool
	batchSize int

	seed int
}

func (f *fakerInput) Connect(ctx context.Context) error {
	return nil
}

func (f *fakerInput) ReadBatch(ctx context.Context) (service.MessageBatch, service.AckFunc, error) {
	batch := make(service.MessageBatch, 0, f.batchSize)

	if f.seed != 0 {
		faker.SetSeed(int64(f.seed))
	}

	for range f.batchSize {

		fd := make(map[string]any)

		if f.name {
			fd[fakerInputFieldUsername] = faker.Username()
		}

		if f.age {
			fd[fakerInputFieldAge] = faker.IntInRange(0, 95)
		}

		jsonData, err := json.Marshal(fd)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshall JSON: %w", err)
		}

		msg := service.NewMessage(jsonData)
		batch = append(batch, msg)
	}

	return batch, func(context.Context, error) error { return nil }, nil
}

func (f *fakerInput) Close(ctx context.Context) error {
	return nil
}
```

---
---
# Bento Custom Plugins
Something more real

```go {*|28-37|21-26|46-48|71-83|104-120|131-134}{maxHeight:'400px'}
package github_events_archive

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/google/go-github/v72/github"
	"github.com/warpstreamlabs/bento/public/service"
)

func inputSpec() *service.ConfigSpec {
	return service.NewConfigSpec().Fields(service.NewAutoRetryNacksToggleField())
}

func newInput(conf *service.ParsedConfig) (service.BatchInput, error) {
	return service.AutoRetryNacksBatchedToggled(conf, &githubArchiveInput{
		// TODO: hard coded for now
		nextArchiveTime: time.Now().Add(-1 * 7 * 24 * time.Hour).UTC(),
	})
}

func init() {
	err := service.RegisterBatchInput("github_events_archive", inputSpec(),
		func(pConf *service.ParsedConfig, res *service.Resources) (service.BatchInput, error) {
			return newInput(pConf)
		})
	if err != nil {
		panic(err)
	}

}

type githubArchiveInput struct {
	nextArchiveTime time.Time

	trackedBatchesLock sync.Mutex
	trackedBatches     []string
}

func (i *githubArchiveInput) Connect(ctx context.Context) error {
	return nil
}

func (i *githubArchiveInput) ReadBatch(ctx context.Context) (service.MessageBatch, service.AckFunc, error) {
	for {
		i.trackedBatchesLock.Lock()
		if len(i.trackedBatches) == 0 {
			i.trackedBatchesLock.Unlock()
			break
		}
		fmt.Println("There is back preasure, waiting...")
		i.trackedBatchesLock.Unlock()
		time.Sleep(5 * time.Second)
	}

	fmt.Println("No backpreasure starting")

	i.trackedBatchesLock.Lock()
	defer i.trackedBatchesLock.Unlock()

	fmt.Println("Got the lock")

	batch := make(service.MessageBatch, 0)

	archiveURL := fmt.Sprintf("https://data.gharchive.org/%s-%d.json.gz", i.nextArchiveTime.Format("2006-01-02"), i.nextArchiveTime.Hour())
	fmt.Printf("Downloading archive from %s\n", archiveURL)
	resp, err := http.Get(archiveURL)
	if err != nil {
		return nil, nil, fmt.Errorf("error making http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// archive not found, so sleep 1 min
		time.Sleep(1 * time.Minute)
		return batch, func(context.Context, error) error { return nil }, nil
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, nil, fmt.Errorf("error reading body during error: %w", err)
		}
		return nil, nil, fmt.Errorf("error response from server (status code %d): %s", resp.StatusCode, string(bodyBytes))
	}

	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating gzip reader: %w", err)
	}
	defer gzipReader.Close()

	jsonDecoder := json.NewDecoder(gzipReader)

	var startID *string
	var endID string

	for {
		var event github.Event
		if err := jsonDecoder.Decode(&event); err != nil {
			if err == io.EOF {
				break // End of stream
			}

			return nil, nil, fmt.Errorf("error decoding json object: %w", err)
		}

		jsonData, err := json.Marshal(event)
		if err != nil {
			return nil, nil, fmt.Errorf("error json marshaling github event %d: %w", event.ID, err)
		}

		msg := service.NewMessage(jsonData)
		batch = append(batch, msg)

		if startID == nil {
			startID = event.ID
		}

		endID = *event.ID

		i.trackedBatches = append(i.trackedBatches, *event.ID)
	}

	i.nextArchiveTime = i.nextArchiveTime.Add(1 * time.Hour)

	fmt.Printf("Sending batch with size %d Start ID: %s End ID: %s\n", len(batch), *startID, endID)
	return batch, func(ctx context.Context, err error) error {
		fmt.Printf("starting ack func: %v\n", err)

		i.trackedBatchesLock.Lock()
		defer i.trackedBatchesLock.Unlock()

		i.trackedBatches = i.trackedBatches[:0]

		fmt.Println("finished ack func")
		return nil
	}, nil
}

func (i *githubArchiveInput) Close(ctx context.Context) error {
	return nil
}
```

---
---
# Import and run bneto with the plugin
Just a simple main.go

```go
package main

import (
	"context"

	// import all components with:
	_ "github.com/warpstreamlabs/bento/public/components/all"

	"github.com/warpstreamlabs/bento/public/service"
	// import your plugins:

	_ "github.com/rmb938/.../github_events_archive"
)

func main() {
	// RunCLI accepts a number of optional functions:
	// https://pkg.go.dev/github.com/warpstreamlabs/bento/public/service#CLIOptFunc
	service.RunCLI(context.Background())
}
```

---
layout: two-cols
---
# Running Bento
As simple as a `go run`

```bash {*|1|6|9-10}
$ go run main.go -c bento-config-gh.yaml 
INFO Running main config from specified file       @service=bento bento_version=v1.7.1 path=bento-config-gh.yaml
INFO Listening for HTTP requests at: http://0.0.0.0:4195  @service=bento
INFO Launching a Bento instance, use CTRL+C to close  @service=bento
INFO Output type kafka_franz is now active         @service=bento label="" path=root.output
INFO Input type github_events_archive is now active  @service=bento label="" path=root.input
No backpreasure starting
Got the lock
Downloading archive from https://data.gharchive.org/2025-05-14-20.json.gz
Sending batch with size 233916 Start ID: 49730103801 End ID: 49732035425
```

::right::

# Configuration
Simple, Boring, Configuration

```yaml  {*|2|5-9}
input:
  github_events_archive: {}

output:
  kafka_franz:
    seed_brokers:
      - localhost:9092
    topic: github_events
    partitioner: uniform_bytes
```

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
layout: two-cols
---
# Managed Data Pipelines
Stream Processing Made Simple

ETL and stream processing from within your WarpStream Agents and cloud account. 

No additional infrastructure needed. 

You own the data pipeline end to end. 

Raw data never leaves your account.

::right::

# Powered By Bento
<br />

![](/images/bento_food.svg)

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
  <img src="/images/qr-code-warpstream-mdp.svg" width="200"/>
</div>

---
---
# Managing Bento at Scale
Easy for a few pipelines, hard when you have many

Bento is completely stateless so it's easy to just spin up a few containers in Kubernetes and call it a day

Updating pipeline configuration is a simple configmap update and then redeploy bento

<br />

This is simple for a few pipelines and deploymets, but becomes and automation and resource allocation issue once you have a few dozen.

<br />

Each Bento deployment is typically 1:1 per pipeline and resource sharing while idle can be difficult.

Kubernetes overprovisioning and right sizing can help improve resource utilization but it's not perfect.

---
---
# Building Managed Data Pipelines
Building a streaming processor platform without vendor lock-in

## Goals

Allow users to run Bento Pipelines with with their existing WarpStream Agent architecture

Allow easy pipeline updates without having to redeploy or reconfigure agents

Use WarpStream's distributed job load balancing and Kubernetes autoscaling and maximize resource utlization

Use the open-source MIT licensed Bento to keep maximum compatibility

---
---
# How to run Bento inside a Golang App
WarpStream agents and Bento are built using Go

```go
env := service.NewEnvironment()
blobEnv := bloblang.NewEnvironment()

env.UseBloglangEnvironment(blobEnv)

builder := env.NewStreamBuilder()

builder.SetYAML("stuff")

stream, err := builder.Build()

stream.Run(ctx)
```

---
---
# Keeping things Secure
Limiting access to the filesystem, environment, and other system resources

This is all built into Bento so if you want to do this yourself it's simple

```go
// Protect the file system or plugin a custom file system like a blob store
env.UseFS(service.NewFS(&fakeFS{}))

// Create a custom importer to prevent importing of custom code outside of the pipeline environment
blobEnv = blobEnv.WithCustomImporter(func(_ string) ([]byte, error) {
  return []byte("map fake {}"), nil
})

// Prevent access to environment variables
builder.SetEnvVarLookupFunc(func(_ string) (string, bool) {
  // TODO: Allow environment variables matching certain patterns
  return "", false
})
```

---
---
# Putting it all Together Yourself
You don't need WarpStream to build something similar

<<< @/snippets/code/main.go go {*|181-202|202-214|216-235|137-157|159-178}{maxHeight:'400px'}

---
---
# Using WarpStream Managed Data Pipelines
It's free with no extra cost when using WarpStream for Kafka

![](</images/Screenshot 2026-04-14 at 3.00.51 PM.png>)

---
---
# Using WarpStream Managed Data Pipelines
It's free with no extra cost when using WarpStream for Kafka

![](</images/Screenshot 2026-04-14 at 3.00.58 PM.png>)

---
---
# Using WarpStream Managed Data Pipelines
It's free with no extra cost when using WarpStream for Kafka

![](</images/Screenshot 2026-04-14 at 3.01.14 PM.png>)

---
---
# Using WarpStream Managed Data Pipelines
It's free with no extra cost when using WarpStream for Kafka

![](</images/Screenshot 2026-04-14 at 3.01.23 PM.png>)

---
---
# Using WarpStream Managed Data Pipelines
It's free with no extra cost when using WarpStream for Kafka

![](</images/Screenshot 2026-04-14 at 3.01.39 PM.png>)

---
---
# Using WarpStream Managed Data Pipelines
It's free with no extra cost when using WarpStream for Kafka

![](</images/Screenshot 2026-04-14 at 3.01.47 PM.png>)

---
---
# Using WarpStream Managed Data Pipelines
It's free with no extra cost when using WarpStream for Kafka

![](</images/Screenshot 2026-04-14 at 3.01.57 PM.png>)

---
---
# Using WarpStream Managed Data Pipelines
It's free with no extra cost when using WarpStream for Kafka

![](</images/Screenshot 2026-04-14 at 3.02.09 PM.png>)

---
---
# WarpStream Managed Data Pipelines Helpers
Connecting to WarpStream and Parallelism Management

````md magic-move
```yaml
input:
    kafka_franz:
        seed_brokers: ["localhost:9092"]
        topics: ["test_topic"]

    processors:
        - mapping: "root = content().capitalize()"

output:
    kafka_franz:
        seed_brokers: ["localhost:9092"]
        topic: "test_topic_capitalized"
```
```yaml
input:
    kafka_franz_warpstream:
        topics: ["test_topic"]

    processors:
        - mapping: "root = content().capitalize()"

output:
    kafka_franz_warpstream:
        topic: "test_topic_capitalized"
```
```yaml
warpstream:
  scheduling:
    run_every: 1s
  cluster_concurrency_target: 1
  pipeline_group: my-agent-group
```
````

---
---
# WarpStream Tableflow uses Bento Internally
bloblang for table transforms

````md magic-move
```json
{
    "schema": {...},
    "payload": {
    	"op": "u",
    	"source": {
    		...
    	},
    	"ts_ms" : "...",
    	"ts_us" : "...",
    	"ts_ns" : "...",
    	"before" : {
    		"field1" : "oldvalue1",
    		"field2" : "oldvalue2"
    	},
    	"after" : {
    		"field1" : "newvalue1",
    		"field2" : "newvalue2"
    	}
	}
}
```
```yaml
source_format: json
input_schema: |
  {
    "type": "object",
    "properties": {
      "payload": {
        "type": "object",
        "properties": {
          "after": {
            "type": "object",
            "properties": {
              "field1": { "type": "string" },
              "field2": { "type": "string" }
            }
          }
        }
      }
    }
  }
```
```yaml {*|8-14|3-6}
source_format: json
transforms:
  - transform_type: bento
    transform: |
      root.field1 = this.payload.after.field1
      root.field2 = this.payload.after.field2
input_schema: |
  {
    "type": "object",
    "properties": {
      "field1": { "type": "string" },
      "field2": { "type": "string" }
    }
  }
```
```go
func ParseDataLakeBentoTransform(transform string) (*bloblang.Executor, error) {
	if err := ValidateDataLakeBentoTransform(transform); err != nil {
		return nil, err
	}

	return bloblang.NewEnvironment().OnlyPure().Parse(transform)
}

....

for i, executor := range bentoExecutors {
  var bentoAny = any(structured)
  err := executor.Overlay(structured, &bentoAny)
  if err != nil {
    e := fmt.Errorf("failed to apply bento transform at index %d: %w", i, err)
    return e
  }

  ...
}
```
````

---
layout: end
---

# Thanks You!
Q & A

Slides, Bento Configs, Plugin Code, Bento Pipeline Automation Available at 

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
