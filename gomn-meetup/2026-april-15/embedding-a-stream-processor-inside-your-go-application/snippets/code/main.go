package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/warpstreamlabs/bento/public/bloblang"
	"github.com/warpstreamlabs/bento/public/service"

	_ "github.com/warpstreamlabs/bento/public/components/all"
)

var pipelineTemplates = map[string]string{
	"counters": `
input:
  generate:
    count: 0
    interval: 1s
    mapping: |
      root.count = counter()
      root.pipeline = "counters"
output:
  label: ""
  stdout:
    codec: lines
`,
	"user-events": `
input:
  generate:
    count: 0
    interval: 2s
    mapping: |
      root.event = if random_int() % 3 == 0 { "login" } else if random_int() % 3 == 1 { "logout" } else { "purchase" }
      root.user_id = random_int(min: 1, max: 1000)
      root.pipeline = "user-events"
output:
  label: ""
  stdout:
    codec: lines
`,
	"metrics": `
input:
  generate:
    count: 0
    interval: 3s
    mapping: |
      root.cpu_percent = random_int(min: 0, max: 100)
      root.memory_mb = random_int(min: 256, max: 16384)
      root.disk_percent = random_int(min: 0, max: 100)
      root.pipeline = "metrics"
output:
  label: ""
  stdout:
    codec: lines
`,
	"logs": `
input:
  generate:
    count: 0
    interval: 1500ms
    mapping: |
      root.level = if random_int() % 4 == 0 { "ERROR" } else if random_int() % 4 == 1 { "WARN" } else { "INFO" }
      root.message = "application log entry"
      root.service = if random_int() % 2 == 0 { "api" } else { "worker" }
      root.pipeline = "logs"
output:
  label: ""
  stdout:
    codec: lines
`,
	"orders": `
input:
  generate:
    count: 0
    interval: 2500ms
    mapping: |
      root.order_id = "ord_" + random_int(min: 10000, max: 99999).string()
      root.amount = random_int(min: 1, max: 500)
      root.status = if random_int() % 2 == 0 { "pending" } else { "completed" }
      root.pipeline = "orders"
output:
  label: ""
  stdout:
    codec: lines
`,
	"heartbeats": `
input:
  generate:
    count: 0
    interval: 5s
    mapping: |
      root.status = "ok"
      root.pipeline = "heartbeats"
output:
  label: ""
  stdout:
    codec: lines
`,
}

// getPipelines returns a random subset of 1-5 pipelines from the template pool.
func getPipelines() map[string]string {
	names := make([]string, 0, len(pipelineTemplates))
	for name := range pipelineTemplates {
		names = append(names, name)
	}

	rand.Shuffle(len(names), func(i, j int) {
		names[i], names[j] = names[j], names[i]
	})

	count := 1 + rand.Intn(5)
	if count > len(names) {
		count = len(names)
	}

	result := make(map[string]string, count)
	for _, name := range names[:count] {
		result[name] = pipelineTemplates[name]
	}
	return result
}

type runningPipeline struct {
	cancel context.CancelFunc
	done   chan struct{}
}

func startPipeline(ctx context.Context, name, yaml string) (*runningPipeline, error) {
	env := service.NewEnvironment()
	blobEnv := bloblang.NewEnvironment()

	env.UseBloblangEnvironment(blobEnv)

	builder := env.NewStreamBuilder()

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

	if err := builder.SetYAML(yaml); err != nil {
		return nil, fmt.Errorf("error setting YAML: %w", err)
	}

	stream, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("error building stream: %w", err)
	}

	pCtx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})

	go func() {
		defer close(done)
		if err := stream.Run(pCtx); err != nil && pCtx.Err() == nil {
			log.Printf("pipeline %s exited with error: %v", name, err)
		}
	}()

	return &runningPipeline{cancel: cancel, done: done}, nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	running := make(map[string]*runningPipeline)

	poll := func() {
		desired := getPipelines()

		// Stop pipelines that are no longer in the desired set
		for name, p := range running {
			if _, ok := desired[name]; !ok {
				log.Printf("stopping pipeline: %s", name)
				p.cancel()
				<-p.done
				delete(running, name)
				log.Printf("stopped pipeline: %s", name)
			}
		}

		// Start pipelines that are new in the desired set
		for name, yaml := range desired {
			if _, ok := running[name]; !ok {
				log.Printf("starting pipeline: %s", name)
				p, err := startPipeline(ctx, name, yaml)
				if err != nil {
					log.Printf("error starting pipeline %s: %v", name, err)
					continue
				}
				running[name] = p
				log.Printf("started pipeline: %s", name)
			}
		}
	}

	poll() // initial poll

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down all pipelines")
			for name, p := range running {
				log.Printf("stopping pipeline: %s", name)
				p.cancel()
				<-p.done
				delete(running, name)
			}
			return
		case <-ticker.C:
			poll()
		}
	}
}

type fakeFS struct{}

func (f *fakeFS) Open(name string) (fs.File, error) {
	return nil, errors.New("fake filesystem")
}

func (f *fakeFS) OpenFile(name string, flag int, perm fs.FileMode) (fs.File, error) {
	return nil, errors.New("fake filesystem")
}

func (f *fakeFS) Stat(name string) (fs.FileInfo, error) {
	return nil, errors.New("fake filesystem")
}

func (f *fakeFS) Remove(name string) error {
	return errors.New("fake filesystem")
}

func (f *fakeFS) MkdirAll(path string) {}
