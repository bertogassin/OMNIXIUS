-- B2 Embedded finance: installment plan for orders (stub; full flow via Trade later).
ALTER TABLE orders ADD COLUMN installment_plan TEXT DEFAULT '';
