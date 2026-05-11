# بخش ۱۵ — اضافه کردن endpoint readiness برای وابستگی Loki

## تاریخ

```text
2026-05-07
```

## هدف

در این بخش endpoint `/readyz` به `loki-ui` اضافه شد تا اتصال سرویس به وابستگی اصلی خود یعنی Loki بررسی شود.

این endpoint با `/healthz` فرق دارد.

## تفاوت health و readiness

`/healthz` پاسخ می‌دهد:

```text
آیا process HTTP برنامه زنده است؟
```

`/readyz` پاسخ می‌دهد:

```text
آیا loki-ui می‌تواند به Loki وصل شود؟
```

## رفتار endpoint

مسیر `/readyz` هنگام request، با timeout کوتاه readiness را از Loki بررسی می‌کند.

اگر Loki آماده باشد، پاسخ موفق برمی‌گردد.

اگر Loki در دسترس نباشد یا آماده نباشد، endpoint پاسخ `503 Service Unavailable` می‌دهد.

## تغییرات کد

برای پشتیبانی از readiness، client مربوط به Loki متد مناسب برای check کردن `/ready` یا وضعیت آماده بودن Loki دریافت کرد.

handler هم endpoint جدید را expose کرد.

## نتیجه

از این بخش به بعد می‌توان بین دو وضعیت فرق گذاشت:

```text
خود سرویس loki-ui زنده است
وابستگی Loki هم آماده است
```

## ارزش عملیاتی

این تفکیک برای debug بسیار مهم است. اگر `/healthz` موفق باشد ولی `/readyz` شکست بخورد، مشکل از خود process نیست و باید Loki یا اتصال داخلی بررسی شود.

## نکته امنیتی

این endpoint نباید به معنی public شدن UI باشد. مسیر readiness هم پشت همان مدل internal-only و SSH tunnel باقی می‌ماند.

## ارزش رزومه‌ای قابل دفاع

```text
Added a /readyz endpoint to an internal Go-based Loki UI to validate dependency readiness against Loki separately from process-level health checks.
```

## محدودیت این بخش

این بخش مدل security، firewall، authentication، authorization یا public access را تغییر نداد.
