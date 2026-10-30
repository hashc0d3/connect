import { many, pool } from "@/lib/db";
import { requireUser } from "@/lib/guard";
import { guessOperator } from "@/lib/phone";

export async function POST(request: Request) {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  const body = await request.json();
  const ids: number[] = Array.isArray(body.ids) ? body.ids.map(Number) : [];
  const contacts = await many<{ id: number; number: string }>(
    `SELECT id, number FROM contacts WHERE user_id = $1 AND id = ANY($2::int[])`,
    [user.id, ids],
  );
  for (const contact of contacts) {
    const operator = guessOperator(contact.number);
    await pool.query(
      `UPDATE contacts SET operator_code = $1, operator_name = $2, hlr_status = $3 WHERE id = $4 AND user_id = $5`,
      [operator.code, operator.name, operator.name === "OTHER" ? 0 : 1, contact.id, user.id],
    );
  }
  return Response.json({ checked: contacts.length });
}
