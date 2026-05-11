درست می‌گویی؛ قبل از commit باید بخش ۱۹ فارسی را اصلاح کنیم، وگرنه commit ناقص می‌شود.

فایل زیر را کامل جایگزین کن:

`docs/progress/fa/SECTION-19-per-project-http-status-charts.md`

````markdown
# بخش ۱۹ — اضافه شدن نمودار وضعیت HTTP برای هر پروژه

## هدف بخش

در این بخش صفحه `Requests` در پروژه `loki-ui` به‌روزرسانی شد تا برای هر پروژه Laravel یک نمودار جداگانه از وضعیت HTTP requestها نمایش داده شود.

قبل از این بخش، صفحه `Requests` لاگ‌ها را بر اساس پروژه جدا می‌کرد و برای هر پروژه شمارنده‌های `2xx`، `4xx` و `5xx` را نشان می‌داد، اما روند زمانی این وضعیت‌ها مشخص نبود.

هدف این بخش این بود که قبل از بررسی جزئیات جدول لاگ‌ها، وضعیت کلی هر پروژه سریع دیده شود؛ یعنی مشخص باشد requestهای موفق، خطاهای سمت client و خطاهای سمت server در چه دقیقه‌هایی رخ داده‌اند.

## محدوده تغییرات

تغییرات این بخش فقط در frontend انجام شد.

فایل اصلی تغییر داده‌شده:

```text
templates/logs.tmpl
```

در این بخش این قسمت‌ها تغییر نکردند:

```text
internal/httpserver/handler.go
internal/httpserver/server.go
internal/loki/client.go
internal/loki/types.go
Alloy configuration
Loki configuration
Laravel logging code
systemd service file
server firewall configuration
SSH tunnel model
```

در این بخش هیچ تغییری در backend، API، LogQL، Alloy، Loki یا مدل امنیتی انجام نشد.

این بخش فقط یک بهبود UX برای صفحه `Requests` بود.

## مدل دسترسی

مدل دسترسی `loki-ui` تغییر نکرد.

برنامه همچنان باید فقط روی loopback سرور گوش بدهد:

```text
127.0.0.1:18090
```

مدل دسترسی همچنان این است:

```text
browser -> SSH tunnel -> 127.0.0.1:18090 on server -> loki-ui -> Loki
```

پورت `18090` نباید عمومی شود.

این بخش نباید به عنوان تغییر در مدل امنیتی، public exposure، authentication، authorization، TLS، rate limiting یا audit logging معرفی شود.

## وضعیت قبل از شروع

در پایان بخش ۱۸، صفحه `Requests` این قابلیت‌ها را داشت:

```text
project-based grouping
per-project local filters
click-to-filter
clearable inputs
focus-preserving re-render
server-side Loki filters
auto-refresh
responsive layout
```

هر پروژه در یک card جدا نمایش داده می‌شد.

داخل هر project card شمارنده‌های جداگانه برای وضعیت‌های زیر وجود داشت:

```text
2xx
4xx
5xx
```

اما فقط عدد نهایی دیده می‌شد.

برای مثال، اگر یک پروژه ۸۴ request داشت، UI نشان می‌داد چه تعداد `2xx`، چه تعداد `4xx` و چه تعداد `5xx` هستند، اما مشخص نبود این requestها در چه دقیقه‌هایی رخ داده‌اند.

## تصمیم UX

برای هر project card یک نمودار مستقل اضافه شد.

جایگاه نمودار داخل هر پروژه:

```text
project card
  -> project header
  -> HTTP status trend chart
  -> project local filters
  -> request log table
```

نمودار بالای فیلترهای محلی و بالای جدول لاگ‌ها قرار گرفت.

این جایگاه باعث می‌شود کاربر قبل از ورود به جزئیات، وضعیت کلی همان پروژه را ببیند.

## نوع نمودار

نوع نمودار انتخاب‌شده:

```text
stacked bar chart
```

هر ستون نمودار نشان‌دهنده یک bucket زمانی یک‌دقیقه‌ای است.

هر ستون از سه بخش ساخته می‌شود:

```text
2xx
4xx
5xx
```

معنی رنگ‌ها:

```text
Green  -> 2xx
Orange -> 4xx
Red    -> 5xx
```

ارتفاع کل ستون یعنی تعداد کل requestهای همان دقیقه.

تقسیم رنگی داخل ستون نشان می‌دهد از requestهای همان دقیقه، چه تعداد موفق، چه تعداد خطای client-side و چه تعداد خطای server-side بوده‌اند.

## رفتار نمودار

نمودار برای هر پروژه جدا محاسبه می‌شود.

داده نمودار از همان لاگ‌هایی ساخته می‌شود که داخل همان project card قابل مشاهده هستند.

این یعنی نمودار به فیلترهای محلی همان پروژه واکنش نشان می‌دهد.

اگر داخل یک project card فیلترهای زیر اعمال شوند، نمودار همان پروژه دوباره بر اساس خروجی visible محاسبه می‌شود:

```text
Search this project
Method
Status
```

نمودار query جدید به Loki ارسال نمی‌کند.

