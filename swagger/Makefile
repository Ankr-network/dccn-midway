.PHONY: restart stop start echo

restart: stop start

start:
	@cp swagger.yaml /tmp/
	@docker container run --name ankr-swagger \
		-d --rm -p 8888:8080 \
		-e SWAGGER_JSON=/foo/swagger.yaml \
		-v /tmp:/foo swaggerapi/swagger-ui > /dev/null
	@echo "Now you can browse \"http://localhost:8888/\"" \
		to view the api specification.

stop:
	@CID="$$(docker ps --all --quiet --filter=name='ankr-swagger')"; \
	if [ -n "$$CID" ]; then \
		docker container stop ankr-swagger > /dev/null; \
	fi
	@rm -rf /tmp/swagger.yaml
	@echo "Server is shutdown"