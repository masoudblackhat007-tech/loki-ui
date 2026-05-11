# بخش ۰۵ — clone روی سرور، build و تست اتصال به Loki

## تاریخ

```text
2026-05-06
```

## هدف

در این بخش repository پروژه `loki-ui` روی سرور clone شد، از روی سورس build شد، به Loki محلی وصل شد و توانست لاگ‌های Laravel که توسط Alloy وارد Loki شده بودند را نمایش دهد.

هدف این بود که پروژه از حالت لوکال خارج شود و روی سرور واقعی، در کنار Loki و Alloy، قابل اجرا باشد.

## محیط سرور

```text
Hostname: 381239
User: deploy
Project path: /home/deploy/apps/loki-ui
Repository: git@github.com:masoudblackhat007-tech/loki-ui.git
```

## کارهای انجام‌شده

در این بخش این کارها انجام شد:

```text
بررسی user deploy روی سرور
ساخت مسیر /home/deploy/apps
clone کردن repository از GitHub
بررسی clean بودن working tree
build گرفتن از سورس روی سرور
تنظیم LOKI_URL برای اتصال به Loki محلی
اجرای loki-ui روی 127.0.0.1:18090
ارسال request به UI
بررسی اینکه UI لاگ‌های Laravel را از Loki می‌خواند
```

## مدل اتصال

مدل اجرا روی سرور:

```text
loki-ui -> 127.0.0.1:3100 -> Loki -> Laravel logs
```

در این مرحله `loki-ui` مستقیماً به اینترنت وصل نشد و نباید public می‌شد.

## نتیجه build

برنامه روی سرور از سورس build شد. این یعنی سورس داخل GitHub برای محیط سرور هم قابل استفاده بود و وابسته به binary لوکال نبود.

## نتیجه integration

UI توانست از Loki لاگ‌های Laravel را بخواند. این یعنی زنجیره زیر در عمل کار می‌کرد:

```text
Laravel JSON logs -> Alloy -> Loki -> loki-ui
```

## نکته امنیتی

در این بخش هنوز سرویس systemd ساخته نشده بود و access model نهایی از طریق SSH tunnel در بخش بعدی کامل‌تر شد.

با این حال اصل مهم حفظ شد: برنامه باید روی loopback بماند و پورت `18090` نباید در firewall عمومی باز شود.

## ارزش فنی

این بخش اولین اثبات server-side integration بود. پروژه فقط build لوکال نداشت؛ روی سرور واقعی هم build شد و به dependency اصلی خود یعنی Loki وصل شد.

## ارزش رزومه‌ای قابل دفاع

```text
Cloned, built, and verified an internal Go-based Loki UI on a Linux server, confirming end-to-end visibility from Laravel JSON logs through Alloy and Loki into the custom UI while keeping the service loopback-only.
```

## محدودیت این بخش

این بخش هنوز systemd service، hardening کامل، graceful shutdown یا SSH tunnel workflow مستندشده را اضافه نکرد. این موارد در بخش‌های بعدی انجام شدند.
