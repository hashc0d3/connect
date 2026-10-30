import { isUnique, many, one } from "@/lib/db";
import { requireUser } from "@/lib/guard";

export async function GET() {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  return Response.json(await many(`SELECT id, name, channel, status, extend_status FROM signs WHERE user_id = $1 ORDER BY id`, [user.id]));
}

export async function POST(request: Request) {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  const body = await request.json();
  const name = String(body.name || "").trim();
  if (!name) return Response.json({ message: "Укажите подпись" }, { status: 422 });
  try {
    const sign = await one(
      `INSERT INTO signs (user_id, name, channel) VALUES ($1,$2,$3) RETURNING id, name, channel, status, extend_status`,
      [user.id, name, body.channel || "PAY_SIGN"],
    );
    return Response.json(sign);
  } catch (error) {
    if (isUnique(error)) return Response.json({ message: "Подпись уже есть" }, { status: 409 });
    throw error;
  }
}
