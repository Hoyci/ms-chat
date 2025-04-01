import { create } from "zustand";

type LayoutKeys = "newChat" | "rooms";

type LayoutState = {
  currentLayout: LayoutKeys;
  setLayout: (layout: LayoutKeys) => void;
};

const useLayoutStore = create<LayoutState>((set) => ({
  currentLayout: "rooms",
  setLayout: (layout: LayoutKeys) => set({ currentLayout: layout }),
}));

export default useLayoutStore;
