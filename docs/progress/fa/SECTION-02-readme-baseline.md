# بخش ۰۲ — README و مستندسازی baseline پروژه

## تاریخ

```text
2026-05-05
```

## هدف

در این بخش، README اولیه پروژه `loki-ui` نوشته شد تا هدف پروژه، معماری، مدل runtime، تنظیمات محیطی، مدل امنیتی، routeها، محدودیت‌ها و روش‌های verification از ابتدا مستند باشند.

هدف این نبود که پروژه بزرگ‌تر از واقعیت معرفی شود؛ هدف این بود که وضعیت واقعی و محدودیت‌های فعلی شفاف نوشته شوند.

## محتوای مستندشده

README اولیه این موارد را پوشش داد:

```text
هدف پروژه loki-ui
ارتباط با پروژه Laravel laravel-log2-loki
مدل خواندن لاگ‌ها از Loki
متغیرهای محیطی لازم
مدل localhost-only
routeهای موجود
قابلیت‌های فعلی
محدودیت‌های فعلی
دستورهای build و run
روش verification ساده
```

## مدل معماری مستندشده

در README توضیح داده شد که `loki-ui` یک ابزار داخلی Go برای خواندن لاگ‌های Laravel از Loki است.

مدل کلی در آن مرحله:

```text
Laravel JSON logs -> Alloy -> Loki -> loki-ui
```

## مدل امنیتی مستندشده

README تأکید کرد که UI نباید public باشد و باید فقط روی loopback اجرا شود.

مدل قابل قبول:

```text
browser -> SSH tunnel -> 127.0.0.1:18090 on server -> loki-ui -> Loki
```

این یعنی `loki-ui` در این مرحله جایگزین authentication، authorization، TLS، rate limiting یا audit logging نبود.

## متغیرهای محیطی

README به جای انتشار `.env` واقعی، از `.env.example` استفاده کرد.

این تصمیم از نظر امنیتی درست بود، چون مقدارهای واقعی runtime یا secretها نباید وارد مستند عمومی یا Git شوند.

## محدودیت‌های مستندشده

در README محدودیت‌های واقعی نوشته شد، از جمله اینکه:

```text
ابزار internal-only است
نباید به اینترنت expose شود
مدل authentication ندارد
مدل authorization ندارد
TLS در خود برنامه اضافه نشده
rate limiting و audit logging کامل وجود ندارد
```

## ارزش فنی

این بخش باعث شد پروژه از ابتدا قابل فهم، قابل build، قابل run و قابل دفاع باشد.

همچنین جلوی ادعاهای اشتباه گرفته شد؛ README پروژه را همان چیزی معرفی کرد که واقعاً بود، نه چیزی که هنوز ساخته نشده بود.

## نتیجه

پروژه یک README پایه گرفت که برای توسعه بعدی، deploy، بررسی security model و رزومه‌سازی قابل استفاده است.

## ارزش رزومه‌ای قابل دفاع

```text
Documented the baseline architecture, runtime configuration, localhost-only access model, current routes, limitations, and verification workflow for an internal Go-based Loki log viewer.
```

## محدودیت این بخش

این بخش هیچ feature جدید runtime اضافه نکرد. همچنین authentication، authorization، TLS، public deployment یا تغییر در Loki و Alloy انجام نشد.
