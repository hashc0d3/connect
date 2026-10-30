import { one } from "./db";
import { readSession } from "./session";

export type SessionUser = {
  id: string;
  email: string;
  name: string;
  role: string;
  balance: string;
};

export async function currentUser() {
  const session = await readSession();
  if (!session) return null;
  return one<SessionUser>(
    `SELECT id::text, email, name, role, balance::text AS balance FROM users WHERE id = $1`,
    [session.uid],
  );
}

export async function requireUser() {
  return currentUser();
}
