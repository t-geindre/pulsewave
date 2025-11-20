BENCH_MASK ?= .
BENCH_COUNT ?= 10
BENCH_DIR := ./bench
BENCH_TIMESTAMP := $(shell date +"%Y-%m-%d-%H_%M_%S")
BENCH_MASK_SAFE := $(shell echo "$(BENCH_MASK)" | tr '/.*^' '_' | cut -c1-40)
BENCH_TARGET := $(BENCH_DIR)/bench-$(BENCH_MASK_SAFE)-$(BENCH_TIMESTAMP).txt

.PHONY: bench
bench:
	mkdir -p $(BENCH_DIR)
	@echo "Running benchmarks (mask=$(BENCH_MASK)) → $(BENCH_TARGET)"
	go test -benchmem -run=^$$ -bench=$(BENCH_MASK) -count=$(BENCH_COUNT) ./... | tee $(BENCH_TARGET)

.PHONY: benchstat
benchstat:
	@latest=$$(ls -t $(BENCH_DIR)/bench-*.txt 2>/dev/null | head -n 1); \
	prev=$$(ls -t $(BENCH_DIR)/bench-*.txt 2>/dev/null | head -n 2 | tail -n 1); \
	if [ -z "$$latest" ] || [ -z "$$prev" ]; then \
		echo "Not enough benchmark files found in $(BENCH_DIR)"; \
		exit 1; \
	fi; \
	echo "Comparing:"; \
	echo "  OLD: $$prev"; \
	echo "  NEW: $$latest"; \
	benchstat $$prev $$latest

.PHONY: pprof-cpu
pprof-cpu:
	sleep 5 && go tool pprof -seconds 20 http://localhost:6060/debug/pprof/profile

.PHONY: pprof-alloc
.pprof-alloc:
	sleep 5 && go tool pprof -alloc_space -seconds 20 http://localhost:6060/debug/pprof/heap

.PHONY: tests
tests:
	go test ./...

.PHONY: proto
proto:
	protoc --go_out=. preset/preset.proto settings/settings.proto

.PHONY: dist-assets
dist-assets:
	mkdir -p dist/assets
	cp -r assets/fonts dist/assets/fonts
	cp -r assets/imgs dist/assets/imgs
	cp -r assets/presets dist/assets/presets
	cp assets/assets.json dist/assets/assets.json

.PHONY:
dist: dist-assets
	CGO_ENABLED=1 go build -trimpath -ldflags="-s -w" -o dist/pulsewave ./cmd/pulsewave

