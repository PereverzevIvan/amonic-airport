import Footer from "./components/Footer/Footer";
import Header from "./components/Header/Header";
import Router from "./components/Router/Router";
import "./scss/init.scss";

function App() {
  return (
    <div className="App">
      <Header />
      <main className="main">
        <Router />
      </main>
      <Footer />
    </div>
  );
}

export default App;
