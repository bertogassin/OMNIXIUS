@echo off
REM Delete SQLite DB so backend creates a fresh one on next start. All users/registrations removed.
if exist "db\omnixius.db" (
  del "db\omnixius.db"
  echo Database deleted. Start backend again (go run . or start-backend.bat) to get a fresh DB.
) else (
  echo No db\omnixius.db found. Nothing to delete.
)
