# بخش ۰۶ — ساخت سرویس systemd برای loki-ui

## تاریخ

```text
2026-05-06
```

## هدف

در این بخش برای `loki-ui` یک سرویس systemd ساخته شد تا برنامه به جای اجرای دستی، به عنوان سرویس مدیریت‌شده روی سرور اجرا شود.

هدف این بود که start، stop، restart، status و auto-start برنامه از طریق systemd قابل کنترل باشد.

## محیط سرور

```text
Hostname: 381239
User: deploy
Project path: /home/deploy/apps/loki-ui
Service file: /etc/systemd/system/loki-ui.service
```

## کارهای انجام‌شده

در این بخش این کارها انجام شد:

```text
ساخت فایل سرویس systemd
تنظیم working directory پروژه
تنظیم اجرای binary ساخته‌شده
تنظیم environment برای LOKI_URL
فعال کردن سرویس
شروع سرویس
بررسی status سرویس
بررسی listen شدن روی 127.0.0.1:18090
```

## مدل سرویس

سرویس باید binary موجود در مسیر پروژه را اجرا کند و به Loki محلی وصل شود.

مدل runtime:

```text
systemd -> loki-ui binary -> Loki on 127.0.0.1:3100
```

## نتیجه

بعد از ساخت سرویس، `loki-ui` با systemd اجرا شد و وضعیت آن active بود.

این یعنی برنامه دیگر به session دستی shell وابسته نبود.

## نکته امنیتی

ساخت سرویس systemd به معنی public کردن UI نبود. برنامه همچنان باید روی `127.0.0.1:18090` اجرا شود.

پورت `18090` نباید در UFW باز شود.

## ارزش فنی

این بخش runtime برنامه را production-like کرد. مدیریت سرویس، restart و status از این مرحله قابل کنترل و قابل مستندسازی شدند.

## ارزش رزومه‌ای قابل دفاع

```text
Created and enabled a systemd-managed runtime for an internal Go-based Loki UI, allowing controlled service lifecycle management while preserving a loopback-only access model.
```

## محدودیت این بخش

این بخش hardening کامل systemd، SSH tunnel access، readiness endpoint یا graceful shutdown اضافه نکرد. آن موارد در بخش‌های بعدی انجام شدند.
