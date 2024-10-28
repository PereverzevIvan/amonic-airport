import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import App from "./App.jsx";
import { ToastProvider } from "./context/ToastContext.jsx";
import { ApiProvider } from "./context/apiContext";
import { AuthProvider } from "./context/authContext";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const queryClient = new QueryClient();

createRoot(document.getElementById("root")).render(
  //<StrictMode>
  <QueryClientProvider client={queryClient}>
    <BrowserRouter>
      <AuthProvider>
        <ApiProvider>
          <ToastProvider>
            <App />
          </ToastProvider>
        </ApiProvider>
      </AuthProvider>
    </BrowserRouter>
  </QueryClientProvider>,
  //</StrictMode>,
);
