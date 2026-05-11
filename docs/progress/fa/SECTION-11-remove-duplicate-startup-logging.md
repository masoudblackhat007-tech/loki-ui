# بخش ۱۱ — حذف لاگ تکراری startup

## تاریخ

```text
2026-05-06
```

## هدف

در این بخش لاگ تکراری startup در `loki-ui` حذف شد.

بعد از تغییرات HTTP server در بخش ۰۹، پیام شروع سرویس دو بار چاپ می‌شد.

## مشکل

در journal این پیام دوبار دیده می‌شد:

```text
loki-ui listening on 127.0.0.1:18090
loki-ui listening on 127.0.0.1:18090
```

این وضعیت می‌توانست گمراه‌کننده باشد و این تصور را ایجاد کند که دو process یا دو socket وجود دارد.

## علت

مشکل port binding نبود.

فقط startup message در دو جای مختلف کد log می‌شد.

## اصلاح انجام‌شده

یکی از محل‌های logging حذف شد تا فقط یک پیام startup معتبر باقی بماند.

## نتیجه

بعد از اصلاح، هنگام شروع سرویس فقط یک پیام listen ثبت شد.

این باعث شد journal تمیزتر و تشخیص وضعیت سرویس ساده‌تر شود.

## ارزش فنی

این بخش کوچک بود، اما از نظر operational hygiene مهم است. لاگ‌های تکراری در محیط واقعی باعث تحلیل اشتباه می‌شوند.

## ارزش رزومه‌ای قابل دفاع

```text
Removed duplicate startup logging from a systemd-managed Go service to keep operational logs accurate and avoid misleading service diagnostics.
```

## محدودیت این بخش

این بخش رفتار functional برنامه، Loki query، Alloy، security model یا access model را تغییر نداد.
