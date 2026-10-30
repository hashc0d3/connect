import { many, one } from "@/lib/db";
import { requireUser } from "@/lib/guard";

export async function GET() {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  return Response.json(await many(
    `SELECT id, title, type, status FROM billing_profiles WHERE user_id = $1 ORDER BY id DESC`,
    [user.id],
  ));
}

export async function POST(request: Request) {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  const body = await request.json();
  const profile = await one(
    `INSERT INTO billing_profiles (user_id, type, title, inn, kpp, ogrn, address, email, phone)
     VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id, title, type, status`,
    [
      user.id,
      ["legal", "entrepreneur", "person"].includes(body.type) ? body.type : "legal",
      String(body.title || ""),
      String(body.inn || ""),
      String(body.kpp || ""),
      String(body.ogrn || ""),
      String(body.address || ""),
      String(body.email || user.email),
      String(body.phone || ""),
    ],
  );
  return Response.json(profile);
}
