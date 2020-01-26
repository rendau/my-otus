build:
	+$(MAKE) -C api
	+$(MAKE) -C scheduler
	+$(MAKE) -C sender

clean:
	+$(MAKE) -C api clean
	+$(MAKE) -C scheduler clean
	+$(MAKE) -C sender clean

test:
	docker-compose -f docker-compose.test.yml up --build --exit-code-from calendar_tests
	@EXIT_CODE=$$?
	docker-compose -f docker-compose.test.yml down
	@exit $$EXIT_CODE
