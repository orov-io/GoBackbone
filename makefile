default:
	@echo "=============Building Production Service============="
	docker build -f ./Dockerfile.prod -t demium/admin_service .

up: 
	@echo "=============Starting Service Locally============="
	docker-compose up -d

build:
	@echo "=============Building Development image============="
	docker-compose build

logs:
	docker-compose logs -f

down:
	@echo "=============Stopping Development Service============="
	docker-compose down

test:
	@echo "=============Stopping Development Service============="
	goconvey

deatachedTest:
	@echo "=============Launching test suite in deatached mode============="
	goconvey >/dev/null &

init:
	dep init -v

update:
	dep ensure --update -v

reload: down up

restart: down build up

clean: down
	@echo "=============cleaning up============="
	rm -f api
	docker system prune -f
	docker volume prune -f
	killall goconvey