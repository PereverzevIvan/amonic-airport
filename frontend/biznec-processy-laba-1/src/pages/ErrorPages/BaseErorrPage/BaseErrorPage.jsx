import "./BaseErrorPage.scss";

export function BaseErrorPage({ errorCode = 1, message = "Ошибка" }) {
  return (
    <>
      <section className="base-error-page">
        <h1 className="title title_big">{errorCode}</h1>
        <p className="common-text common-text_big">{message}</p>
      </section>
    </>
  );
}
