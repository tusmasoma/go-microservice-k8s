SERVICE_PATH_PREFIX := .
SERVICES := catalog customer order

DB_USER   := root
DB_PASS   := root
CONTAINER := mysql_test

.PHONY: all start migrate dump stop clean

# 全体実行: 初期化 → コンテナ起動 → 差分検知を元にマイグレーション実行 → 各サービスの最終DDLダンプ
gen: clean start migrate dump

# 1. DockerでMySQLコンテナを起動し、各サービス用のデータベースを作成
start:
	@echo "Starting MySQL Docker container..."
	docker run --name $(CONTAINER) -e MYSQL_ROOT_PASSWORD=$(DB_PASS) -p 3306:3306 -d mysql:latest
	@echo "Waiting for MySQL to be ready..."
	@sleep 20
	@for service in $(SERVICES); do \
	    echo "Creating database $$service..."; \
	    docker exec -i $(CONTAINER) mysql -u$(DB_USER) -p$(DB_PASS) -e "CREATE DATABASE IF NOT EXISTS \`$$service\`;"; \
	done

# 2. Gitの差分検知で、変更のあったサービスのみマイグレーション実行
migrate:
ifdef SERVICE
	@echo "Migrating for service: $(SERVICE)"; \
	$(MIGRATE_FUNC); \
	_migrate $(SERVICE)
else
	@for service in $(SERVICES); do \
	    echo "Migrating for service: $$service"; \
	    $(MIGRATE_FUNC); \
	    _migrate $$service; \
	done
endif

# 3. 各サービスごとに mysqldump でスキーマのみをダンプし、各ディレクトリの [service].ddl に出力
dump:
	@echo "Dumping schemas for each service..."
	@for service in $(SERVICES); do \
	    echo "Dumping schema for $$service to $(SERVICE_PATH_PREFIX)/$$service/$$service.ddl..."; \
	    docker exec $(CONTAINER) mysqldump -u$(DB_USER) -p$(DB_PASS) --no-data --compact --skip-comments --skip-set-charset $$service | \
	    sed '/^DROP TABLE/d' | sed '/^\/\*/d' | sed '/^--/d' | sed 's/ ENGINE=.*$$//' | sed 's/ DEFAULT CHARSET=.*$$//' | sed 's/`//g' > $(SERVICE_PATH_PREFIX)/$$service/$$service.ddl; \
	done
	@echo "Schema dump completed."

# コンテナの停止と削除
stop:
	@echo "Stopping and removing Docker container..."
	docker stop $(CONTAINER)
	docker rm -v $(CONTAINER)

# クリーンアップ：コンテナ停止＋生成ファイルの削除
clean: stop
	@echo "Cleaning up generated dump files..."
	@for service in $(SERVICES); do \
	    rm -f $(SERVICE_PATH_PREFIX)/$$service/$$service.ddl; \
	done

define MIGRATE_FUNC
_migrate() { \
  echo "Running migration for service: $$1"; \
  if [ -n "$$(git diff --name-only origin/main -- $(SERVICE_PATH_PREFIX)/$$1)" ]; then \
      echo "Differences detected in $$1"; \
      for file in $$(ls $(SERVICE_PATH_PREFIX)/$$1/*.ddl | sort); do \
          echo "Applying $$file to database $$1..."; \
          docker exec -i $(CONTAINER) mysql -u$(DB_USER) -p$(DB_PASS) -D$$1 < $$file; \
      done; \
  else \
      echo "No migration differences in $$1"; \
  fi; \
}
endef
export MIGRATE_FUNC