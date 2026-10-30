import bcrypt from "bcryptjs";
import { one } from "@/lib/db";
import { writeSession } from "@/lib/session";

export async function POST(request: Request) {
  const body = await request.json();
  const email = String(body.email || "").trim().toLowerCase();
  const password = String(body.password || "");
  const user = await one<{ id: string; role: string; passwordHash: string }>(
    `SELECT id::text, role, password_hash FROM users WHERE email = $1`,
    [email],
  );
  if (!user || !(await bcrypt.compare(password, user.passwordHash))) {
    return Response.json({ message: "Неверный email или пароль" }, { status: 401 });
  }
  await writeSession({ uid: user.id, role: user.role, exp: Math.floor(Date.now() / 1000) + 60 * 60 * 24 * 7 });
  return Response.json({ ok: true });
}
