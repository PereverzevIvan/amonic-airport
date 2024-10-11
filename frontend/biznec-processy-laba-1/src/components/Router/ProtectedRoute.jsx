import { Route } from "react-router-dom";
import { useAuth } from "../../context/authContext";
import { Navigate } from "react-router-dom";
import * as rp from "../../global/routerPaths/routerPaths.jsx";

export function ProtectedRoute({ children, requiredRole }) {
  const { isAuth, role } = useAuth();

  if (!isAuth) {
    return <Navigate to={rp.LOGIN_ROUTE} />;
  }

  if (role != "admin" && role != requiredRole) {
    return <Navigate to={rp.FORBIDDEN_ERROR_ROUTE} />;
  }

  return children;
}
