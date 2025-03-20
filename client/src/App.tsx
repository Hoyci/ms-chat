import Header from "./layouts/Header";

function App() {
  return (
    <div className="h-screen w-full bg-background flex items-center justify-center">
      <div className="top-[19px] w-[calc(100%-38px)] max-w-[1600px] h-[calc(100%-38px)]">
        <Header />
      </div>
    </div>
  );
}

export default App;
