# Status

- Сервис читает транзакции из Kafka через Consumer Group, обрабатывает их и пишет результат в PostgreSQL.
- Топик: `transactions`, 3 партиции, ключ партиционирования — `merchant_id`.
- На каждую партицию (`ConsumeClaim`) поднимается отдельный `partitionWorker` с очередью сообщений (`internal/adapters/kafka`), что реализует Worker Pool — по одному воркеру на партицию, сообщения одного мерчанта обрабатываются строго последовательно.
- Логика обработки транзакции (`internal/application/processor.go`):
  - мерчант ищется в LRU-кэше (`internal/cache/lru.go`, собственная реализация на `container/list` + `map`, потокобезопасная через `sync.RWMutex`);
  - при промахе кэша мерчант загружается из PostgreSQL и кладется в кэш;
  - если мерчант неактивен (`is_active = false`) — транзакция сохраняется со статусом `rejected`;
  - если активен — баланс мерчанта и транзакция обновляются атомарно в одной SQL-транзакции (`TxManager.WithinTransaction`), статус — `approved`, кэш обновляется новым балансом.
- **At-Least-Once**: offset сообщения коммитится (`session.MarkMessage`) только после успешной обработки транзакции и записи в PostgreSQL; при ошибке offset не фиксируется, сообщение будет обработано повторно. Повторная вставка транзакции защищена от дублей через `ON CONFLICT (id) DO NOTHING`.
- **Graceful shutdown** реализован полноценно: по `SIGINT/SIGTERM` отменяется контекст консьюмера, воркеры каждой партиции дорабатывают уже поставленные в очередь сообщения и закрываются (`partitionWorker.Close`), после чего ждём завершения через `sync.WaitGroup` с ограничением по времени `APP_SHUTDOWN_TIMEOUT` (если не уложились — принудительное завершение); закрываются соединения с Kafka consumer group и пулом PostgreSQL (`pgxpool`).
- Для работы с Kafka используется пакет `github.com/IBM/sarama`, для PostgreSQL — `github.com/jackc/pgx/v5` (`pgxpool`).
- Для работы подготовлен `docker/docker-compose.yml`, который запускает Apache Kafka 4.3.1, создает топик `transactions` с 3 партициями, если такого еще нет, и поднимает PostgreSQL 17 с инициализацией таблиц `merchants` и `transactions`.
- Инициализация схемы (`docker/postgres/initdb/01-schema.sql`) заранее заполняет таблицу `merchants` пятью тестовыми записями (баланс генерируется случайно при первом запуске, если таблица пуста):
  - `1` — Sunrise Bakery, активен;
  - `2` — Blue Horizon Electronics, активен;
  - `3` — Green Leaf Grocers, активен;
  - `4` — Old Town Bookstore, заблокирован (`is_active = false`);
  - `5` — Silver Line Motors, заблокирован (`is_active = false`).

- При запуске можно задать набор переменных окружения. Если использовать `APP_ENV = debug`, то автоматически применяются конфиги для локального запуска из Docker, игнорируя остальные параметры запуска.

- Переменные окружения:
  - `APP_ENV` — обязательная переменная; если равна `debug`, приложение использует конфиг по умолчанию для локального окружения, иначе загружает значения из остальных переменных.
  - `APP_SHUTDOWN_TIMEOUT` — целое число секунд; максимальное время ожидания завершения воркеров при graceful shutdown.
  - `CACHE_CAPACITY` — целое число больше 0; задает емкость LRU-кэша мерчантов.
  - `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_SSLMODE` — параметры подключения к PostgreSQL. `POSTGRES_PORT` должен быть положительным числом; `POSTGRES_SSLMODE` обязателен, например `disable`.
  - `KAFKA_BROKERS` — список брокеров Kafka через запятую, например `localhost:9092,localhost:9093`.
  - `KAFKA_TOPIC` — имя топика Kafka.
  - `KAFKA_GROUP_ID` — group ID consumer'а.

- Для режима `APP_ENV=debug` используются следующие значения по умолчанию:
  - `APP_SHUTDOWN_TIMEOUT=5s`
  - `CACHE_CAPACITY=1000`
  - PostgreSQL: `localhost:5432`, база `transactions`, пользователь `postgres`, пароль `password`, SSL `disable`
  - Kafka: брокер `localhost:9092`, топик `transactions`, group ID `payment-consumer-group`
  - Размер очереди воркера (`WorkerConfig.QueueSize`) — `1024`; чтение этого параметра из переменных окружения пока не реализовано (`TODO` в `loadWorkerConfig`), всегда используется значение по умолчанию.

- Известные ограничения:
  - Идентификаторы `Merchant.Id` и `Transaction.ID`/`MerchantID` — `int64`, приходят в JSON-сообщении Kafka как числа.
  - Тесты (`test/`) пока отсутствуют.

