# بخش ۱۸ — بهبود تجربه کاربری فیلترهای صفحه Requests

## هدف بخش

در این بخش رابط کاربری صفحه `Requests` در پروژه `loki-ui` بهتر شد تا بررسی لاگ‌های HTTP چند پروژه Laravel ساده‌تر، سریع‌تر و قابل‌فهم‌تر شود.

قبل از این بخش، صفحه `Requests` لاگ‌ها را از Loki دریافت می‌کرد، اما وقتی چند پروژه Laravel هم‌زمان لاگ تولید می‌کردند، تشخیص پروژه، فیلتر کردن لاگ‌ها، و رسیدن به درخواست موردنظر سخت‌تر از چیزی بود که باید باشد.

هدف این بخش این بود که بدون تغییر در backend، بدون تغییر در Loki، بدون تغییر در Alloy، بدون تغییر در Laravel logging، و بدون تغییر در مدل امنیتی، فقط تجربه کاربری صفحه لاگ‌ها بهتر شود.

فایل اصلی تغییر داده‌شده در این بخش:

```text
templates/logs.tmpl
```

## محدوده تغییرات

تغییرات این بخش فقط روی رابط کاربری صفحه `Requests` انجام شد.

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

مدل دسترسی همچنان همان مدل قبلی باقی ماند:

```text
browser -> SSH tunnel -> 127.0.0.1:18090 on server -> loki-ui -> Loki
```

پورت `18090` همچنان نباید عمومی شود.

## وضعیت قبل از شروع

صفحه `Requests` لاگ‌های HTTP را از endpoint زیر دریافت می‌کرد:

```text
/api/logs
```

چند مشکل UX وجود داشت:

- لاگ‌های چند پروژه Laravel در یک نمای کلی سخت‌تر قابل بررسی بودند.
- تشخیص اینکه هر لاگ متعلق به کدام پروژه است به اندازه کافی سریع نبود.
- فیلتر کردن داخل یک پروژه مشخص راحت نبود.
- کاربر مجبور بود مقدارهایی مثل `path`، `service`، `method`، `status` یا `time` را دستی کپی کند.
- بعضی فیلترهای بالای صفحه کاربرد واضحی نداشتند.
- تفاوت فیلترهای server-side و فیلترهای client-side واضح نبود.
- هنگام تایپ داخل بعضی inputها، به دلیل re-render شدن DOM، focus از input خارج می‌شد.
- معنی `service_name` و `limit` در UI برای کاربر واضح نبود.
- کلیک روی مقدار `Service` ابتدا داخل فیلتر همان پروژه اعمال می‌شد، در حالی که باید فیلتر سراسری Loki را تغییر می‌داد.

## تغییر اول — گروه‌بندی لاگ‌ها بر اساس پروژه

صفحه `Requests` طوری تغییر کرد که لاگ‌ها بر اساس پروژه جدا شوند.

مبنای اصلی تشخیص پروژه مقدار label زیر است:

```text
service_name
```

این مقدار از Loki labels خوانده می‌شود.

اگر `service_name` وجود نداشته باشد، UI تلاش می‌کند از مقدارهای جایگزین استفاده کند:

```text
log.service
context.service
context.project
labels.service
labels.project
```

اگر هیچ مقدار معتبری پیدا نشود، مقدار پیش‌فرض زیر نمایش داده می‌شود:

```text
unknown-project
```

بعد از این تغییر، هر پروژه در یک card جدا نمایش داده می‌شود.

Commit مربوط:

```text
2676526 Improve project-based request log view
```

## تغییر دوم — اضافه شدن فیلترهای محلی برای هر پروژه

برای هر project card سه فیلتر محلی اضافه شد:

```text
Search this project
Method
Status
```

این فیلترها فقط روی همان project card اعمال می‌شوند.

این یعنی اگر چند پروژه Laravel هم‌زمان لاگ داشته باشند، می‌توان داخل لاگ‌های یک پروژه خاص جستجو کرد، بدون اینکه لاگ‌های پروژه‌های دیگر در نتیجه دخالت کنند.

این فیلترها سمت مرورگر هستند و فقط روی لاگ‌هایی اعمال می‌شوند که قبلاً از Loki گرفته شده‌اند.

## تغییر سوم — کلیک روی مقدارهای مهم برای فیلتر سریع

برای کاهش تایپ دستی، چند ستون قابل کلیک شدند.

مقدارهای زیر قابل کلیک شدند:

