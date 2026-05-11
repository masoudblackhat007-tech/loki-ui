# بخش ۰۴ — اتصال به GitHub و اولین push

## تاریخ

```text
2026-05-05
```

## هدف

در این بخش repository پروژه `loki-ui` به GitHub وصل شد و branch اصلی `main` برای اولین بار push شد.

هدف این بود که پروژه از حالت فقط لوکال خارج شود و یک remote قابل استفاده برای deploy، server pull و تاریخچه قابل بررسی داشته باشد.

## کارهای انجام‌شده

در این بخش این کارها انجام شد:

```text
ساخت repository جدید GitHub برای loki-ui
اضافه کردن remote با نام origin
بررسی URL مربوط به remote
push کردن branch main
تنظیم tracking برای origin/main
بررسی clean بودن working tree بعد از push
اضافه کردن progress doc بخش ۰۴
```

## مدل remote

remote پروژه روی GitHub تنظیم شد تا سرور production بعداً فقط از repository رسمی pull کند.

این مدل برای deploy امن‌تر است، چون production نباید محل نوشتن کد application یا push کردن تغییرات باشد.

## وضعیت بعد از push

بعد از push اولیه، branch محلی `main` به `origin/main` وصل شد و working tree تمیز باقی ماند.

این وضعیت برای ادامه workflow مهم بود، چون تغییرات بعدی باید به شکل مرحله‌ای commit و push شوند.

## تصمیم امنیتی

هیچ فایل `.env` واقعی یا secret وارد GitHub نشد.

remote فقط شامل سورس و فایل‌های امن baseline بود.

## ارزش فنی

این بخش workflow توسعه را از حالت دستی و محلی به یک مسیر قابل تکرار تبدیل کرد:

```text
local development -> commit -> push to GitHub -> server pull/deploy
```

این مدل بعداً برای deployهای UI، chart و docs viewer استفاده شد.

## ارزش رزومه‌ای قابل دفاع

```text
Created and connected a GitHub remote for an internal Go-based Loki UI project, pushed the clean main branch, and established a repeatable local-to-remote workflow without committing secrets or generated artifacts.
```

## محدودیت این بخش

این بخش هنوز deploy سرور، systemd service، SSH tunnel یا اتصال عملی production به Loki را انجام نداد. این کارها در بخش‌های بعدی انجام شدند.
