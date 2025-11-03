DB_PROPS=migrations/liquibase.properties
MIGRATIONS_DIR=migrations/changes
MASTER_FILE=migrations/db.changelog-master.yaml

migrate:
	liquibase --defaultsFile=$(DB_PROPS) update

status:
	liquibase --defaultsFile=$(DB_PROPS) status --verbose

rollback:
	liquibase --defaultsFile=$(DB_PROPS) rollbackCount 1

update-sql:
	liquibase --defaultsFile=$(DB_PROPS) updateSQL

validate:
	liquibase --defaultsFile=$(DB_PROPS) validate

history:
	liquibase --defaultsFile=$(DB_PROPS) history

new-migration:
	@timestamp=$$(date +%Y%m%d%H%M%S); \
	forward="$$timestamp-$(name).sql"; \
	rollback="$$timestamp-$(name)-rollback.sql"; \
	touch $(MIGRATIONS_DIR)/$$forward; \
	touch $(MIGRATIONS_DIR)/$$rollback; \
	echo "Created: $(MIGRATIONS_DIR)/$$forward and $(MIGRATIONS_DIR)/$$rollback"; \
	echo "" >> $(MASTER_FILE); \
	echo "  - changeSet:" >> $(MASTER_FILE); \
	echo "      id: $$timestamp-$(name)" >> $(MASTER_FILE); \
	echo "      author: $(USER)" >> $(MASTER_FILE); \
	echo "      changes:" >> $(MASTER_FILE); \
	echo "        - sqlFile:" >> $(MASTER_FILE); \
	echo "            path: migrations/changes/$$forward" >> $(MASTER_FILE); \
	echo "      rollback:" >> $(MASTER_FILE); \
	echo "        - sqlFile:" >> $(MASTER_FILE); \
	echo "            path: migrations/changes/$$rollback" >> $(MASTER_FILE); \
	echo "  -> Added to $(MASTER_FILE)"