```text
Service
Path
Time
Verb
Status
```

رفتار نهایی هرکدام به شکل زیر تنظیم شد:

| ستون | رفتار نهایی |
|---|---|
| `Service` | مقدار label پروژه را داخل فیلتر سراسری `Project label` قرار می‌دهد و query جدید به Loki می‌فرستد |
| `Path` | مقدار مسیر را داخل فیلتر محلی همان پروژه قرار می‌دهد |
| `Time` | مقدار زمان نمایش‌داده‌شده را داخل فیلتر محلی همان پروژه قرار می‌دهد |
| `Verb` | مقدار method مثل `GET` را داخل فیلتر `Method` همان پروژه قرار می‌دهد |
| `Status` | مقدار status مثل `200` یا `404` را داخل فیلتر `Status` همان پروژه قرار می‌دهد |

علت جدا بودن رفتار `Service` این است که مقدار پروژه باید query اصلی Loki را محدود کند، نه فقط لاگ‌های همان card را.

## تغییر چهارم — اصلاح رفتار کلیک روی Service

در نسخه میانی این بخش، کلیک روی مقدار `Service` مقدار را داخل فیلتر محلی همان پروژه قرار می‌داد.

این رفتار درست نبود، چون `Service` در این UI عملاً به `service_name` در Loki مربوط است و باید فیلتر server-side را تغییر دهد.

رفتار اصلاح شد:

```text
Click Service -> fill Project label -> fetch /api/logs again
```

برای این کار مقدار واقعی label از این مسیر خوانده می‌شود:

```text
labels.service_name
```

اگر این مقدار وجود نداشته باشد، از مقدار نمایش‌داده‌شده سرویس استفاده می‌شود.

Commit مربوط:

```text
e442aa5 Fix request log service filter behavior
```

## تغییر پنجم — اضافه شدن دکمه پاک‌کردن داخل inputها

برای همه inputهای قابل استفاده، دکمه ضربدر اضافه شد.

این دکمه وقتی مقدار input خالی باشد نمایش داده نمی‌شود و وقتی input مقدار داشته باشد ظاهر می‌شود.

با کلیک روی ضربدر:

- مقدار input پاک می‌شود.
- رویدادهای `input` و `change` اجرا می‌شوند.
- focus داخل همان input باقی می‌ماند.

این تغییر برای فیلترهای بالای صفحه و فیلترهای داخل project card اعمال شد.

Commit مربوط:

```text
74e49dc Add clear buttons to request log filters
```

## تغییر ششم — جلوگیری از خروج focus هنگام تایپ

یک مشکل UX مهم این بود که هنگام تایپ داخل فیلترهای project card، بعد از وارد کردن یک حرف، DOM دوباره render می‌شد و focus از input خارج می‌شد.

برای حل این مشکل، وضعیت focus ذخیره و بعد از render دوباره برگردانده شد.

اطلاعاتی که قبل از render ذخیره می‌شود:

```text
project key
filter name
selection start
selection end
```

بعد از render، همان input دوباره پیدا می‌شود و focus و cursor position برگردانده می‌شود.

برای جلوگیری از render بیش از حد، debounce هم اضافه شد.

مقدار delay برای render محلی:

```text
180ms
```

Commit مربوط:

```text
50a37d0 Improve request log filter interactions
```

## تغییر هفتم — حذف جستجوی کلی بالای صفحه

ابتدا یک global search بالای همه پروژه‌ها وجود داشت.

این فیلتر روی همه پروژه‌ها اعمال می‌شد، اما از نظر UX واضح نبود و با مدل جدید project cardها تداخل ذهنی ایجاد می‌کرد.

چون هدف اصلی این صفحه بررسی جداگانه پروژه‌ها بود، global search حذف شد.

بعد از حذف آن، جستجو به دو سطح واضح تقسیم شد:

```text
server-side filters
project-level filters
```

فیلترهای server-side داده را از Loki محدود می‌کنند.

فیلترهای project-level فقط روی داده‌های گرفته‌شده داخل همان card اعمال می‌شوند.

Commit مربوط:

```text
659b883 Remove unused global request log filters
```

## تغییر هشتم — شفاف‌سازی فیلترهای بالای صفحه

فیلترهای بالای صفحه بازنویسی شدند تا مشخص باشد این‌ها فیلترهای سمت سرور هستند.

عنوان اضافه‌شده:

```text
Server-side Loki filters
```

توضیح اضافه‌شده:

