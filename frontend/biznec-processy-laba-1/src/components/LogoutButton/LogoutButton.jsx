import { useAuth } from "../../context/authContext";
import { useNavigate } from "react-router-dom";
import * as rp from "../../global/routerPaths/routerPaths.jsx";
import Button from "../Button/Button.jsx";
import { useApi } from "../../context/apiContext.jsx";

export function LogoutButton() {
  const { isAuth, logout } = useAuth();
  const { apiClient } = useApi();

  const navigate = useNavigate();

  async function onClickLogout() {
    if (isAuth) {
      await logout(apiClient);
    }
    navigate(rp.LOGIN_ROUTE);
  }

  return (
    <Button color="red" onClick={onClickLogout}>
      Выйти
    </Button>
  );
}
