import { isUnique, many, one, pool } from "@/lib/db";
import { requireUser } from "@/lib/guard";
import { normalizePhone } from "@/lib/phone";

export async function GET() {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  return Response.json(await many(
    `SELECT id, number, comment FROM blacklist_entries WHERE user_id = $1 ORDER BY id DESC`,
    [user.id],
  ));
}

export async function POST(request: Request) {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  const body = await request.json();
  const number = normalizePhone(String(body.number || ""));
  if (!number) return Response.json({ message: "Некорректный номер" }, { status: 422 });
  try {
    const entry = await one(
      `INSERT INTO blacklist_entries (user_id, number, comment) VALUES ($1,$2,$3) RETURNING id, number, comment`,
      [user.id, number, String(body.comment || "")],
    );
    return Response.json(entry);
  } catch (error) {
    if (isUnique(error)) return Response.json({ message: "Номер уже в списке" }, { status: 409 });
    throw error;
  }
}

export async function DELETE(request: Request) {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  const id = Number(new URL(request.url).searchParams.get("id"));
  await pool.query(`DELETE FROM blacklist_entries WHERE id = $1 AND user_id = $2`, [id, user.id]);
  return Response.json({ ok: true });
}
