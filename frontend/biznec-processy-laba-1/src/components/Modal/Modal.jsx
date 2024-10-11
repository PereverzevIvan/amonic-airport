import "./Modal.scss";
import { useEffect } from "react";

function Modal({
  children,
  show = false,
  handleClose,
  title = "Модальное окно",
}) {
  const showHideClassName = show ? "modal modal_showed" : "modal modal_hidden";

  useEffect(() => {
    if (show) {
      document.body.style.overflow = "hidden"; // Отключаем прокрутку
    } else {
      document.body.style.overflow = "auto"; // Восстанавливаем прокрутку
    }

    // Очистка после закрытия модального окна
    return () => {
      document.body.style.overflow = "auto";
    };
  }, [show]);

  return (
    <>
      <div className={showHideClassName}>
        <div className="modal__window">
          <div className="modal__header">
            <p className="common-text common-text_big modal_title">{title}</p>
            <button className="close-button" onClick={handleClose}>
              <svg
                className="cross-svg"
                fill="#000000"
                width="800px"
                height="800px"
                viewBox="0 0 16 16"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M0 14.545L1.455 16 8 9.455 14.545 16 16 14.545 9.455 8 16 1.455 14.545 0 8 6.545 1.455 0 0 1.455 6.545 8z"
                  fillRule="evenodd"
                />
              </svg>
            </button>
          </div>
          <div className="modal__body">{children}</div>
        </div>
      </div>
    </>
  );
}

export default Modal;
