import { createHmac, timingSafeEqual } from "crypto";
import { cookies } from "next/headers";

export type Session = { uid: string; role: string; exp: number };

const cookieName = "connect_session";

function secret() {
  return process.env.SESSION_SECRET || "dev-session-secret-connect";
}

export function seal(session: Session) {
  const body = Buffer.from(JSON.stringify(session)).toString("base64url");
  const sig = createHmac("sha256", secret()).update(body).digest("base64url");
  return `${body}.${sig}`;
}

export function open(token?: string): Session | null {
  if (!token) return null;
  const [body, sig] = token.split(".");
  if (!body || !sig) return null;
  const expected = createHmac("sha256", secret()).update(body).digest("base64url");
  const a = Buffer.from(sig);
  const b = Buffer.from(expected);
  if (a.length !== b.length || !timingSafeEqual(a, b)) return null;
  const session = JSON.parse(Buffer.from(body, "base64url").toString()) as Session;
  if (session.exp < Date.now() / 1000) return null;
  return session;
}

export async function writeSession(session: Session) {
  const jar = await cookies();
  jar.set(cookieName, seal(session), {
    httpOnly: true,
    sameSite: "lax",
    path: "/",
    maxAge: 60 * 60 * 24 * 7,
  });
}

export async function clearSession() {
  const jar = await cookies();
  jar.delete(cookieName);
}

export async function readSession() {
  const jar = await cookies();
  return open(jar.get(cookieName)?.value);
}
