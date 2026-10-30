import { Pool, type QueryResultRow } from "pg";

function connectionString() {
  const raw = process.env.DATABASE_URL || "";
  const url = new URL(raw);
  url.searchParams.delete("schema");
  if (!url.searchParams.has("sslmode")) url.searchParams.set("sslmode", "disable");
  return url.toString();
}

const globalForDb = globalThis as unknown as { pool?: Pool };

export const pool = globalForDb.pool ?? new Pool({ connectionString: connectionString() });

if (process.env.NODE_ENV !== "production") globalForDb.pool = pool;

export function row<T>(source: QueryResultRow): T {
  const out: Record<string, unknown> = {};
  for (const [key, value] of Object.entries(source)) {
    out[key.replace(/_([a-z0-9])/g, (_, char: string) => char.toUpperCase())] = value;
  }
  return out as T;
}

export async function many<T>(sql: string, params: unknown[] = []) {
  const result = await pool.query(sql, params);
  return result.rows.map((item) => row<T>(item));
}

export async function one<T>(sql: string, params: unknown[] = []) {
  const rows = await many<T>(sql, params);
  return rows[0] ?? null;
}

export function isUnique(error: unknown) {
  return typeof error === "object" && error !== null && "code" in error && (error as { code: string }).code === "23505";
}
