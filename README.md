# wvu-flights ✈️


* a SQLite(3) database of flight records from [West Virginia University](https://wvu.edu)
    * download the current [flights.db (307 MB)](https://l.abs.codes/data/wvu-data/flights.db)
* a web-accessible way to browse and query the database (with rich forms and SQL queries) using the [Datsette](https://datasette.io) project
    * see https://l.abs.codes/data/wvu-data/flights
* a static website built with [Hugo](https://gohugo.io) to display trip, passenger, and invoice data in an enriched way with support for basic text searches
    * see https://wvu-flights.pages.dev
* a simple command-line tool to help project admins manage database updates from the various formats (in scanned PDFs) as used by the [WVU Office of Procurement Contracting & Payment Services](https://procurement.wvu.edu/) and the [West Virginia State Auditor's Office](https://www.wvsao.gov/).
    * see [./cmd/wvuflights](./cmd/wvuflights/)

### Questions & Comments

See the [Frequently Asked Questions](https://github.com/AustinDizzy/wvu-flights/wiki/Frequently-Asked-Questions) in the repo wiki. Public questions & comments also available via email mailing list at `~abs/wvu-flights@lists.sr.ht`.

### License
![Creative Commons Zero v1.0](https://licensebuttons.net/p/zero/1.0/88x15.png)

This project is released into the public domain via Creative Commons Zero. See the [LICENSE](./LICENSE) file.

## Data

All information has been sourced from legal requests made under W.Va. Code § 29B-1-1 (WVFOIA)), and the intent is to keep the database updated on a rolling basis as information is released.

Current data includes 678 flights spanning from **Oct 5, 2016** to **Jul 26, 2023**, with over $7 million dollars worth of flight costs. See [./data/schema.sql](./data/schema.sql) for the database schema.

<table>
<tr><th>trips by month</th></tr>
<tr><td>

| months  | total_cost | num_trips |
|---------|------------|-----------|
| 2023-07 | $43,296    |         5 |
| 2023-06 | $157,875   |        10 |
| 2023-05 | $112,671   |        13 |
| 2023-04 | $158,926   |         9 |
| 2023-03 | $193,904   |        16 |
| 2023-02 | $149,162   |        10 |
| 2023-01 | $251,430   |        15 |
| 2022-12 | $149,516   |        12 |
| 2022-11 | $80,095    |         6 |
| 2022-10 | $101,330   |         9 |
| 2022-09 | $113,805   |        10 |
| 2022-08 | $31,102    |         4 |
| 2022-07 | $140,926   |         7 |
| 2022-06 | $149,525   |        11 |
| 2022-05 | $48,348    |         7 |
| 2022-04 | $96,824    |         7 |
| 2022-03 | $159,975   |        10 |
| 2022-02 | $56,071    |         7 |
| 2022-01 | $48,928    |         6 |
| 2021-11 | $102,550   |         7 |
| 2021-10 | $71,259    |         4 |
| 2021-09 | $37,098    |         5 |
| 2021-08 | $13,337    |         2 |
| 2021-07 | $73,723    |         7 |
| 2021-06 | $77,205    |         6 |
| 2021-05 | $34,603    |         5 |
| 2021-04 | $112,669   |         9 |
| 2021-03 | $99,294    |         5 |
| 2021-02 | $61,720    |         4 |
| 2021-01 | $12,515    |         2 |
| 2020-12 | $99,281    |         6 |
| 2020-11 | $45,761    |         2 |
| 2020-10 | $12,032    |         2 |
| 2020-09 | $89,854    |         8 |
| 2020-08 | $95,561    |         7 |
| 2020-07 | $17,268    |         3 |
| 2020-06 | $11,915    |         2 |
| 2020-05 | $11,640    |         2 |
| 2020-04 | $17,775    |         3 |
| 2020-03 | $36,008    |         6 |
| 2020-02 | $160,305   |        10 |
| 2020-01 | $53,130    |         4 |
| 2019-12 | $45,658    |         6 |
| 2019-11 | $111,697   |         9 |
| 2019-10 | $143,333   |        12 |
| 2019-09 | $71,160    |         7 |
| 2019-08 | $75,336    |        11 |
| 2019-07 | $96,139    |        13 |
| 2019-06 | $111,943   |         8 |
| 2019-05 | $18,246    |         2 |
| 2019-04 | $96,446    |         7 |
| 2019-03 | $162,191   |        12 |
| 2019-02 | $182,948   |         9 |
| 2019-01 | $100,673   |        10 |
| 2018-12 | $55,907    |         5 |
| 2018-11 | $70,855    |         7 |
| 2018-10 | $104,893   |        12 |
| 2018-09 | $69,924    |         9 |
| 2018-08 | $66,321    |        10 |
| 2018-07 | $87,891    |         8 |
| 2018-06 | $160,162   |        15 |
| 2018-05 | $79,242    |         8 |
| 2018-04 | $89,480    |        14 |
| 2018-03 | $168,281   |        14 |
| 2018-02 | $128,230   |        12 |
| 2017-11 | $125,678   |        12 |
| 2017-06 | $113,667   |        13 |
| 2017-05 | $64,479    |        10 |
| 2017-04 | $121,661   |        13 |
| 2017-03 | $85,477    |        15 |
| 2017-02 | $150,024   |        10 |
| 2017-01 | $54,886    |         7 |
| 2016-12 | $144,401   |        13 |
| 2016-11 | $60,656    |         8 |
| 2016-10 | $123,521   |        12 |
| 2016-09 | $92,506    |        15 |
| 2016-08 | $56,415    |         9 |
| 2016-07 | $64,086    |         8 |
| 2016-06 | $109,390   |        15 |
| 2016-05 | $121,041   |        13 |
| **Total**   | **$7,375,090** |       **678** |

<details> 
  <summary>View SQL Query</summary>

   ```sql
SELECT
    months,
    total_cost,
    num_trips
FROM
(
    SELECT
        strftime('%Y-%m',
                   CASE
                       WHEN instr(trips.date, ';') > 0
                       THEN substr(trips.date, instr(trips.date, ';') + 1)
                       ELSE trips.date
                   END
        ) AS months,
        PRINTF("$%,2d", SUM(fuel + landing + crew_expense + dom_tax + billing_amount)) AS total_cost,
        COUNT(*) AS num_trips,
        1 AS sort_order
    FROM trips
    GROUP BY months

    UNION

    SELECT
        'Total' AS months,
        PRINTF("$%,2d", SUM(fuel + landing + crew_expense + dom_tax + billing_amount)) AS total_cost,
        COUNT(*) AS num_trips,
        2 AS sort_order
    FROM trips
) AS combined
ORDER BY sort_order, months DESC;
   ```
</details>
</td></tr>
</table>