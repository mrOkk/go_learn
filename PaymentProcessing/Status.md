# Status

- Сервис читает из Kafka и записывает в лог.
- Graceful shutdown реализован костылем в виде ожидания в течение 2 секунд; на текущий момент нечего закрывать, корректное решение пока не требуется.
- Топик: `transactions`.
- Для работы с Kafka используется пакет `github.com/IBM/sarama`.
- Для работы подготовлен `docker/docker-compose.yml`, который запускает Apache Kafka 4.3.1, создает топик с 3 партициями, если такого еще нет, и поднимает PostgreSQL с инициализацией таблиц `merchants` и `transactions`.

- При запуске можно задать набор переменных окружения. Если использовать `APP_ENV = debug`, то автоматически применяются конфиги для локального запуска из Docker, игнорируя остальные параметры запуска.

- Переменные окружения:
  - `APP_ENV` — обязательная переменная; если равна `debug`, приложение использует конфиг по умолчанию для локального окружения, иначе загружает значения из остальных переменных.
  - `CACHE_CAPACITY` — целое число больше 0; задает емкость кэша.
  - `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_SSLMODE` — параметры подключения к PostgreSQL. `POSTGRES_PORT` должен быть положительным числом; `POSTGRES_SSLMODE` обязателен, например `disable`.
  - `KAFKA_BROKERS` — список брокеров Kafka через запятую, например `localhost:9092,localhost:9093`.
  - `KAFKA_TOPIC` — имя топика Kafka.
  - `KAFKA_GROUP_ID` — group ID consumer'а.

- Для режима `APP_ENV=debug` используются следующие значения по умолчанию:
  - `CACHE_CAPACITY=1000`
  - PostgreSQL: `localhost:5432`, база `transactions`, пользователь `postgres`, пароль `password`, SSL `disable`
  - Kafka: брокер `localhost:9092`, топик `transactions`, group ID `payment-consumer-group`

