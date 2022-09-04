import { atom, atomFamily, selectorFamily } from "recoil";
import { Profile } from "./profiles";
import { recoilPersist } from "recoil-persist";
import { fetchProxy } from "../../routes/instances/InstanceAPI";
import { Service } from "./services";

const { persistAtom } = recoilPersist();

export type Instance = {
  id: number;
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
  id: 0,
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

export type InstanceProxyService = {
  id: string;
  port: number;
  status: string;
  service: Service;
};

export const instanceProxySelector = selectorFamily({
  key: "getProxyAPI",
  get:
    (id: number) =>
    ({ get }) => {
      // Get the instance we are looking for
      const instance = get(instances(id));

      // Return the proxy included in the instance from the API
      return fetchProxy(instance);
    },
});
