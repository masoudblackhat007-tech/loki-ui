# بخش ۰۹ — اضافه کردن timeoutهای HTTP server

## تاریخ

```text
2026-05-06
```

## هدف

در این بخش timeoutهای explicit برای HTTP server پروژه `loki-ui` اضافه شد.

قبل از این بخش، server با رفتار پیش‌فرض Go و `http.ListenAndServe` اجرا می‌شد.

## مشکل

رفتار پیش‌فرض HTTP server برای ابزار عملیاتی مناسب نیست، حتی اگر سرویس فقط loopback باشد.

بدون timeoutهای مشخص، اتصال‌های کند یا نیمه‌باز می‌توانند resource را بیش از حد نگه دارند.

## تغییر انجام‌شده

به جای اجرای ساده با `http.ListenAndServe`، از `http.Server` با timeoutهای مشخص استفاده شد.

timeoutهای مهم شامل موارد زیر بودند:

```text
ReadHeaderTimeout
ReadTimeout
WriteTimeout
IdleTimeout
```

## نتیجه

رفتار HTTP server کنترل‌شده‌تر شد و runtime برنامه برای محیط عملیاتی سالم‌تر شد.

## نکته امنیتی

این hardening جایگزین authentication، authorization، TLS یا rate limiting نیست.

اما برای کاهش ریسک اتصال‌های کند و رفتارهای نامطلوب HTTP server لازم است.

## اعتبارسنجی

بعد از تغییر، پروژه build شد و سرویس همچنان روی loopback اجرا شد.

مسیرهای UI و API بدون تغییر در contract قبلی کار کردند.

## ارزش فنی

این بخش runtime HTTP برنامه را از حالت ساده توسعه‌ای به حالت قابل دفاع‌تر برای ابزار داخلی تبدیل کرد.

## ارزش رزومه‌ای قابل دفاع

```text
Added explicit Go HTTP server timeouts to an internal Loki UI service to avoid relying on default server behavior and improve operational robustness while preserving the loopback-only access model.
```

## محدودیت این بخش

این بخش public exposure، TLS، auth، rate limiting یا تغییر در Loki و Alloy اضافه نکرد.