```text
These filters change the query sent to Loki. Project filters inside each card only filter the already fetched rows.
```

نام فیلدها واضح‌تر شد:

| نام قبلی | نام جدید |
|---|---|
| `Service label` | `Project label` |
| `Range` | `Lookback range` |
| `Limit` | `Max rows fetched` |
| `Text` | `Loki text contains` |

placeholder فیلدها هم واضح‌تر شد.

برای `Project label`:

```text
service_name, e.g. laravel-log3-loki
```

برای `Request ID`:

```text
exact correlation id
```

برای `Loki text contains`:

```text
server-side raw log text filter
```

برای `Max rows fetched`:

```text
maximum rows returned by Loki
```

Commit مربوط:

```text
955c482 Clarify request log server filters
```

## معنی فیلترهای server-side

فیلترهای server-side مستقیماً روی query ارسالی به Loki اثر می‌گذارند.

این فیلترها برای کم‌کردن حجم داده‌ای هستند که از Loki گرفته می‌شود.

### Project label

این فیلتر مقدار label زیر را محدود می‌کند:

```text
service_name
```

نمونه مقدارها:

```text
laravel-log2-loki
laravel-log3-loki
```

وقتی روی badge سرویس کلیک شود، مقدار `service_name` داخل این input قرار می‌گیرد و لاگ‌ها دوباره از Loki دریافت می‌شوند.

### Lookback range

این فیلتر مشخص می‌کند Loki از چه بازه زمانی لاگ‌ها را بخواند.

گزینه‌های فعلی:

```text
Last 15 minutes
Last 1 hour
Last 6 hours
Last 24 hours
```

### Max rows fetched

این فیلتر مشخص می‌کند حداکثر چند ردیف از Loki گرفته شود.

این فیلتر تعداد واقعی لاگ‌های کل سیستم را نشان نمی‌دهد؛ فقط سقف تعداد لاگ‌هایی است که query فعلی از Loki برمی‌گرداند.

اگر مقدار خیلی کم باشد، ممکن است بعضی لاگ‌های مرتبط دیده نشوند.

اگر مقدار خیلی زیاد باشد، UI سنگین‌تر می‌شود و query می‌تواند کندتر شود.

### Request ID

این فیلتر برای پیدا کردن یک request مشخص استفاده می‌شود.

مقدار آن باید correlation id دقیق باشد.

از این مقدار برای باز کردن جزئیات یک request و دنبال کردن لاگ‌های مرتبط با همان request استفاده می‌شود.

### Loki text contains

این فیلتر روی متن خام لاگ در Loki اعمال می‌شود.

این فیلتر client-side نیست و query ارسالی به Loki را محدود می‌کند.

برای جستجوی سریع یک عبارت خام در لاگ‌ها استفاده می‌شود.

## معنی فیلترهای داخل هر پروژه

فیلترهای داخل هر project card فقط روی همان پروژه اعمال می‌شوند.

این فیلترها query جدید به Loki نمی‌فرستند.

فیلترهای داخل project card:

```text
Search this project
Method
Status
```

### Search this project

این فیلتر روی داده‌های همان پروژه جستجو می‌کند.

مقدارهای قابل جستجو شامل موارد زیر است:

```text
message
path
service
method
status
formatted time
service_name
error text
context JSON
```

### Method

این فیلتر مقدار HTTP method را دقیق بررسی می‌کند.

نمونه مقدار:

```text
GET
POST
```

### Status

این فیلتر HTTP status code را دقیق بررسی می‌کند.

نمونه مقدار:

```text
200
404
500
```

## تغییرات مرتبط با status summary

در project cardها، شمارنده‌های وضعیت HTTP نمایش داده می‌شوند:

```text
2xx
4xx
5xx
```

این شمارنده‌ها برای هر پروژه جدا محاسبه می‌شوند.

در summary بالای صفحه هم وضعیت کلی لاگ‌های HTTP fetched شده نمایش داده می‌شود.

متن summary به شکل واضح‌تر تغییر کرد:

```text
Total fetched HTTP entries
```

## تغییرات مرتبط با auto-refresh

قابلیت auto-refresh حفظ شد.

دکمه pause/resume همچنان وجود دارد.

دکمه refresh دستی همچنان وجود دارد.

بازه refresh همان مدل قبلی باقی ماند:

```text
3000ms
```

یعنی صفحه هر ۳ ثانیه یک بار، در صورت pause نبودن، لاگ‌ها را دوباره fetch می‌کند.

## نکات امنیتی

