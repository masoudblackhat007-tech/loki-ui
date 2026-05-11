# بخش ۱۶ — اضافه کردن request logging امن برای loki-ui

## تاریخ

```text
2026-05-07
```

## هدف

در این بخش logging ساختاریافته و امن برای requestهای خود `loki-ui` اضافه شد.

هدف این بود که رفتار خود سرویس قابل مشاهده باشد، بدون اینکه داده‌های حساس request ثبت شوند.

## مشکل

قبل از این بخش، سرویس فقط startup و shutdown را log می‌کرد.

برای بررسی رفتار عملیاتی، لازم بود requestهای ورودی هم ثبت شوند.

اما request logging اگر خام انجام شود خطرناک است، چون ممکن است headerها، cookie، authorization، query string، payload یا response body را log کند.

## تصمیم امنیتی

در این بخش فقط metadata کم‌ریسک ثبت شد:

```text
request_id
method
path
status
duration_ms
```

این موارد ثبت نشدند:

```text
headers
cookies
authorization values
raw query strings
request body
response body
```

## تغییر انجام‌شده

یک middleware برای request logging اضافه شد.

برای هر request یک `X-Request-Id` تولید شد و در response header قرار گرفت.

status code و duration اندازه‌گیری و در قالب JSON log نوشته شد.

## نتیجه

از این بخش به بعد خود `loki-ui` هم requestهایش را به شکل structured و محدود log می‌کند.

## نکته امنیتی

`request_id` secret نیست، اما مقدار correlation است و نباید بی‌دلیل در مستندات عمومی یا screenshotهای رزومه‌ای منتشر شود.

## ارزش فنی

این بخش observability خود ابزار observability را بهتر کرد، بدون اینکه raw sensitive data وارد لاگ شود.

## ارزش رزومه‌ای قابل دفاع

```text
Added safe structured request logging to an internal Go-based Loki UI, capturing request id, method, path, status, and duration while intentionally excluding headers, cookies, authorization data, request bodies, and response bodies.
```

## محدودیت این بخش

این بخش audit logging کامل، authentication، authorization، rate limiting یا تغییر در Laravel/Loki/Alloy اضافه نکرد.
