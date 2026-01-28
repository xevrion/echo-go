APP=echo-go

a:
	PORT=8080 go run ./cmd/$(APP)

b:
	PORT=8081 go run ./cmd/$(APP) $(ADDR)

two:
	@echo "Starting peer A (8080)"
	@PORT=8080 go run ./cmd/$(APP) & \
	sleep 1; \
	echo "Starting peer B (8081)"; \
	PORT=8081 go run ./cmd/$(APP)
