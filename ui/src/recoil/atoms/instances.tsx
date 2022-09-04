import { atom, atomFamily, selectorFamily } from "recoil";
import { Profile } from "./profiles";
import { recoilPersist } from "recoil-persist";
import { fetchServices } from "../../routes/instances/InstanceAPI";
import { DefaultService } from "./services";
import { InteractionOption, NetworkOption } from "../../constants/globals";

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

export type InstanceService = {
  name: string;
  interaction: InteractionOption;
  network: NetworkOption;
  host: string;
  port: Number;
  running: boolean;
  proxy: number | undefined;
};

const DefaultInstanceService = {
  ...DefaultService,
  running: false,
  proxy: undefined,
};

/**
 * Instance services.
 * This atom represents the services registered in an instance.
 * The atom will be populated through the API, comparing the services registered
 * in the local storage to the ones received from the API.
 */
export const instanceServices = atomFamily<InstanceService, number>({
  key: "instanceServices",
  default: DefaultInstanceService as InstanceService & { [key: string]: any },
});

export const instanceServicesIds = atom({
  key: "instanceServicesIds",
  default: [],
});

export const instanceServicesSelector = selectorFamily({
  key: "getServicesAPI",
  get:
    (id: number) =>
    ({ get }) => {
      // Get the instance we are looking for
      const instance = get(instances(id));

      // Return the services included in the instance from the API
      return fetchServices(instance);
    },
});
