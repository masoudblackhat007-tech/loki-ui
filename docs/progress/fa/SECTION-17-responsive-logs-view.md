# بخش ۱۷ - ریسپانسیو کردن نمای لاگ‌ها


## تاریخ

2026-05-09

## هدف

بهبود نمایش صفحه `/logs` در عرض‌های کوچک‌تر، بدون تغییر در منطق Go، queryهای Loki، API، DTO یا مدل امنیتی پروژه.

## مشکل

صفحه `/logs` برای دسکتاپ طراحی شده بود. در عرض‌های کوچک، فیلترها فشرده می‌شدند، header انعطاف کافی نداشت، search input عرض ثابت داشت، summary bar خوب wrap نمی‌شد و جدول به‌جای اسکرول کنترل‌شده داخل wrapper با overflow پنهان می‌شد.

## علت اصلی

CSS صفحه فقط یک media query محدود برای مخفی‌کردن sidebar داشت و layout اصلی برای tablet/mobile کامل نشده بود.

## فایل‌های تغییرکرده

- `templates/logs.tmpl`

## اصلاح انجام‌شده

اصلاح فقط در CSS انجام شد: header و search منعطف‌تر شدند، فیلترها قابلیت wrap گرفتند، inputها عرض کامل داخل label گرفتند، summary bar قابل wrap شد، table wrapper اسکرول افقی کنترل‌شده گرفت و برای عرض‌های زیر `1100px` و `760px` رفتار responsive اضافه شد.

## اعتبارسنجی لوکال

دستورهای `go test ./...` و `go build -o bin/loki-ui ./cmd/loki-ui` در محیط لوکال WSL موفق بودند.

## اعتبارسنجی baseline سرور

قبل از deploy تغییر جدید، مسیر فعلی `/logs` از طریق SSH tunnel تست شد و پاسخ `HTTP/1.1 200 OK` برگشت.

## دستورهای deploy روی سرور

هنوز انجام نشده است.

## اعتبارسنجی روی سرور

هنوز انجام نشده است.

## بررسی امنیتی

این بخش فقط CSS را تغییر داد. LogQL، API، DTO، request logging، systemd، bind address، UFW و پردازش secretها تغییر نکردند. هیچ cookie، token، authorization header، request body، response body حساس، محتوای `.env` یا secret وارد سند نشد.

## محدودیت باقی‌مانده

جدول هنوز به card layout موبایلی تبدیل نشده است؛ در موبایل همچنان جدول باقی می‌ماند اما داخل wrapper اسکرول افقی کنترل‌شده دارد.

## نتیجه

نمای `/logs` در عرض‌های کوچک قابل استفاده‌تر شد، بدون اینکه منطق backend یا سطح امنیتی پروژه تغییر کند.

## ارزش رزومه‌ای
این بخش نشان می‌دهد که بهبود UI در ابزارهای observability باید کم‌ریسک، قابل تست و بدون تغییر ناخواسته در مسیر داده یا امنیت انجام شود.

## تکمیل نهایی بعد از deploy

بعد از بازبینی بصری، responsive اولیه کافی نبود؛ چون جدول روی موبایل هنوز تجربه‌ی مناسبی نداشت. اصلاح تکمیلی انجام شد و در موبایل، جدول به card layout تبدیل شد.

Commitهای مرتبط:

- `29ed046 Improve responsive logs view`
- `bffb2b5 Improve mobile logs card layout`
- `e40e6f5 Fix mobile logs view overflow`

اعتبارسنجی نهایی:

- `go test ./...` موفق بود.
- `go build -o bin/loki-ui ./cmd/loki-ui` موفق بود.
- سرویس `loki-ui` بعد از restart فعال بود.
- سرویس فقط روی `127.0.0.1:18090` گوش می‌داد.
- هیچ rule برای `18090/tcp` در UFW وجود نداشت.
- مسیر `/logs` از طریق SSH tunnel در مرورگر بررسی شد.
- در عرض حدود `390px`، requestها به شکل کارت عمودی با labelهای `Verb`، `Service`، `Path`، `Status` و `Duration` نمایش داده شدند.

نتیجه‌ی نهایی:

نمای `/logs` در موبایل از جدول فشرده و اسکرول افقی ضعیف به card layout قابل استفاده تبدیل شد، بدون تغییر در backend، API، LogQL، systemd، bind address یا firewall.
