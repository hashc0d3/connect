export default function DocsPage() {
  return (
    <main className="page doc">
      <a className="brand" href="/">Connect<i>→</i></a>
      <h1>API</h1>
      <p>Базовый URL REST: <code>http://localhost:8080</code>. gRPC: <code>localhost:50051</code>, пакет <code>connect.v1.SmsGateway</code>.</p>
      <p>Ключ передаётся заголовком <code>Authorization: Bearer &lt;ключ&gt;</code>. Номер — цифры с кодом страны, без «+», не короче 10 символов. Текст для EasySMS — не короче 4 символов, получателей не больше 1000.</p>
      <pre className="code">{`POST /v2/sms/send
{
  "receivers": [{ "phone": "79990000000" }],
  "text": "Test Sms",
  "sign": "CONNECT"
}`}</pre>
      <pre className="code">{`grpcurl -plaintext -H "authorization: Bearer <ключ>" \\
  -d '{"receivers":[{"phone":"79990000000"}],"text":"Test Sms","sign":"CONNECT"}' \\
  localhost:50051 connect.v1.SmsGateway/SendSms`}</pre>
      <p>Перед вызовом easysms/v3 Connect проверяет подпись, чёрный список, баланс, наценку и договор. Одобренная подпись уходит поставщику как originator. Статусы доставки приходят вебхуком. Если логин поставщика не задан, отправка помечается успешной локально.</p>
      <p>Также доступны <code>/v2/sms/status</code>, <code>/v2/sms/list</code>, <code>/v2/balance</code>, <code>/v2/tariffs</code>, <code>/v2/otp/create</code>, <code>/v2/otp/validate</code>, контакты, группы и подписи.</p>
    </main>
  );
}
