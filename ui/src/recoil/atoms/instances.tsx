import {
  atom,
  selectorFamily,
  useRecoilCallback,
  useRecoilValue,
} from "recoil";
import { Profile } from "./profiles";
import { recoilPersist } from "recoil-persist";
import { fetchProxy } from "../../routes/instances/InstanceAPI";
import { DefaultService, Service } from "./services";
import { pseudoRandomBytes } from "crypto";

const { persistAtom } = recoilPersist();

export type Instance = {
  id: number;
  name: string;
  host: string;
  port: number;
  description: string;
  profile: Profile | undefined;
};

export type InstanceProxy = {
  id: string;
  port: number;
  status: string;
  service: Service;
};

export const DefaultInstance = {
  id: 0,
  name: "Default",
  host: "localhost",
  port: 2022,
  description: "Default instance",
  profile: undefined,
} as Instance;

export const DefaultInstanceProxy = {
  id: "",
  port: 0,
  status: "stopped",
  service: DefaultService,
};

// Default instance that will be used when navigating to "home/instance" path
export const instance = atom<Instance>({
  key: "instance",
  default: DefaultInstance,
  effects_UNSTABLE: [persistAtom],
});

// Default instance to be used as dummy to pre-populate forms
export const instanceFormFields = atom<Instance>({
  key: "instanceFormFields",
  default: DefaultInstance,
});

// Returns the address of the instance in the form of `<host>:<port>`
// Example: localhost:2022
export const GetInstanceAddress = () => {
  const inst = useRecoilValue(instance);
  return inst.host + ":" + inst.port;
};

// Callback selector that returns two functions. One to register and the other to remove
// a service from the Instance
export const useInstanceProxy = () => {
  const pxs = useRecoilValue(proxies);

  // Callback to `register` a service in the Instance
  const registerProxy = useRecoilCallback(({ set }) => (px: InstanceProxy) => {
    const prev = pxs.find((p) => p.id === px.id);

    if (!prev) {
      set(proxies, [...pxs, px]);
    }
  });

  // Callback to `remove` a service from the Instance
  const removeProxy = useRecoilCallback(({ set }) => (px: InstanceProxy) => {
    set(proxies, (prev) => prev.filter((x) => x.id !== px.id));
  });

  const removeProxyFromService = useRecoilCallback(
    ({ set }) =>
      (id: string) => {
        set(proxies, (prev) => prev.filter((p) => p.service.id !== id));
      }
  );

  return {
    registerProxy,
    removeProxy,
    removeProxyFromService,
  };
};

const remoteProxiesEffect =
  () =>
  ({ setSelf }: { setSelf: any }) => {
    setSelf(async () => {
      const response = await fetchProxy("localhost:2022");

      if (response.error) {
        throw response.error;
      }

      return response.map((proxy: any) => {
        return { ...DefaultInstanceProxy, ...proxy };
      });
    });
  };

export const proxies = atom<InstanceProxy[]>({
  key: "proxies",
  default: [],
  effects: [remoteProxiesEffect()],
});

export const instanceProxySelector = selectorFamily({
  key: "proxy/default",
  get:
    (id: string) =>
    ({ get }) => {
      const pxs = get(proxies);
      return pxs.find((x) => x.id === id);
    },
  set:
    (id: string) =>
    ({ get, set }, newValue) => {
      const pxs: any = get(proxies);
      const pxInd = pxs.findIndex((x: InstanceProxy) => x.id === id);
      let cp = [...pxs];
      cp[pxInd] = { ...cp[pxInd], newValue };
      return set(proxies, cp);
    },
});

//
// FORMS
//
type FormError = {
  [key: string]: string;
};

export const DefaultFormError = {} as FormError;

export const intanceFormErrors = atom<FormError>({
  key: "intanceFormErrors",
  default: DefaultFormError,
});

export const intanceFormFieldErrors = selectorFamily({
  key: "profileFormFieldErrors",
  get:
    (field: string) =>
    ({ get }) =>
      get(intanceFormErrors)[field],
});
