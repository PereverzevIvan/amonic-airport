import Footer from "./components/Footer/Footer";
import Header from "./components/Header/Header";
import Router from "./components/Router/Router";
import { ApiProvider } from "./context/apiContext";
import { AuthProvider } from "./context/authContext";
import "./scss/init.scss";

function App() {
  return (
    <div className="App">
      <AuthProvider>
        <ApiProvider>
          <Header />
          <main className="main">
            <Router />
          </main>
          <Footer />
        </ApiProvider>
      </AuthProvider>
    </div>
  );
}

export default App;
