# بخش ۰۳ — build و اجرای لوکال پروژه

## تاریخ

```text
2026-05-05
```

## هدف

در این بخش بررسی شد که پروژه `loki-ui` از روی سورس قابل build و اجرا باشد و تنظیمات runtime از طریق environment انجام شود.

هدف این بود که قبل از GitHub، server deploy یا systemd، ابتدا baseline لوکال پروژه واقعاً قابل اجرا باشد.

## کارهای انجام‌شده

در این بخش این کارها انجام شد:

```text
نصب toolchain سالم Go داخل WSL
تشخیص خراب بودن یا ناقص بودن tarballهای قبلی Go
دریافت tarball سالم Go از منبع رسمی
بررسی gzip archive
بررسی SHA256 checksum
نصب Go در /usr/local/go
اضافه کردن Go به PATH
build گرفتن از سورس loki-ui
اجرای binary با LOKI_URL
بررسی اتصال runtime به Loki
```

## مشکل محیطی

مشکل اصلی این بخش مربوط به سورس پروژه نبود. مشکل در toolchain محلی Go بود.

برخی فایل‌های قبلی Go ناقص یا خراب بودند و نمی‌شد به build آن‌ها اعتماد کرد. به همین دلیل toolchain سالم دوباره نصب و validate شد.

## تنظیمات runtime

برنامه برای اجرا به `LOKI_URL` نیاز داشت.

نمونه مدل runtime:

```text
LOKI_URL=http://127.0.0.1:3100
```

این مقدار داخل `.env` واقعی یا environment قرار می‌گیرد و نباید به عنوان secret یا config واقعی در Git منتشر شود.

## نتیجه build

بعد از نصب Go سالم، پروژه از سورس build شد و binary ساخته شد.

این نشان داد که ساختار packageهای Go، entrypoint و templateهای موجود از نظر build baseline سالم هستند.

## نتیجه runtime

برنامه با environment لازم اجرا شد و مدل خواندن از Loki در حالت لوکال بررسی شد.

در این مرحله هدف فقط اثبات قابل اجرا بودن برنامه بود، نه deploy production.

## نکته امنیتی

اجرای لوکال و اتصال به Loki نباید به معنی public بودن UI تلقی شود. مدل دسترسی همچنان باید loopback-only و internal بماند.

## ارزش فنی

این بخش نشان داد که پروژه فقط در حد فایل‌های repository نیست و واقعاً از روی source قابل build است.

همچنین مشکل toolchain از مشکل application جدا شد؛ سورس برنامه بی‌دلیل مقصر شناخته نشد.

## ارزش رزومه‌ای قابل دفاع

```text
Validated local Go build and runtime execution for an internal Loki UI, including environment-based Loki configuration and verification of a clean Linux Go toolchain inside WSL.
```

## محدودیت این بخش

این بخش شامل GitHub remote، server deploy، systemd service، hardening یا SSH tunnel نبود. این موارد در بخش‌های بعدی انجام شدند.
