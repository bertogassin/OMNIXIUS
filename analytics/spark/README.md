# OMNIXIUS — Spark analytics

Analytics & recommendations (Scala). Run with Spark 3.x:

```bash
sbt package
spark-submit --class omnixius.RecommendationsJob target/scala-2.12/omnixius-analytics_2.12-0.1.0.jar
```
