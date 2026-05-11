````markdown
# بخش ۱۹ — اضافه شدن نمودار وضعیت HTTP برای هر پروژه

## هدف بخش

در این بخش صفحه `Requests` در پروژه `loki-ui` به‌روزرسانی شد تا برای هر پروژه Laravel یک نمودار جداگانه از وضعیت HTTP requestها نمایش داده شود.

قبل از این بخش، صفحه `Requests` لاگ‌ها را بر اساس پروژه جدا می‌کرد و برای هر پروژه شمارنده‌های `2xx`، `4xx` و `5xx` را نشان می‌داد، اما روند زمانی این وضعیت‌ها مشخص نبود.

هدف این بخش این بود که کاربر بتواند برای هر پروژه، قبل از بررسی جدول لاگ‌ها، سریع ببیند requestهای موفق، خطاهای سمت client و خطاهای سمت server در چه دقیقه‌هایی رخ داده‌اند.

فایل اصلی تغییر داده‌شده در این بخش:

```text
templates/logs.tmpl
```

## محدوده تغییرات

تغییرات این بخش فقط در frontend انجام شد.

در این بخش فایل‌ها و بخش‌های زیر تغییر نکردند:

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

مدل دسترسی همچنان همان مدل قبلی باقی ماند:

```text
browser -> SSH tunnel -> 127.0.0.1:18090 on server -> loki-ui -> Loki
```

پورت `18090` همچنان نباید عمومی شود.

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

داخل هر project card شمارنده‌های زیر وجود داشت:

```text
2xx
4xx
5xx
```

اما فقط عدد نهایی دیده می‌شد.

برای مثال، اگر یک پروژه ۸۴ request داشت، UI نشان می‌داد چه تعداد `2xx`، چه تعداد `4xx` و چه تعداد `5xx` هستند، اما مشخص نبود این requestها در چه دقیقه‌هایی بیشتر شده‌اند.

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

نمودار بالای فیلترهای محلی و بالای جدول لاگ‌ها قرار گرفت تا کاربر قبل از ورود به جزئیات، وضعیت کلی همان پروژه را ببیند.

## نوع نمودار

نوع نمودار انتخاب‌شده:

```text
stacked bar chart
```

هر ستون نمودار نشان‌دهنده یک بازه زمانی یک‌دقیقه‌ای است.

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

تقسیم رنگی داخل ستون نشان می‌دهد از requestهای همان دقیقه، چه تعداد `2xx`، چه تعداد `4xx` و چه تعداد `5xx` بوده‌اند.

## رفتار نمودار

نمودار برای هر پروژه جدا محاسبه می‌شود.

داده نمودار از همان لاگ‌هایی ساخته می‌شود که داخل project card قابل مشاهده هستند.

این یعنی نمودار به فیلترهای محلی همان پروژه هم واکنش نشان می‌دهد.

اگر داخل یک project card فیلتر `Status` یا `Method` یا `Search this project` اعمال شود، نمودار همان پروژه فقط بر اساس لاگ‌های visible دوباره محاسبه می‌شود.

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

نمودار آخرین bucketهای زمانی را نمایش می‌دهد.

حداکثر تعداد bucketهای نمایش‌داده‌شده:

```text
12
```

یعنی نمودار حداکثر ۱۲ ستون زمانی آخر را نشان می‌دهد.

## توابع اضافه‌شده در frontend

برای ساخت نمودار، چند تابع به JavaScript داخل `templates/logs.tmpl` اضافه شد.

تابع ساخت bucketهای زمانی:

```text
buildStatusChartBuckets
```

وظیفه این تابع:

```text
read log timestamp
group logs by minute
count 2xx, 4xx, 5xx per minute
return last 12 buckets
```

تابع render نمودار:

```text
renderStatusChart
```

وظیفه این تابع:

```text
build SVG chart
draw chart grid
draw stacked bars
draw axis labels
draw legend
handle empty chart state
```

## محل render نمودار

نمودار داخل تابع render پروژه‌ها اضافه شد.

جایگاه آن داخل project body و قبل از local filters است.

ترتیب نهایی داخل هر project card:

```text
HTTP status trend chart
Search this project
Method
Status
Request table
```

## خروجی قابل مشاهده در UI

بعد از deploy، برای هر پروژه یک نمودار با عنوان زیر نمایش داده می‌شود:

```text
HTTP status trend
```

زیر عنوان نمودار توضیح زیر نمایش داده می‌شود:

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

اگر ستون نارنجی زیاد باشد، یعنی نسبت `4xx` بالا است.

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

این مقدار تحت تأثیر فیلترهای زیر است:

```text
Project label
Lookback range
Max rows fetched
Request ID
Loki text contains
```

همچنین تحت تأثیر فیلترهای محلی همان project card است:

```text
Search this project
Method
Status
```

بنابراین اگر مقدار `Max rows fetched` کم باشد، ممکن است بعضی لاگ‌ها اصلاً وارد UI نشوند و در نمودار هم دیده نشوند.

## تلاش ناموفق اولیه

در اولین تلاش برای اضافه کردن نمودار، فایل `templates/logs.tmpl` به شکل اشتباه بازنویسی شد.

اشتباه اصلی این بود که ساختار template واقعی پروژه حفظ نشد.

برنامه انتظار داشت template با نام زیر وجود داشته باشد:

```text
logs
```

اما نسخه اشتباه، template را با نام‌های دیگری تعریف کرده بود.

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

این تصمیم باعث شد تغییر فقط داخل frontend template باقی بماند و نیازی به build pipeline جدید، asset جدید یا package جدید نباشد.

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

## Commit نهایی

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

```bash
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

## جمع‌بندی رزومه‌ای قابل دفاع

ادعای قابل دفاع برای این بخش:

```text
Added per-project HTTP status trend charts to an internal Loki-based Laravel request log viewer, using frontend-only SVG stacked bar charts for 2xx, 4xx, and 5xx request counts while preserving the existing internal-only SSH tunnel access model and avoiding backend, Loki, Alloy, or Laravel logging changes.
```

این ادعا فقط به محدوده همین بخش مربوط است.

این تغییر نباید به عنوان تغییر backend، تغییر security model، اضافه شدن authentication، یا اضافه شدن AI analysis معرفی شود.
````
