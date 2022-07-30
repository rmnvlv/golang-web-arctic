# Web-Arctic

## Как запустить локально

Установить переменные окружения

```shel
set -o allexport; source .env; set +o allexport 
```

Пример **.env** файла
```shel
HCAPTCHA_SECRET_KEY="your_hcaptcha_secret_key"
HCAPTCHA_SITE_KEY="your_hcaptcha_site_key"

YANDEX_OAUTH_TOKEN="your_yandex_oauth_token"

SMTP_USER="your_smtp_user@your_domain.com"
SMTP_PASSWORD="your_smtp_user_password"
SMTP_HOST="smtp.your_domain.com"
SMTP_PORT=25

ADMIN_PASSWORD="123456"

DATABASE_URL="test.db"

PORT=8080
````

Запустить сервер
```shell
go run .
```

Когда меняются стили, tailwind должен знать об этом. 
Для этого нужно, чтобы tailwind следил за всеми изменения в html/css/js файлах и генерировал обновленный css файл.
```shell
npx tailwindcss -i ./styles/tailwind.css -o ./assets/css/tailwind.css --watch
```