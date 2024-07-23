DATASE_ARGS := "mysql://user:password@tcp(localhost:3306)/todo"

# マイグレーションUp実行
.PHONY: migrate-up
migrate-up:
	migrate -database $(DATASE_ARGS) -path db/migrations up

# マイグレーションDown実行
.PHONY: migrate-down
migrate-down:
	migrate -database $(DATASE_ARGS) -path db/migrations down -all

# マイグレーションファイルのクリア
.PHONY: migrate-clear
migrate-clear:
	migrate -database $(DATASE_ARGS) -path db/migrations force 1

# マイグレーションバージョン確認
.PHONY: migrate-version
migrate-version:
	migrate -database $(DATASE_ARGS) -path db/migrations version

# マイグレーションを実行せずに現在のバージョンをセットする
.PHONY: migrate-force
migrate-force:
	migrate -database $(DATASE_ARGS) -path db/migrations force $(VERSION)

# マイグレーションファイルの作成(tableName = テーブル名(複数形)として引数に渡す)
.PHONY: create-migrate-file
create-migrate-file:
	migrate create -ext sql -dir db/migrations -seq ${tableName}

# モックの生成
.PHONY: gen-mock
gen-mock:
	@rm -rf mock/$(MOCK_DIR)/*.go
	@for r in $$(find $(MOCK_DIR) -type f -name '*.go' ! -name '*_test.go' -exec basename {} \;); do \
		echo "mockgen -source $(MOCK_DIR)/$$r -destination mock/$(MOCK_DIR)/$$r"; \
		mockgen -source $(MOCK_DIR)/$$r -destination mock/$(MOCK_DIR)/$$r; \
	done

# usecaseのモックの生成
.PHONY: gen-mock-usecase
gen-mock-usecase:
	@make gen-mock MOCK_DIR=usecase

# usecaseのモックの生成
.PHONY: gen-mock-repository
gen-mock-repository:
	@make gen-mock MOCK_DIR=repository

# テスト
.PHONY: test
test:
	go test -race ./...
