BENCH_MASK ?= .
BENCH_TARGET ?= /tmp/synth-bench.txt
BENCH_OLD_TARGET ?= /tmp/previous-synth-bench.txt
BENCH_COUNT ?= 10

bench:
	touch $(BENCH_TARGET)
	cp $(BENCH_TARGET) $(BENCH_OLD_TARGET)
	go test -benchmem -bench=$(BENCH_MASK) -count=$(BENCH_COUNT) ./... | tee $(BENCH_TARGET)
	benchstat $(BENCH_OLD_TARGET) $(BENCH_TARGET)

tests:
	go test ./...