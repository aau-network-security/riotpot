import { atom, atomFamily, selector, selectorFamily } from "recoil";
import { Profile } from "./profiles";
import { recoilPersist } from "recoil-persist";
import { fetchProxy } from "../../routes/instances/InstanceAPI";
import { DefaultService, Service } from "./services";

const { persistAtom } = recoilPersist();

export type Instance = {
  id: number;
  name: string;
  host: string;
  port: number;
  description: string;
  profile: Profile | undefined;
};

export type InstanceProxyService = {
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

export const DefaultInstanceProxyService = {
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

// MULTIPLE INSTANCES HANDLERS
//

// Array of registered instances IDs
//    NOTE: Typically used together with `atomFamilies` to track IDs in the collection
export const instanceIds = atom<number[]>({
  key: "instanceIds",
  default: [],
  effects_UNSTABLE: [persistAtom],
});

// Array of registered services in an instance IDs
//    NOTE: Typically used together with `atomFamilies` to track IDs in the collection
export const instanceServiceIDs = atom<string[]>({
  key: "instanceProxyServiceIDs",
  default: [],
});

// Collection of instances registered
export const instances = atomFamily<Instance, number>({
  key: "instance",
  default: DefaultInstance,
  effects_UNSTABLE: [persistAtom],
});

// Collection of services registered
export const instanceService = atomFamily<InstanceProxyService, string>({
  key: "instanceProxyService",
  default: DefaultInstanceProxyService,
});

// A selector that fetch the proxy and services details of an instance
export const instanceProxySelector = selectorFamily({
  key: "getProxyAPI",
  get:
    (id: number) =>
    ({ get }) => {
      // Get the instance we are looking for
      const instance = get(instances(id));

      // Return the proxy included in the instance from the API
      return fetchProxy(instance.host);
    },
});

// Selector that returns the list of registered services
export const instanceProxyServiceSelector = selector({
  key: "getProxyServices",
  get: ({ get }) => {
    const ids = get(instanceServiceIDs);
    let services = [];
    for (let id of ids) {
      services.push(get(instanceService(id)));
    }

    return services;
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
