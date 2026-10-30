import Link from "next/link";

export default function HomePage() {
  return (
    <main className="login">
      <a className="brand" href="/"><span>Connect</span><i>→</i></a>
      <p>SMS-платформа для рассылок, OTP и API. Сообщения проходят проверки Connect и уходят в EasySMS.</p>
      <div className="row">
        <Link className="btn green" href="/login">Войти</Link>
        <Link className="btn" href="/register">Регистрация</Link>
        <Link className="btn" href="/docs">Документация</Link>
      </div>
    </main>
  );
}
