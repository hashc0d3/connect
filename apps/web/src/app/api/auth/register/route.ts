import { createHash, randomBytes } from "crypto";
import bcrypt from "bcryptjs";
import { isUnique, pool } from "@/lib/db";
import { writeSession } from "@/lib/session";

export async function POST(request: Request) {
  const body = await request.json();
  const email = String(body.email || "").trim().toLowerCase();
  const password = String(body.password || "");
  const name = String(body.name || "").trim();
  if (!email || password.length < 8 || !name) {
    return Response.json({ message: "Заполните имя, email и пароль от 8 символов" }, { status: 422 });
  }
  const rawKey = `cnt_${randomBytes(12).toString("hex")}`;
  const client = await pool.connect();
  try {
    await client.query("BEGIN");
    const created = await client.query<{ id: string; role: string }>(
      `INSERT INTO users (email, name, password_hash) VALUES ($1,$2,$3) RETURNING id::text, role`,
      [email, name, await bcrypt.hash(password, 10)],
    );
    const user = created.rows[0];
    await client.query(
      `INSERT INTO api_keys (user_id, name, prefix, key_hash) VALUES ($1,$2,$3,$4)`,
      [user.id, "Основной", rawKey.slice(0, 12), createHash("sha256").update(rawKey).digest("hex")],
    );
    await client.query("COMMIT");
    await writeSession({ uid: user.id, role: user.role, exp: Math.floor(Date.now() / 1000) + 60 * 60 * 24 * 7 });
    return Response.json({ apiKey: rawKey });
  } catch (error) {
    await client.query("ROLLBACK");
    if (isUnique(error)) return Response.json({ message: "Такой email уже зарегистрирован" }, { status: 409 });
    throw error;
  } finally {
    client.release();
  }
}
