#### Поднять проект:
```
docker compose up --build
```

#### Регистрация:
```
POST http://0.0.0.0:80/user/register
{
    "username": "new_user",
    "password": "53047"
}
```

#### Вход:
```
POST http://0.0.0.0:80/user/login
{
    "username": "new_user",
    "password": "53047"
} + JWT
```

#### Создать дом:
```
POST http://0.0.0.0:80/houses
{
    "address": "more test house",
    "square": 23.0
} + JWT
```

#### Получить дома:
```
GET http://0.0.0.0:80/user/houses + JWT
```

#### Получить все предоставляемые модули компанией (слеш!):
```
GET http://0.0.0.0:80/device/modules/
```

#### Подключить модуль к дому:
```
POST http://0.0.0.0:80/device/modules/houses/5d19d994-12ef-40fc-9569-67bcbc800cfe/modules/15584fb6-d251-43a1-98f7-96c8497b6b43/assign
```

#### Убедиться в подключении дома к модулю:
```
GET http://0.0.0.0:80/device/modules/houses/5d19d994-12ef-40fc-9569-67bcbc800cfe
```

#### Выключить модуль:
```
POST http://0.0.0.0:80/device/modules/houses/5d19d994-12ef-40fc-9569-67bcbc800cfe/modules/15584fb6-d251-43a1-98f7-96c8497b6b43/turn-off
```

#### Включить модуль:
```
POST http://0.0.0.0:80/device/modules/houses/5d19d994-12ef-40fc-9569-67bcbc800cfe/modules/15584fb6-d251-43a1-98f7-96c8497b6b43/turn-on
```

#### Получить текущее состояние подключенного модуля к дому:
```
http://0.0.0.0:80/device/modules/houses/5d19d994-12ef-40fc-9569-67bcbc800cfe/modules/8176acb6-b8ca-44a3-8038-3f3b845dc1b6/state
```
