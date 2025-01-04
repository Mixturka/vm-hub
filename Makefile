START_TEST_ENV_SCRIPT := ./scripts/test/test_env_up.sh
STOP_TEST_ENV_SCRIPT := ./scripts/test/test_env_down.sh
MIGRATE_TEST_SCRIPT := ./scripts/test/migrate_test.sh
RESET_TEST_DB_SCRIPT := ./scripts/test/reset_test_db.sh
RUN_TESTS_SCRIPT := ./scripts/test/test.sh
LOAD_ENV := ./scripts/load_env.sh

.PHONY: test-env-up test-env-down migrate-test test clean-tests reset-db

test-env-up:
	@$(START_TEST_ENV_SCRIPT)

test-env-down:
	@$(STOP_TEST_ENV_SCRIPT)

migrate-test:
	@$(MIGRATE_TEST_SCRIPT)

reset-db:
	@$(RESET_TEST_DB_SCRIPT)

test:
	@. $(LOAD_ENV) && $(RUN_TESTS_SCRIPT)

clean-tests:
	@$(MAKE) reset-db
	@$(MAKE) test-env-down
