```markdown
# Order Service System
---

## Запуск системы

1. Перейдите в корневую папку проекта.
2. Запустите систему командой:  
   ```
   sudo docker compose --env-file ./order-service/.env.wb.example --env-file ./producer-simulator/.env.producer.example up
   ```

---

## Использование API

После запуска системы API будет доступно по адресу `http://localhost:8080`.

### Получение списка заказов

Для получения списка заказов откройте:  
`http://localhost:8080/api/orders/`

На странице будут ссылки на отдельные заказы.

---

### Просмотр информации о заказе

Если известен ID заказа, вы можете перейти по следующему адресу:  
`http://localhost:8080/api/orders/<example_id>`  
(замените `<example_id>` на реальный ID заказа).

---

### Параметры пагинации

Вы можете настроить пагинацию, добавив параметры `page` и `limit` в запрос:  
Пример:  
`http://localhost:8080/api/orders/?page=0&limit=10`

- **page** — номер страницы (начинается с 0).  
- **limit** — количество элементов на странице.

---

## Настройка конфигурации

### Конфигурация Order Service

Чтобы изменить параметры `order-service`, отредактируйте файл `./order-service/.env.wb.example`.

---

### Конфигурация Producer Simulator

Для изменения скорости записи в Kafka или параметров подключения к брокерам, измените файл `./producer-simulator/.env.producer.example`. Пример параметров:

```env
PRODUCER_SERVERS=kafka:9092
PRODUCER_PROTOCOL=PLAINTEXT
PRODUCER_ACKS=all

SCHEDULAR_PERIOD=10s
```

- **PRODUCER_SERVERS** — адрес Kafka-брокеров.  
- **PRODUCER_PROTOCOL** — протокол подключения (например, `PLAINTEXT`).  
- **PRODUCER_ACKS** — режим подтверждений Kafka.  
- **SCHEDULAR_PERIOD** — интервал записи данных продьюсером.

---  

## Контакты для обратной связи

Если возникли вопросы или проблемы, обратитесь к разработчику.  
email: ubelwertyas@gmail.com