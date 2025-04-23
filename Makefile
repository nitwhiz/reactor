.PHONY: clean
clean:
	rm out/reactor

build: out/reactor

out/reactor:
	CGO_ENABLED=false go build -o out/reactor ./cmd/main

.PHONY: run
run: build
	out/reactor

.PHONY: pprof_cpu
pprof_cpu: build
	docker compose up --build pprof_cpu

.PHONY: pprof_mem
pprof_mem: build
	docker compose up --build pprof_mem
