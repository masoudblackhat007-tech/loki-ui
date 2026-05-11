# بخش ۱۴ — اضافه کردن endpoint سلامت healthz

## تاریخ

```text
2026-05-07
```

## هدف

در این بخش یک endpoint سبک برای health check به `loki-ui` اضافه شد.

قبل از این بخش، بررسی سلامت سرویس با صفحه واقعی UI مثل `/logs` انجام می‌شد.

## مشکل

مسیر `/logs` یک صفحه واقعی است و ممکن است به template rendering یا منطق application وابسته باشد.

برای smoke test سرویس، بهتر است endpoint ساده‌تر و مستقل‌تری وجود داشته باشد.

## endpoint اضافه‌شده

```text
/healthz
```

## رفتار endpoint

این endpoint فقط نشان می‌دهد process HTTP برنامه زنده است.

برای methodهای `GET` و `HEAD` پاسخ موفق می‌دهد و برای methodهای دیگر `405 Method Not Allowed` برمی‌گرداند.

## نتیجه

از این بخش به بعد می‌توان زنده بودن process را بدون نیاز به render صفحه اصلی بررسی کرد.

## تفاوت با readiness

`/healthz` فقط خود process را بررسی می‌کند.

بررسی اتصال به Loki در بخش بعدی با `/readyz` اضافه شد.

## نکته امنیتی

این endpoint اطلاعات حساس برنمی‌گرداند. فقط پاسخ ساده health می‌دهد.

## ارزش فنی

این بخش operational check تمیزتری برای سرویس فراهم کرد.

## ارزش رزومه‌ای قابل دفاع

```text
Added a lightweight /healthz endpoint to a Go-based internal Loki UI to support process-level smoke checks without depending on full UI rendering or Loki connectivity.
```

## محدودیت این بخش

این بخش readiness، auth، TLS، rate limiting یا public exposure اضافه نکرد.
