import { pool } from "@/lib/db";
import { requireUser } from "@/lib/guard";

export async function POST(request: Request) {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  const body = await request.json();
  const text = String(body.text || "").trim();
  if (text.length < 3) return Response.json({ message: "Слишком короткое сообщение" }, { status: 422 });
  await pool.query(`INSERT INTO feedback (user_id, text) VALUES ($1,$2)`, [user.id, text]);
  return Response.json({ ok: true });
}