در این بخش هیچ تغییری در مدل امنیتی انجام نشد.

موارد زیر همچنان باید رعایت شوند:

- سرویس `loki-ui` نباید public شود.
- پورت `18090` نباید در UFW باز شود.
- دسترسی باید از طریق SSH tunnel باقی بماند.
- مقدارهای حساس نباید در docs، screenshots، resume، commit message یا issueها منتشر شوند.
- مقدارهای raw مربوط به `request_id`، `session_hash`، token، cookie، authorization header، raw headers، raw bodies، private keys، API keys و `.env` نباید در مستندات عمومی کپی شوند.

این بخش فقط UX را بهتر کرده و نباید به عنوان authentication، authorization، audit logging، rate limiting یا hardening جدید معرفی شود.

## تست محلی

بعد از تغییرات، تست‌های Go در WSL اجرا شدند.

دستور اجراشده:

```bash
go test ./...
```

خروجی موفق:

```text
?       loki-ui/cmd/loki-ui             [no test files]
?       loki-ui/internal/httpserver     [no test files]
?       loki-ui/internal/loki           [no test files]
```

این تست نشان داد تغییرات template باعث شکست build/testهای فعلی Go نشده‌اند.

## بررسی diff محلی

برای بررسی حجم تغییرات از دستور زیر استفاده شد:

```bash
git diff --stat
```

در طول این بخش چند مرحله تغییر روی فایل زیر انجام شد:

```text
templates/logs.tmpl
```

تغییرات در چند commit کوچک و قابل پیگیری ثبت شدند.

## Commitهای این بخش

Commitهای انجام‌شده در این بخش:

```text
2676526 Improve project-based request log view
701fcec Add click-to-filter request log fields
74e49dc Add clear buttons to request log filters
50a37d0 Improve request log filter interactions
e442aa5 Fix request log service filter behavior
659b883 Remove unused global request log filters
955c482 Clarify request log server filters
```

## Deploy روی سرور

بعد از push شدن تغییرات، روی سرور deploy انجام شد.

مسیر پروژه روی سرور:

```text
/home/deploy/apps/loki-ui
```

فرآیند deploy شامل این مراحل بود:

```bash
git pull --ff-only origin main
go build -o bin/loki-ui ./cmd/loki-ui
sudo systemctl restart loki-ui
systemctl status loki-ui --no-pager
```

بعد از restart، سرویس فعال بود.

نمونه وضعیت موفق سرویس:

```text
Active: active (running)
```

## نتیجه نهایی

در پایان این بخش، صفحه `Requests` این قابلیت‌ها را دارد:

- لاگ‌ها را بر اساس پروژه جدا نمایش می‌دهد.
- هر پروژه فیلترهای محلی خودش را دارد.
- کلیک روی `Path` مقدار مسیر را داخل فیلتر همان پروژه قرار می‌دهد.
- کلیک روی `Time` مقدار زمان را داخل فیلتر همان پروژه قرار می‌دهد.
- کلیک روی `Verb` مقدار method را داخل فیلتر method همان پروژه قرار می‌دهد.
- کلیک روی `Status` مقدار status را داخل فیلتر status همان پروژه قرار می‌دهد.
- کلیک روی `Service` مقدار `service_name` را داخل فیلتر server-side پروژه قرار می‌دهد.
- inputها دکمه clear داخلی دارند.
- هنگام تایپ در inputهای project card، focus از input خارج نمی‌شود.
- فیلتر global search حذف شد چون کاربرد آن نسبت به مدل جدید واضح نبود.
- فیلترهای server-side واضح‌تر نام‌گذاری شدند.
- معنی `Project label`، `Lookback range`، `Max rows fetched`، `Request ID` و `Loki text contains` برای کاربر قابل‌فهم‌تر شد.
- تغییرات بدون public کردن سرویس و بدون باز کردن پورت جدید انجام شدند.

## جمع‌بندی رزومه‌ای قابل دفاع

در این بخش یک UI داخلی برای observability چند پروژه Laravel بهتر شد.

ادعای قابل دفاع:

```text
Improved an internal Loki-based request log viewer by adding project-based grouping, per-project client-side filters, click-to-filter interactions, clearable inputs, focus-preserving re-renders, and clearer server-side Loki filter controls while keeping the service internal-only behind an SSH tunnel.
```

این ادعا فقط به همین محدوده مربوط است و نباید به عنوان تغییر backend، تغییر security model، یا اضافه شدن authentication معرفی شود.