نمودار فقط روی داده‌هایی کار می‌کند که قبلاً از endpoint زیر دریافت شده‌اند:

```text
/api/logs
```

## منطق bucket زمانی

برای ساخت نمودار، timestamp هر log خوانده می‌شود.

هر log داخل bucket یک‌دقیقه‌ای خودش قرار می‌گیرد.

برای هر bucket این مقدارها محاسبه می‌شوند:

```text
2xx count
4xx count
5xx count
total count
```

نمودار حداکثر آخرین ۱۲ bucket زمانی را نمایش می‌دهد.

یعنی نمودار حداکثر ۱۲ ستون آخر را نشان می‌دهد.

## توابع frontend اضافه‌شده

برای ساخت نمودار، چند تابع JavaScript داخل `templates/logs.tmpl` اضافه شد.

تابع ساخت bucketهای زمانی:

```text
buildStatusChartBuckets
```

وظیفه این تابع:

```text
خواندن timestamp لاگ
گروه‌بندی لاگ‌ها بر اساس دقیقه
شمردن 2xx، 4xx و 5xx در هر دقیقه
برگرداندن آخرین ۱۲ bucket
```

تابع render نمودار:

```text
renderStatusChart
```

وظیفه این تابع:

```text
ساخت SVG chart
رسم grid نمودار
رسم stacked bars
رسم axis labels
رسم legend
نمایش empty state
```

## محل render نمودار

نمودار داخل render هر project card اضافه شد.

ترتیب نهایی داخل هر project card:

```text
HTTP status trend chart
Search this project
Method
Status
Request table
```

این ساختار باعث شد نمودار به عنوان خلاصه دیداری همان پروژه عمل کند، نه یک نمودار global جدا از context.

## خروجی قابل مشاهده در UI

بعد از deploy، برای هر پروژه یک نمودار با عنوان زیر نمایش داده می‌شود:

```text
HTTP status trend
```

زیر عنوان نمودار این توضیح نمایش داده می‌شود:

```text
Stacked per-minute request counts for the currently visible rows.
```

در legend نمودار این سه وضعیت دیده می‌شود:

```text
2xx
4xx
5xx
```

## نمونه تفسیر نمودار

اگر در یک project card مقدارهای زیر دیده شود:

```text
Showing 84 of 84 logs
2xx: 42
4xx: 42
5xx: 0
```

معنی آن این است که برای آن پروژه، از لاگ‌های fetch شده و visible فعلی:

```text
42 request موفق بوده‌اند.
42 request خطای client-side داشته‌اند.
0 request خطای server-side داشته‌اند.
```

اگر در نمودار رنگ قرمز دیده نشود، یعنی در داده‌های visible فعلی `5xx` وجود ندارد.

اگر ستون‌های نارنجی زیاد باشند، یعنی نسبت `4xx` بالا است.

اگر `4xx`ها مربوط به تست‌های عمدی باشند، این وضعیت طبیعی است.

اگر `4xx`ها واقعی باشند، می‌تواند نشانه یکی از این موارد باشد:

```text
missing route
broken link
wrong client request
bot scan
probing
invalid endpoint
```

## نکته مهم درباره معنی نمودار

این نمودار تعداد کل واقعی لاگ‌های موجود در Loki را نشان نمی‌دهد.

نمودار فقط روی لاگ‌هایی ساخته می‌شود که query فعلی از Loki گرفته و UI دریافت کرده است.

این مقدار تحت تأثیر فیلترهای server-side است:

```text
Project label
Lookback range
Max rows fetched
Request ID
Loki text contains
```

همچنین تحت تأثیر فیلترهای local همان project card است:

```text
Search this project
Method
Status
```

اگر مقدار `Max rows fetched` کم باشد، ممکن است بعضی لاگ‌ها اصلاً وارد UI نشوند و در نمودار هم دیده نشوند.

پس این نمودار ابزار تحلیل سریع روی داده‌های visible فعلی است، نه آمار کامل تمام لاگ‌های موجود در Loki.

## تلاش ناموفق اولیه

در اولین تلاش برای اضافه کردن نمودار، فایل `templates/logs.tmpl` به شکل اشتباه بازنویسی شد.

اشتباه اصلی این بود که ساختار template واقعی پروژه حفظ نشد.

برنامه انتظار داشت template با نام زیر وجود داشته باشد:

```text
logs
```

اما نسخه اشتباه template را با نام درست تعریف نکرده بود.

خطای runtime:

```text
render error: html/template: "logs" is undefined
```

این خطا مربوط به Loki، Alloy، backend یا لاگ‌های Laravel نبود.

علت خطا فقط خراب شدن ساختار template صفحه `Requests` بود.

برای recovery، فایل `templates/logs.tmpl` روی سرور از نسخه سالم قبلی برگردانده شد.

بعد از build و restart، سرویس دوباره فعال شد.

Commit ناموفق:

```text
b9acec1 Add per-project request status charts
```

این commit نباید به عنوان تغییر موفق مستند شود.

## اصلاح نهایی

بعد از recovery، فایل `templates/logs.tmpl` دوباره با حفظ ساختار واقعی template بازنویسی شد.

