# BPMN to Hyperledger Fabric Translator

Инструмент для автоматической генерации смарт-контрактов Hyperledger Fabric
из BPMN-хореографий. Включает веб-редактор диаграмм и сервис генерации кода.

## Запуск

```bash
docker-compose up --build
```

- Редактор: http://localhost:9013  
- API: http://localhost:8000


## Структура проекта

| Путь | Описание |
|------|----------|
| `backend/src/` | Транслятор|
| `backend/generated_contracts/` | Результаты генерации |
| `frontend/` | Веб-редактор на базе chor-js |
| `tests/examples/` | Сгенерированные контракты для 4 контрольных моделей |
| `tests/start_testnet.sh` | Скрипт запуска тестовой сети Hyperledger Fabric |
| `tests/invoke_cc.sh` | Скрипт вызова транзакций |
