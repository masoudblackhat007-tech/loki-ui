# بخش ۱۳ — اصلاح آلودگی محیط Go در WSL

## تاریخ

```text
2026-05-07
```

## هدف

در این بخش مشکل محیط Go داخل WSL اصلاح شد تا `go test` و `go build` از toolchain لینوکسی درست استفاده کنند.

## مشکل

در زمان validation، خطاهای زیر دیده شد:

```text
go: no such tool "vet"
go: no such tool "compile"
```

این خطاها از کد پروژه نبودند.

## علت

محیط WSL به path یا toolchain ویندوزی Go آلوده شده بود. در نتیجه Go داخل WSL ابزارهای لازم لینوکسی مثل `compile` و `vet` را درست پیدا نمی‌کرد.

## اصلاح انجام‌شده

محیط Go پاک‌سازی شد و pathها طوری اصلاح شدند که WSL از toolchain لینوکسی درست استفاده کند.

بعد از اصلاح، build و test دیگر به ابزارهای ناقص یا اشتباه اشاره نکردند.

## نتیجه

مشکل environment از مشکل application جدا شد.

این تصمیم مهم بود، چون اگر بدون تحلیل، سورس پروژه مقصر دانسته می‌شد، مسیر debug اشتباه می‌شد.

## ارزش فنی

این بخش نشان داد که validation قابل اعتماد فقط وقتی معنی دارد که toolchain سالم باشد.

## ارزش رزومه‌ای قابل دفاع

```text
Diagnosed and fixed a WSL Go environment contamination issue where Windows Go paths broke Linux go build and go test validation, separating toolchain failures from application code defects.
```

## محدودیت این بخش

این بخش تغییری در سورس application، UI، Loki، Alloy، systemd یا security model ایجاد نکرد.
