# MEDODS Go Test task

**[Пример конфигурационного файла config.env](config.env)**

## Установка

### Docker
```bash
git clone https://github.com/AV3RAGE-ENJOYER/medods_test_task
cd medods_test_task
docker compose up -d
```

### Source

Для корректной работы приложения нужен работающий сервер **PostgreSQL** на **localhost**. Также нужно изменить хост в переменной **POSTGRES_URL** в файле [config.env](config.env)  

```bash
git clone https://github.com/AV3RAGE-ENJOYER/medods_test_task
cd medods_test_task
go mod download && go mod verify
go build main.go
./main
```

## Тестирование

```bash
go test -v ./tests
```

## Контакты

Домбровский Андрей
**Телефон: +79663838555**
**[Telegram](https://t.me/dombrovskii_andrei)**
**[Резюме на HeadHunter](https://hh.ru/resume/0f7bc270ff0d56019a0039ed1f737977446667?hhtmFrom=resume_list)**
