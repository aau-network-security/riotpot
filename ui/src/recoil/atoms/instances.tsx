import { atom } from "recoil";

export const instances = atom({
  key: "instancesList",
  default: [],
});
