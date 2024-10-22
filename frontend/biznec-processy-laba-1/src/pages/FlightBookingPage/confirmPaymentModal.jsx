import { useEffect, useState } from "react";
import Modal from "../../components/Modal/Modal";
import Button from "../../components/Button/Button";
import { confirmTickets } from "../../api/tickets";
import { useApi } from "../../context/apiContext";
import { useToast } from "../../context/ToastContext";

export function ConfirmPaymentModal({
  show,
  handleClose,
  outboundFlight,
  inboundFlight,
  tickets = [],
  sub = false,
  closeParentModal,
}) {
  const { apiClient } = useApi();
  const { addToast } = useToast();

  const [sum, setSum] = useState(0);

  useEffect(() => {
    if (!show) {
      setSum(0);
    } else {
      if (outboundFlight) {
        if (inboundFlight) {
          setSum(outboundFlight.price + inboundFlight.price);
        } else {
          setSum(outboundFlight.price);
        }
      }
    }
  }, [show]);

  function handleSubmit() {
    let requestBody = {
      tickets: tickets.map((ticket) => ticket.id) || null,
    };
    confirmTickets(apiClient, requestBody)
      .then((response) => {
        console.log(response);
        if (response.status == 200) {
          addToast("Билеты подтверждены");
          handleClose();
          closeParentModal();
        }
      })
      .catch((error) => {
        console.log(error);
        if (error.response) {
          addToast(error.response.data, "error");
        }
      });
  }

  return (
    <>
      <Modal
        sub={sub}
        show={show}
        handleClose={handleClose}
        title="Подтверждение оплаты"
      >
        <form className="form">
          <fieldset className="form__fieldset tickets-confirm-form">
            <legend className="form__legend">Подтвердите оплату</legend>
            <label className="form__label">Общая сумма: </label>
            <p className="common-text">${sum}</p>

            <label className="form__label">Способ оплаты: </label>
            <div
              style={{ display: "flex", flexDirection: "column", gap: "5px" }}
            >
              <label>
                <input
                  type="radio"
                  name="payment-type"
                  value="Наичные"
                  onChange={() => {}}
                  checked
                />
                Наличные
              </label>
              <label>
                <input
                  type="radio"
                  name="payment-type"
                  value="Кредитка"
                  onChange={() => {}}
                />
                Кредитка
              </label>
              <label>
                <input
                  type="radio"
                  name="payment-type"
                  value="Воунчер"
                  onChange={() => {}}
                />
                Воунчер
              </label>
            </div>
          </fieldset>
          <Button color="green" onClick={handleSubmit}>
            Подтвердить билеты
          </Button>
        </form>
      </Modal>
    </>
  );
}
