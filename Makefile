code = k8s-token-review-mock
ldflags = '-s -w -linkmode external -extldflags "-static"'

build:
	go build -ldflags ${ldflags} -o ${code} ${code}.go

# build in alpien with musl
build-musl: build-builder
	docker run -ti --rm -v "$$PWD:/app" -w /app --env CALLER_UID=$(shell id -u) --env outputFilename="${code}" alpine_go_builder make build docker-perms

build-builder:
	docker build -t alpine_go_builder .

docker-perms:
	[ -z "$${CALLER_UID}" ] || chown $${CALLER_UID}:$${CALLER_UID} $${outputFilename}
	chmod 750 $${outputFilename}

clean:
	docker rmi alpine_go_builder
	rm -f ${code}

