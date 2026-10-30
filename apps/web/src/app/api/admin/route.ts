import { gateway } from "@/lib/gateway";
import { many, pool } from "@/lib/db";
import { requireUser } from "@/lib/guard";

async function admin() {
  const user = await requireUser();
  if (!user || user.role !== "admin") return null;
  return user;
}

export async function GET() {
  const user = await admin();
  if (!user) return Response.json({ message: "Нет доступа" }, { status: 403 });
  const [messages, signs, profiles, settings] = await Promise.all([
    many(`SELECT id, text, number FROM messages WHERE status = 8 ORDER BY id DESC LIMIT 50`),
    many(`SELECT id, name FROM signs WHERE status = 0 ORDER BY id DESC LIMIT 50`),
    many(`SELECT id, title FROM billing_profiles WHERE status = 'moderation' ORDER BY id DESC LIMIT 50`),
    many(`SELECT key, value FROM platform_settings ORDER BY key`),
  ]);
  return Response.json({ messages, signs, profiles, settings });
}

export async function POST(request: Request) {
  const user = await admin();
  if (!user) return Response.json({ message: "Нет доступа" }, { status: 403 });
  const body = await request.json();
  if (body.kind === "sign") {
    await pool.query(
      body.approve
        ? `UPDATE signs SET status = 1, extend_status = 'active' WHERE id = $1`
        : `UPDATE signs SET status = 2, extend_status = 'reject', reject_reason = $2 WHERE id = $1`,
      body.approve ? [Number(body.id)] : [Number(body.id), String(body.reason || "")],
    );
  }
  if (body.kind === "profile") {
    await pool.query(
      `UPDATE billing_profiles SET status = $2, reject_reason = $3 WHERE id = $1`,
      [Number(body.id), body.approve ? "approved" : "reject", String(body.reason || "")],
    );
  }
  if (body.kind === "message") {
    if (!body.approve) {
      await pool.query(`UPDATE messages SET status = 6, extend_status = 'reject', upstream_state = 'skipped' WHERE id = $1`, [Number(body.id)]);
    } else {
      await gateway(`/internal/messages/${Number(body.id)}/release`, user.id, {});
    }
  }
  if (body.kind === "setting") {
    await pool.query(
      `INSERT INTO platform_settings (key, value) VALUES ($1,$2) ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value`,
      [String(body.key), String(body.value)],
    );
  }
  return Response.json({ ok: true });
}
