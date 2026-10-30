import { createHash, randomBytes } from "crypto";
import { many, pool } from "@/lib/db";
import { requireUser } from "@/lib/guard";

export async function GET() {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  return Response.json(await many(
    `SELECT id, name, prefix, created_at, last_used_at FROM api_keys WHERE user_id = $1 ORDER BY id DESC`,
    [user.id],
  ));
}

export async function POST(request: Request) {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  const body = await request.json();
  const raw = `cnt_${randomBytes(12).toString("hex")}`;
  await pool.query(
    `INSERT INTO api_keys (user_id, name, prefix, key_hash) VALUES ($1,$2,$3,$4)`,
    [user.id, String(body.name || "Ключ"), raw.slice(0, 12), createHash("sha256").update(raw).digest("hex")],
  );
  return Response.json({ apiKey: raw });
}

export async function DELETE(request: Request) {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  const id = Number(new URL(request.url).searchParams.get("id"));
  await pool.query(`DELETE FROM api_keys WHERE id = $1 AND user_id = $2`, [id, user.id]);
  return Response.json({ ok: true });
}
