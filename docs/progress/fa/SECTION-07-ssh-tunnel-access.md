# بخش ۰۷ — دسترسی به loki-ui از طریق SSH tunnel

## تاریخ

```text
2026-05-06
```

## هدف

در این بخش دسترسی امن به `loki-ui` از طریق SSH local port forwarding بررسی و مستند شد، بدون اینکه پورت `18090` به صورت عمومی باز شود.

هدف این بود که UI روی سرور internal بماند، اما از مرورگر لوکال قابل مشاهده باشد.

## مدل دسترسی

سرویس روی سرور فقط روی loopback گوش می‌دهد:

```text
127.0.0.1:18090
```

دسترسی از ماشین لوکال از طریق tunnel انجام می‌شود:

```text
local browser -> SSH tunnel -> server 127.0.0.1:18090 -> loki-ui -> Loki -> Laravel logs
```

## کارهای انجام‌شده

در این بخش این موارد بررسی شد:

```text
فعال بودن سرویس loki-ui روی سرور
گوش دادن سرویس روی 127.0.0.1:18090
عدم نیاز به باز کردن پورت 18090 در firewall
ساخت SSH tunnel از سیستم لوکال
باز کردن UI از طریق مرورگر لوکال
مشاهده لاگ‌ها از مسیر tunnel
```

## تصمیم امنیتی

پورت `18090` عمومی نشد.

این تصمیم حیاتی است، چون `loki-ui` در این مرحله authentication، authorization، TLS، rate limiting و audit logging کامل ندارد.

بنابراین expose کردن آن روی اینترنت اشتباه امنیتی است.

## نتیجه

کاربر توانست از طریق مرورگر لوکال UI را ببیند، در حالی که سرویس روی سرور فقط internal باقی ماند.

## ارزش فنی

این بخش access model درست برای ابزار observability داخلی را تثبیت کرد.

UI قابل استفاده شد، اما سطح attack surface عمومی اضافه نشد.

## ارزش رزومه‌ای قابل دفاع

```text
Implemented and validated SSH-tunneled access for an internal Go-based Loki UI, keeping the service bound to 127.0.0.1 on the server and avoiding public exposure of the observability interface.
```

## محدودیت این بخش

این بخش authentication، authorization، TLS، reverse proxy یا public access اضافه نکرد. دسترسی همچنان فقط internal و از طریق SSH tunnel قابل قبول است.
