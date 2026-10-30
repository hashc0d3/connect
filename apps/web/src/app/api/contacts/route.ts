import { isUnique, many, one, pool } from "@/lib/db";
import { requireUser } from "@/lib/guard";
import { guessOperator, normalizePhone } from "@/lib/phone";

export async function GET(request: Request) {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  const url = new URL(request.url);
  const q = url.searchParams.get("q") || "";
  const groupId = Number(url.searchParams.get("groupId") || 0);
  const personal = url.searchParams.get("personal") === "1";
  const page = Math.max(1, Number(url.searchParams.get("page") || 1));
  const pageSize = [20, 50, 100, 200].includes(Number(url.searchParams.get("pageSize"))) ? Number(url.searchParams.get("pageSize")) : 20;
  const filter = [user.id, personal, groupId, q];
  const where = `
    WHERE c.user_id = $1
      AND ($2::bool = false OR c.personal)
      AND ($3::int = 0 OR EXISTS (SELECT 1 FROM contact_groups cg WHERE cg.contact_id = c.id AND cg.group_id = $3))
      AND ($4 = '' OR position(lower($4) in lower(c.number)) > 0 OR position(lower($4) in lower(c.fname)) > 0 OR position(lower($4) in lower(c.lname)) > 0)
  `;
  const total = await one<{ total: number }>(`SELECT count(*)::int AS total FROM contacts c ${where}`, filter);
  const contacts = await many(
    `SELECT c.id, c.number, c.lname, c.fname, c.sname, c.operator_name, c.hlr_status
     FROM contacts c ${where} ORDER BY c.id DESC LIMIT $5 OFFSET $6`,
    [...filter, pageSize, (page - 1) * pageSize],
  );
  const groups = (await many<{ id: number; name: string; members: number }>(
    `SELECT g.id, g.name, count(cg.contact_id)::int AS members
     FROM "groups" g LEFT JOIN contact_groups cg ON cg.group_id = g.id
     WHERE g.user_id = $1 GROUP BY g.id ORDER BY g.id`,
    [user.id],
  )).map((group) => ({ id: group.id, name: group.name, _count: { members: group.members } }));
  const counts = await one<{ all: number; own: number }>(
    `SELECT count(*)::int AS "all", count(*) FILTER (WHERE personal)::int AS own FROM contacts WHERE user_id = $1`,
    [user.id],
  );
  return Response.json({ total: total?.total || 0, page, pageSize, all: counts?.all || 0, own: counts?.own || 0, groups, contacts });
}

export async function POST(request: Request) {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  const body = await request.json();
  const number = normalizePhone(String(body.number || ""));
  if (!number) return Response.json({ message: "Номер: только цифры, код страны, без +" }, { status: 422 });
  const operator = guessOperator(number);
  const client = await pool.connect();
  try {
    await client.query("BEGIN");
    const created = await client.query(
      `INSERT INTO contacts (user_id, number, lname, fname, sname, sex, param1, param2, param3, operator_code, operator_name, personal)
       VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id, number, lname, fname, sname`,
      [user.id, number, String(body.lname || ""), String(body.fname || ""), String(body.sname || ""),
        body.sex === "male" || body.sex === "female" ? body.sex : "unknown",
        String(body.param1 || ""), String(body.param2 || ""), String(body.param3 || ""),
        operator.code, operator.name, !body.groupId],
    );
    if (body.groupId) {
      await client.query(
        `INSERT INTO contact_groups (contact_id, group_id) SELECT $1, id FROM "groups" WHERE id = $2 AND user_id = $3`,
        [created.rows[0].id, Number(body.groupId), user.id],
      );
    }
    await client.query("COMMIT");
    return Response.json(created.rows[0]);
  } catch (error) {
    await client.query("ROLLBACK");
    if (isUnique(error)) return Response.json({ message: "Такой номер уже есть" }, { status: 409 });
    throw error;
  } finally {
    client.release();
  }
}

export async function DELETE(request: Request) {
  const user = await requireUser();
  if (!user) return Response.json({ message: "Нужна авторизация" }, { status: 401 });
  const ids = (new URL(request.url).searchParams.get("ids") || "").split(",").map(Number).filter(Boolean);
  await pool.query(`DELETE FROM contacts WHERE user_id = $1 AND id = ANY($2::int[])`, [user.id, ids]);
  return Response.json({ ok: true });
}
