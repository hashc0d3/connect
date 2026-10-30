import { many, one, pool } from "@/lib/db";
import { requireUser } from "@/lib/guard";

export async function GET() {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  return Response.json(await many(`SELECT id, name, body FROM templates WHERE user_id = $1 ORDER BY id DESC`, [user.id]));
}

export async function POST(request: Request) {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  const body = await request.json();
  if (body.id && body.remove) {
    await pool.query(`DELETE FROM templates WHERE id = $1 AND user_id = $2`, [Number(body.id), user.id]);
    return Response.json({ ok: true });
  }
  const template = await one(
    `INSERT INTO templates (user_id, name, body) VALUES ($1,$2,$3) RETURNING id, name, body`,
    [user.id, String(body.name || "Шаблон"), String(body.body || "")],
  );
  return Response.json(template);
}
