import { atom, atomFamily, selectorFamily } from "recoil";
import { Profile } from "./profiles";
import { recoilPersist } from "recoil-persist";

const { persistAtom } = recoilPersist();

export type Instance = {
  id: number | undefined;
  name: string;
  host: string;
  description: string;
  profile: Profile | undefined;
};

export const instanceIds = atom<number[]>({
  key: "instanceIds",
  default: [],
  effects_UNSTABLE: [persistAtom],
});

export const DefaultInstance = {
  id: undefined,
  name: "",
  host: "",
  description: "",
  profile: undefined,
};

export const instances = atomFamily<Instance, number>({
  key: "instance",
  default: DefaultInstance,
  effects_UNSTABLE: [persistAtom],
});

export const intanceFormErrors = atom({
  key: "intanceFormErrors",
  default: {} as { [key: string]: string },
});

export const intanceFormFieldErrors = selectorFamily({
  key: "profileFormFieldErrors",
  get:
    (field: string) =>
    ({ get }) =>
      get(intanceFormErrors)[field],
});

export const instanceFormFields = atom({
  key: "instanceFormFields",
  default: DefaultInstance as { [key: string]: any },
});
