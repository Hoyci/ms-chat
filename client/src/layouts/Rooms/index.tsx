import Header from "./Header";
import List from "./List";

function Rooms() {
  return (
    <div className="flex-shrink-0 flex-grow-0 basis-[30%] bg-neutral-700 text-neutral-100 px-4 flex flex-col h-full">
      <Header />
      <List />
    </div>
  );
}

export default Rooms;
