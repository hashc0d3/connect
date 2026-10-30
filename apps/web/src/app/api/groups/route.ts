import { isUnique, one, pool } from "@/lib/db";
import { requireUser } from "@/lib/guard";

export async function POST(request: Request) {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  const body = await request.json();
  const name = String(body.name || "").trim();
  if (!name) return Response.json({ message: "Укажите имя группы" }, { status: 422 });
  try {
    const group = await one(`INSERT INTO "groups" (user_id, name) VALUES ($1,$2) RETURNING id, name`, [user.id, name]);
    return Response.json(group);
  } catch (error) {
    if (isUnique(error)) return Response.json({ message: "Группа уже существует" }, { status: 409 });
    throw error;
  }
}

export async function DELETE(request: Request) {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  const id = Number(new URL(request.url).searchParams.get("id"));
  await pool.query(`DELETE FROM "groups" WHERE id = $1 AND user_id = $2`, [id, user.id]);
  return Response.json({ ok: true });
}