نسخه اصلاح‌شده با این ساختار شروع و تمام می‌شود:

```gotemplate
{{ define "logs" }}
...
{{ end }}
```

در نسخه اصلاح‌شده، نمودار واقعی با SVG ساخته شد.

برای نمودار از dependency خارجی استفاده نشد.

این تصمیم باعث شد تغییر فقط داخل frontend template باقی بماند و نیازی به package جدید، build pipeline جدید یا asset خارجی نباشد.

## اعتبارسنجی لوکال

قبل از commit نهایی، وضعیت Git بررسی شد.

فقط فایل زیر تغییر کرده بود:

```text
templates/logs.tmpl
```

تلاش برای اجرای مستقیم برنامه بدون `LOKI_URL` با خطای environment متوقف شد:

```text
panic: LOKI_URL is required
```

این خطا مربوط به template نبود.

بعد از دادن `LOKI_URL`، اجرای لوکال به دلیل اشغال بودن port پیش‌فرض متوقف شد:

```text
listen tcp 127.0.0.1:18090: bind: address already in use
```

به دلیل محدودیت runtime لوکال، تست نهایی UI روی سرور انجام شد.

## Commit موفق

بعد از اصلاح فایل، commit نهایی ساخته شد.

Commit موفق:

```text
4f17fb0 Add per-project request status charts
```

این commit روی GitHub push شد.

## Deploy روی سرور

Deploy روی سرور در مسیر زیر انجام شد:

```text
/home/deploy/apps/loki-ui
```

فرآیند deploy شامل این مراحل بود:

```text
git reset --hard origin/main
git pull --ff-only origin main
go build -o bin/loki-ui ./cmd/loki-ui
sudo systemctl restart loki-ui
sudo systemctl status loki-ui --no-pager
```

بعد از pull، سرور از commit خراب قبلی به commit نهایی رسید:

```text
b9acec1 -> 4f17fb0
```

Build بدون خطا انجام شد.

سرویس restart شد.

وضعیت نهایی سرویس:

```text
Active: active (running)
```

## نتیجه نهایی

در پایان این بخش، صفحه `Requests` این قابلیت جدید را دارد:

```text
per-project HTTP status trend chart
```

برای هر پروژه، یک نمودار جدا بالای جدول لاگ‌ها نمایش داده می‌شود.

نمودار برای هر پروژه وضعیت‌های زیر را نشان می‌دهد:

```text
2xx
4xx
5xx
```

نمودار به صورت stacked bar chart نمایش داده می‌شود.

هر ستون مربوط به یک دقیقه است.

نمودار بر اساس لاگ‌های visible همان project card ساخته می‌شود.

فیلترهای محلی هر پروژه روی نمودار همان پروژه اثر می‌گذارند.

قابلیت‌های قبلی حفظ شدند:

```text
project grouping
project summary
local project filters
click-to-filter
clearable inputs
focus-preserving render
server-side Loki filters
auto-refresh
responsive layout
```

## نکات امنیتی

در این بخش هیچ تغییری در مدل امنیتی انجام نشد.

موارد زیر همچنان برقرار هستند:

```text
loki-ui must stay internal-only
port 18090 must not be opened publicly
access must stay behind SSH tunnel
observability UI must not be exposed to the internet
```

این بخش نباید به عنوان authentication، authorization، TLS، rate limiting، audit logging یا hardening جدید معرفی شود.

این بخش فقط یک بهبود UX برای تحلیل سریع‌تر لاگ‌های HTTP است.

## محدودیت‌های فعلی

نمودار فعلی فقط روی داده‌های fetch شده در UI کار می‌کند.

نمودار فعلی count واقعی کل لاگ‌های Loki را نشان نمی‌دهد.

نمودار فعلی anomaly detection انجام نمی‌دهد.

نمودار فعلی AI-based analysis ندارد.

نمودار فعلی فقط وضعیت‌های HTTP را در سه گروه زیر نمایش می‌دهد:

```text
2xx
4xx
5xx
```

## ارزش فنی این بخش

ارزش این بخش فقط اضافه شدن یک نمودار نبود.

ارزش اصلی این بود که UI بتواند برای هر پروژه، روند زمانی وضعیت requestها را بدون تغییر backend و بدون سنگین‌تر کردن queryها نشان دهد.

این کار باعث شد صفحه `Requests` برای چند پروژه Laravel قابل تحلیل‌تر شود.

## ارزش رزومه‌ای قابل دفاع

ادعای قابل دفاع برای این بخش:

```text
Added per-project HTTP status trend charts to an internal Loki-based Laravel request log viewer, using frontend-only SVG stacked bar charts for 2xx, 4xx, and 5xx request counts while preserving the existing internal-only SSH tunnel access model and avoiding backend, Loki, Alloy, or Laravel logging changes.
```

این ادعا فقط به محدوده همین بخش مربوط است.

این تغییر نباید به عنوان تغییر backend، تغییر security model، اضافه شدن authentication، اضافه شدن authorization، اضافه شدن TLS، اضافه شدن audit logging یا اضافه شدن AI analysis معرفی شود.
````
