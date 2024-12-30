# Smart Home Monolith

## Описание

Проект "Smart Home Monolith" представляет собой монолитное приложение для управления отоплением и мониторинга температуры в умном доме. Пользователи
могут удаленно включать/выключать отопление, устанавливать желаемую температуру и просматривать текущую температуру через веб-интерфейс.

## Структура проекта

Проект организован в соответствии со стандартной структурой Maven/Gradle проекта:

- **Основной класс приложения** аннотирован `@SpringBootApplication`.
- **Пакеты**:
    - `controller` - контроллеры для обработки HTTP-запросов.
    - `service` - сервисы для бизнес-логики.
    - `repository` - репозитории для доступа к базе данных.
    - `entity` - сущности JPA.
    - `dto` - объекты передачи данных (DTO).

## Зависимости

Проект использует следующие зависимости:

- Java 17
- Maven >=3.8.1 && <4.0.0
- `spring-boot-starter-web` - для создания веб-приложений.
- `spring-boot-starter-data-jpa` - для работы с JPA и базой данных.
- `postgresql` - драйвер для работы с PostgreSQL.
- `springfox-boot-starter` - для интеграции Swagger.
- `lombok` - для упрощения написания кода.
- `spring-boot-starter-test`, `spring-boot-testcontainers`, `junit-jupiter`, `testcontainers` - для тестирования.

## Запуск

1.

``` shell
docker run --name postgres-db -e POSTGRES_DB=smart_home -e POSTGRES_USER=your_username -e POSTGRES_PASSWORD=your_password -p 5432:5432 -d postgres:13-alpine
```

2.

``` shell
mvn spring-boot:run
```

## Тесты

``` shell
mvn test
```

## Конфигурация

### Подключение к базе данных

Конфигурация подключения к PostgreSQL находится в файле `application.yml`:

```yaml
spring:
  datasource:
    url: jdbc:postgresql://localhost:5432/smart_home
    username: user
    password: password
  jpa:
    hibernate:
      ddl-auto: update
    show-sql: true
```

## API
* GET /api/heating/{id} - Fetching heating system with id {id} - Получение информации о системе отопления
* PUT /api/heating/{id} - Updating heating system with id {id} - Обновление системы отопления
* POST /api/heating/{id}/turn-on - Turning on heating system with id {id} - Включение системы отопления
* POST /api/heating/{id}/turn-off - Turning off heating system with id {id} - Выключение системы отопления
* POST /api/heating/{id}/set-temperature - Setting target temperature to {temperature} for heating system with id {id} - Установка целевой температуры
* GET /api/heating/{id}/current-temperature - Fetching current temperature for heating system with id {id} - Получение текущей температуры


## Domains
* управление системой отопления
* взаимодействие с датчиками температуры

## Логирование

В проекте используется SLF4J для логирования.

## Тестирование

В проекте написаны модульные тесты для сервисов и контроллеров.
