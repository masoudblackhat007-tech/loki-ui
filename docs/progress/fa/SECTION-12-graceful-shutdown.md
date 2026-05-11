# بخش ۱۲ — اضافه کردن graceful shutdown

## تاریخ

```text
2026-05-06
```

## هدف

در این بخش graceful shutdown برای `loki-ui` اضافه شد.

قبل از این بخش، HTTP server مستقیم با `ListenAndServe` اجرا می‌شد و هنگام stop یا restart شدن توسط systemd، shutdown path مشخصی نداشت.

## مشکل

سرویس تحت مدیریت systemd است و ممکن است هنگام deploy یا maintenance restart شود.

بدون graceful shutdown، process می‌تواند بدون مسیر کنترل‌شده متوقف شود.

## تغییر انجام‌شده

کد server طوری تغییر کرد که context دریافت کند و هنگام cancel شدن context، server با timeout مشخص shutdown شود.

همچنین channel برای دریافت error از `ListenAndServe` استفاده شد تا خطاهای واقعی server از shutdown عادی جدا شوند.

## نتیجه

در زمان restart یا stop، سرویس پیام shutdown request و shutdown completed را ثبت می‌کند و با مسیر کنترل‌شده متوقف می‌شود.

## نکته عملیاتی

این تغییر برای deployهای بعدی مهم بود، چون سرویس بارها build و restart شد.

## نکته امنیتی

graceful shutdown مدل امنیتی را تغییر نمی‌دهد، اما reliability سرویس داخلی را بهتر می‌کند.

## اعتبارسنجی

بعد از تغییر، build انجام شد، سرویس restart شد و وضعیت systemd بررسی شد.

سرویس همچنان active بود و روی loopback اجرا می‌شد.

## ارزش فنی

این بخش lifecycle مدیریت‌شده برای برنامه ایجاد کرد و رفتار سرویس را هنگام restart قابل پیش‌بینی‌تر کرد.

## ارزش رزومه‌ای قابل دفاع

```text
Implemented graceful shutdown for a systemd-managed Go HTTP service using context cancellation and bounded server shutdown to improve deploy and maintenance reliability.
```

## محدودیت این بخش

این بخش UI، API، Loki query، Alloy، authentication یا authorization را تغییر نداد.
