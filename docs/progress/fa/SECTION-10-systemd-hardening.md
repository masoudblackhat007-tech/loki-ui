# بخش ۱۰ — سخت‌سازی systemd برای loki-ui

## تاریخ

```text
2026-05-06
```

## هدف

در این بخش سرویس systemd مربوط به `loki-ui` سخت‌سازی شد، بدون اینکه مدل دسترسی internal-only تغییر کند.

هدف کاهش سطح دسترسی runtime سرویس و حفظ workflow مبتنی بر SSH tunnel بود.

## شرط اصلی

سرویس باید همچنان فقط روی این آدرس گوش بدهد:

```text
127.0.0.1:18090
```

پورت زیر نباید در UFW باز شود:

```text
18090
```

## وضعیت شروع

قبل از این بخش، سرویس systemd وجود داشت و برنامه به صورت managed service اجرا می‌شد.

اما hardening آن کامل نبود و می‌شد محدودیت‌های بیشتری برای process تعریف کرد.

## تغییرات انجام‌شده

در service unit تنظیمات hardening اضافه یا اصلاح شد تا دسترسی سرویس محدودتر شود.

تمرکز روی این بود که برنامه فقط همان چیزی را داشته باشد که برای خواندن templateها، اجرای binary و اتصال به Loki لازم دارد.

## اصول امنیتی رعایت‌شده

در این بخش:

```text
سرویس public نشد
پورت 18090 در firewall باز نشد
مدل SSH tunnel حفظ شد
کد application مستقیم روی سرور تغییر نکرد
hardening با workflow فعلی سازگار نگه داشته شد
```

## اعتبارسنجی

بعد از تغییر service file، systemd reload شد، سرویس restart شد و وضعیت آن بررسی شد.

سرویس active باقی ماند و همچنان روی loopback گوش داد.

## نتیجه

سرویس `loki-ui` با محدودیت‌های systemd قوی‌تر اجرا شد، در حالی که دسترسی operational از طریق SSH tunnel خراب نشد.

## ارزش فنی

این بخش نشان داد که hardening نباید فقط روی کد باشد. runtime service هم باید محدود و قابل دفاع باشد.

## ارزش رزومه‌ای قابل دفاع

```text
Hardened the systemd runtime for an internal Go-based Loki UI while preserving a loopback-only SSH-tunneled access model and avoiding public exposure of the observability interface.
```

## محدودیت این بخش

این بخش authentication، authorization، TLS، rate limiting یا audit logging کامل اضافه نکرد. همچنین public reverse proxy ساخته نشد.
