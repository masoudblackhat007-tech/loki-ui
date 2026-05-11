# بخش ۰۱ — راه‌اندازی Git برای پروژه loki-ui

## تاریخ

```text
2026-05-05
```

## هدف

در این بخش، مخزن Git برای پروژه Go با نام `loki-ui` در مسیر درست راه‌اندازی شد و یک baseline تمیز و source-only ساخته شد.

هدف این بود که از همان شروع پروژه، فقط فایل‌های لازم و قابل انتشار وارد repository شوند و فایل‌های حساس، خروجی build، metadata محیط توسعه و فایل‌های runtime وارد Git نشوند.

## مسیر پروژه

```text
/mnt/d/project-learn/loki-ui
```

## وضعیت اولیه

در شروع کار، دو نسخه محلی از پروژه وجود داشت:

```text
/mnt/d/project-learn/loki-ui
/home/masoud/projects/loki-ui
```

مسیر زیر به عنوان source of truth انتخاب شد:

```text
/mnt/d/project-learn/loki-ui
```

این تصمیم مهم بود، چون ادامه توسعه، commit، push، build و deploy باید از یک مسیر مشخص و قابل تکرار انجام شود.

## کارهای انجام‌شده

در این بخش این کارها انجام شد:

```text
بررسی مسیر درست پروژه
راه‌اندازی Git
تنظیم branch اصلی روی main
ساخت فایل .gitignore
ساخت فایل .env.example
بررسی فایل‌های ignored
ساخت اولین commit تمیز
اصلاح اشتباه rule مربوط به .gitignore
amend کردن commit اولیه
```

## فایل‌های واردشده به Git

baseline اولیه شامل این فایل‌ها بود:

```text
.env.example
.gitignore
cmd/loki-ui/main.go
go.mod
internal/httpserver/handler.go
internal/httpserver/server.go
internal/loki/client.go
internal/loki/types.go
templates/layout.tmpl
templates/log_detail.tmpl
templates/logs.tmpl
```

## فایل‌های عمداً ignore شده

این موارد وارد Git نشدند:

```text
.env
.idea/
bin/
loki-ui
loki-ui.new
```

## نکته امنیتی مهم

فایل واقعی `.env` وارد repository نشد. این یعنی مقدارهای runtime، secretها، آدرس‌های داخلی و credentialها در Git ذخیره نشدند.

خروجی‌های build و فایل‌های IDE هم وارد Git نشدند تا repository قابل تکرار و source-only باقی بماند.

## خطای پیدا و اصلاح‌شده

در نسخه اولیه `.gitignore`، pattern زیر وجود داشت:

```gitignore
loki-ui
```

این rule خطرناک بود، چون ممکن بود مسیر زیر را هم ignore کند:

```text
cmd/loki-ui/main.go
```

این فایل entrypoint اصلی برنامه است و اگر وارد Git نمی‌شد، repository ناقص می‌شد.

نسخه اصلاح‌شده این بود:

```gitignore
/loki-ui
```

این rule فقط binary خروجی در root پروژه را ignore می‌کند و دیگر به سورس داخل `cmd/loki-ui` آسیب نمی‌زند.

## اعتبارسنجی

اعتبارسنجی نشان داد که `cmd/loki-ui/main.go` ignore نشده و فایل‌هایی مثل `.env`، `.idea/`، `bin/` و binaryهای خروجی درست ignore شده‌اند.

## commit نهایی بخش

```text
4ac811d Initial loki-ui project
```

## نتیجه

پروژه با یک baseline تمیز Git شروع شد. repository از همان ابتدا source-only بود و فایل‌های حساس یا generated وارد Git نشدند.

## ارزش رزومه‌ای قابل دفاع

```text
Initialized a clean source-only Git baseline for an internal Go-based Loki UI project, including secure ignore rules, environment example separation, and validation to prevent secrets, IDE metadata, build artifacts, and runtime files from entering source control.
```

## محدودیت این بخش

این بخش deploy، observability integration، systemd hardening، authentication یا authorization اضافه نکرد. فقط baseline امن و قابل دفاع برای source control ساخته شد.
