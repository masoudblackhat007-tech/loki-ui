# بخش ۰۸ — اصلاح استخراج DTO از context تو در توی HTTP

## تاریخ

```text
2026-05-06
```

## هدف

در این بخش مشکل استخراج فیلدهای HTTP از لاگ‌های Laravel در `loki-ui` اصلاح شد.

API لاگ‌ها را از Loki می‌خواند، اما برخی فیلدهای top-level در DTO اشتباه یا خالی بودند.

## مشکل

فیلدهای زیر در خروجی درست پر نمی‌شدند:

```text
route: empty
method: empty
status: 0
duration_ms: 0
```

در حالی که مقدارهای درست داخل context تو در توی Laravel وجود داشتند:

```text
context.http.method
context.http.route
context.http.status_code
context.http.duration_ms
```

## علت

کد قبلی فقط به دنبال فیلدهای top-level می‌گشت و context nested را کامل پشتیبانی نمی‌کرد.

این باعث می‌شد UI داده را داشته باشد، اما در سطح DTO مقدارهای مهم نمایش داده نشوند.

## تغییر انجام‌شده

منطق استخراج DTO اصلاح شد تا ابتدا top-level و سپس context nested را بررسی کند.

helperهایی برای خواندن string و number از context اصلی و nested اضافه یا اصلاح شدند.

## نتیجه

بعد از اصلاح، فیلدهای HTTP مانند method، route، status و duration_ms در UI و API درست نمایش داده شدند.

## اعتبارسنجی

اعتبارسنجی با لاگ‌های واقعی Laravel انجام شد و خروجی API نشان داد که مقدارهای HTTP دیگر خالی یا صفر نیستند.

## نکته امنیتی

این تغییر فقط mapping داده را اصلاح کرد. هیچ secret، raw body، cookie، authorization header یا مقدار حساس جدیدی expose نشد.

## ارزش فنی

این بخش نشان داد که مشکل از Loki یا Laravel logging نبود؛ مشکل در mapping سمت UI بود و با تحلیل ساختار JSON لاگ اصلاح شد.

## ارزش رزومه‌ای قابل دفاع

```text
Fixed nested Laravel HTTP log context extraction in a Go-based Loki UI so method, route, status code, and duration fields are correctly mapped from structured JSON logs into API and UI DTOs.
```

## محدودیت این بخش

این بخش backend logging Laravel، Alloy، Loki، authentication یا security model را تغییر نداد.
