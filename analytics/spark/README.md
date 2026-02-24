# OMNIXIUS — Spark analytics

Аналитика и рекомендации (Scala/Java). По ARCHITECTURE.md — серверное направление. **Уровень:** заготовка (один джоб); дальше — реальные пайплайны.

Run with Spark 3.x:

```bash
sbt package
spark-submit --class omnixius.RecommendationsJob target/scala-2.12/omnixius-analytics_2.12-0.1.0.jar
```
