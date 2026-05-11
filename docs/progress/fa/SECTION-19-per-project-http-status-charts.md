# بخش ۱۹ — اضافه کردن نمودار وضعیت HTTP برای هر پروژه

## تاریخ

```text
2026-05-11
```

## هدف

در این بخش صفحه `Requests` در پروژه `loki-ui` بهتر شد تا برای هر پروژه Laravel یک نمودار جداگانه از وضعیت HTTP requestها نمایش داده شود.

قبل از این بخش، صفحه `Requests` لاگ‌ها را بر اساس پروژه جدا می‌کرد و برای هر پروژه شمارنده‌های `2xx`، `4xx` و `5xx` را نشان می‌داد، اما روند زمانی requestهای موفق و خطادار مشخص نبود.

## مشکل

بعد از اضافه شدن پروژه دوم Laravel به pipeline لاگ، نمایش عددی به‌تنهایی کافی نبود.

کاربر می‌توانست ببیند هر پروژه چند request موفق یا ناموفق دارد، اما نمی‌توانست سریع تشخیص دهد این وضعیت‌ها در چه دقیقه‌هایی رخ داده‌اند یا آیا خطاها به شکل موجی و متمرکز ایجاد شده‌اند.

## تغییر انجام‌شده

برای هر project card یک نمودار `HTTP status trend` اضافه شد.

این نمودار بالای فیلترهای محلی همان پروژه و بالای جدول لاگ‌ها قرار گرفت.

نمودار به شکل `stacked bar chart` ساخته شد و وضعیت‌های `2xx`، `4xx` و `5xx` را برای هر پروژه جدا نشان می‌دهد.

هر ستون نمودار یک bucket زمانی یک‌دقیقه‌ای است.

## پیاده‌سازی

تغییر فقط در فایل زیر انجام شد:

```text
templates/logs.tmpl
```

برای نمودار از SVG داخلی استفاده شد و dependency خارجی اضافه نشد.

داده نمودار از همان rowهای visible داخل project card ساخته می‌شود؛ یعنی اگر فیلتر محلی پروژه تغییر کند، نمودار همان پروژه هم بر اساس همان خروجی visible دوباره ساخته می‌شود.

## رفتار نمودار

نمودار query جدید به Loki ارسال نمی‌کند.

نمودار فقط روی داده‌هایی کار می‌کند که قبلاً از `/api/logs` گرفته شده‌اند.

حداکثر ۱۲ bucket زمانی آخر نمایش داده می‌شود.

رنگ سبز وضعیت `2xx` را نشان می‌دهد.

رنگ نارنجی وضعیت `4xx` را نشان می‌دهد.

رنگ قرمز وضعیت `5xx` را نشان می‌دهد.

## تلاش ناموفق اولیه

در تلاش اولیه، فایل `templates/logs.tmpl` با ساختار اشتباه بازنویسی شد.

مشکل این بود که template با نام `logs` حذف شده بود و برنامه هنگام render صفحه `Requests` خطا داد.

خطای runtime این بود:

```text
html/template: "logs" is undefined
```

این خطا مربوط به Loki، Alloy، backend یا لاگ‌های Laravel نبود.

علت فقط خراب شدن ساختار template صفحه بود.

## بازیابی

برای recovery، فایل `templates/logs.tmpl` روی سرور از نسخه سالم قبلی برگردانده شد.

بعد از build و restart، سرویس دوباره active شد.

commit ناموفق این بود:

```text
b9acec1 Add per-project request status charts
```

این commit نباید به عنوان تغییر موفق معرفی شود.

## اصلاح نهایی

بعد از recovery، فایل `templates/logs.tmpl` دوباره با حفظ ساختار واقعی template اصلاح شد.

ساختار لازم template حفظ شد:

```gotemplate
{{ define "logs" }}
```

نسخه اصلاح‌شده نمودار SVG را درست render کرد و صفحه `Requests` دوباره بدون خطا بالا آمد.

commit موفق این بود:

```text
4f17fb0 Add per-project request status charts
```

## اعتبارسنجی

بعد از اصلاح، build روی سرور بدون خطا انجام شد.

سرویس `loki-ui` restart شد و وضعیت systemd فعال بود.

صفحه `Requests` از طریق SSH tunnel دیده شد و نمودار برای پروژه Laravel نمایش داده شد.

## نتیجه

در پایان این بخش، هر پروژه Laravel در صفحه `Requests` یک نمودار وضعیت HTTP مخصوص خودش دارد.

این نمودار روند زمانی requestهای `2xx`، `4xx` و `5xx` را نشان می‌دهد.

فیلترهای محلی، click-to-filter، clearable inputs، auto-refresh و layout responsive حفظ شدند.

## نکته امنیتی

این بخش مدل امنیتی را تغییر نداد.

پورت `18090` همچنان نباید public شود.

دسترسی به `loki-ui` همچنان باید فقط از طریق SSH tunnel باشد.

این بخش authentication، authorization، TLS، rate limiting یا audit logging اضافه نکرد.

## محدودیت این بخش

این نمودار آمار کامل تمام لاگ‌های موجود در Loki را نشان نمی‌دهد.

نمودار فقط روی داده‌هایی ساخته می‌شود که query فعلی از Loki گرفته و UI نمایش داده است.

این بخش anomaly detection یا AI analysis اضافه نکرد.

## ارزش فنی

این بخش خوانایی صفحه `Requests` را برای چند پروژه بهتر کرد و بدون تغییر backend، بدون تغییر LogQL و بدون اضافه کردن dependency، یک دید سریع زمانی از وضعیت HTTP هر پروژه ایجاد کرد.

## ارزش رزومه‌ای قابل دفاع

```text
Added frontend-only per-project HTTP status trend charts to an internal Loki-based Laravel request log viewer using SVG stacked bars for 2xx, 4xx, and 5xx counts while preserving the existing SSH-tunnel-only access model.
```
