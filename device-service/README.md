* Получение всех модулей: 
GET /modules

* Получение модулей, подключенных к дому:
GET /modules/:houseID

* Включение модуля:
PUT /modules/:houseID/:moduleID/turn-on

* Отключение модуля:
PUT /modules/:houseID/:moduleID/turn-off

* Добавление нового модуля:
POST /houses/:houseID/modules
