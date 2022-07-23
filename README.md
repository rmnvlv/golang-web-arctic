# Web-Arctic

### Запуск локально

Установить переменные окружения
```shel
HCAPTCHA_SECRET_KEY="your_hcaptcha_secret_key"
HCAPTCHA_SITE_KEY="your_hcaptcha_site_key"
ADMIN_PASSWORD="123456"
````

```shell
go run .
```

Чтобы увидеть измения стилей, сначала нужно сгенерировать добавленные стили
Комада ниже будет отслеживать все изменения в html файлах и генерирывать новый css файл
```shell
npx tailwindcss -i ./styles/tailwind.css -o ./assets/css/tailwind.css --watch
```