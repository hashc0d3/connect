import { one } from "@/lib/db";
import { requireUser } from "@/lib/guard";

export async function GET() {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  return Response.json(await one(
    `SELECT enabled, template_id, sign_name FROM auto_rules WHERE user_id = $1`,
    [user.id],
  ));
}

export async function POST(request: Request) {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  const body = await request.json();
  const rule = await one(
    `INSERT INTO auto_rules (user_id, enabled, template_id, sign_name)
     VALUES ($1,$2,$3,$4)
     ON CONFLICT (user_id) DO UPDATE
       SET enabled = EXCLUDED.enabled, template_id = EXCLUDED.template_id, sign_name = EXCLUDED.sign_name
     RETURNING enabled, template_id, sign_name`,
    [user.id, Boolean(body.enabled), body.templateId ? Number(body.templateId) : null, String(body.signName || "")],
  );
  return Response.json(rule);
}
