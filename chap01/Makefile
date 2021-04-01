-include ../Makefile

data_race: ## run data_race
	go run data_race.go

sync: ## run sync
	go run sync.go

deadlock: ## run deadlock
	go run deadlock.go

livelock: ## run livelock
	GOMAXPROCS=1 go run livelock.go
