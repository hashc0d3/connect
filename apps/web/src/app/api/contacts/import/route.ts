import { isUnique, pool } from "@/lib/db";
import { requireUser } from "@/lib/guard";
import { guessOperator, normalizePhone } from "@/lib/phone";

export async function POST(request: Request) {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  const body = await request.json();
  const rows = Array.isArray(body.contacts) ? body.contacts : [];
  let added = 0;
  let skipped = 0;
  for (const row of rows) {
    const number = normalizePhone(String(row.number || ""));
    if (!number) {
      skipped += 1;
      continue;
    }
    const operator = guessOperator(number);
    try {
      await pool.query(
        `INSERT INTO contacts (user_id, number, operator_code, operator_name, lname, fname, sname, param1, param2, param3, personal)
         VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,true)`,
        [user.id, number, operator.code, operator.name, String(row.lname || ""), String(row.fname || ""), String(row.sname || ""), String(row.param1 || ""), String(row.param2 || ""), String(row.param3 || "")],
      );
      added += 1;
    } catch (error) {
      if (!isUnique(error)) throw error;
      skipped += 1;
    }
  }
  return Response.json({ added, skipped });
}
