# Optimizations

 - [X] Audio block/batch processing
 - [X] Lock free
   - [X] SPSC queues
   - [X] Noise local rng (no lock on global rng src)
 - [X] Zero alloc hotpath
   - [X] Polyphonic voice management, mixer and oscillator hotpaths are tested with no allocations
   - [ ] Add more tests
 - [X] LUT caching
 - [/] Wave tables (sine wave only)
   - [ ] Add more waveforms
 - [X] SIMD optimizations in mixer
 - [/] Write benchmarks and profile (ocs partially done)
 - [X] Pan from equal power to linear



