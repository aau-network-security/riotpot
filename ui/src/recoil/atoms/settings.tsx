import { atom } from "recoil";
import { recoilPersist } from "recoil-persist";

const { persistAtom } = recoilPersist();

export type Backend = {
  host: string;
  port: number;
};

export const DefaultBackend = {
  host: "localhost",
  port: 2022,
} as Backend;

export const backend = atom({
  key: "backend",
  default: DefaultBackend,
  effects_UNSTABLE: [persistAtom],
});
