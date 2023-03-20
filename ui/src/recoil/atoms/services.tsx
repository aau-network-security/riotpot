import { atom, selector, selectorFamily } from "recoil";
import {
  InteractionOption,
  InteractionOptions,
  NetworkOption,
  NetworkOptions,
} from "../../constants/globals";
import { recoilPersist } from "recoil-persist";
import { ValidatorProps } from "./common";

const { persistAtom } = recoilPersist();

export type Service = {
  id: string;
  name: string;
  network: NetworkOption;
  interaction: InteractionOption;
  host: string;
  port: number;
};

export const DefaultService = {
  id: "",
  name: "",
  network: NetworkOptions[0],
  interaction: InteractionOptions[0],
  host: "",
  port: 1,
} as Service;

export const services = atom({
  key: "servicesList",
  default: [] as Service[],
  effects_UNSTABLE: [persistAtom],
});

export const servicesFilter = selectorFamily({
  key: "service/default",
  get:
    (id: string) =>
    ({ get }) => {
      const servs = get(services);
      return servs.find((x: Service) => x.id === id);
    },
  set:
    (id: string) =>
    ({ get, set }, newValue) => {
      const servs: any = get(services);

      // Replaces a service
      const servInd = servs.findIndex((x: Service) => x.id === id);
      let cp = [...servs];
      cp[servInd] = { ...cp[servInd], ...newValue };

      return set(services, cp);
    },
});

export const servicesList = selector({
  key: "servicesListSelector",
  get: ({ get }) => get(services),
  set: ({ set }, newService) => set(services, newService),
});

export const servicesFormFields = atom({
  key: "serviceFormFields",
  default: DefaultService as { [key: string]: any },
});

export const serviceFormField = selectorFamily({
  key: "serviceFormField",
  get:
    (field?: string) =>
    ({ get }) =>
      !!field ? get(servicesFormFields)[field] : get(servicesFormFields),

  set:
    (field?: string) =>
    ({ set }, newValue) =>
      !!field
        ? // If a field is passed, set the value for the field
          set(servicesFormFields, (prevState: any) => ({
            ...prevState,
            [field]: newValue,
          }))
        : // Otherwise, update the item values
          set(servicesFormFields, (prev) => ({ ...prev, newValue })),
});

export const serviceFormErrors = atom({
  key: "serviceFormErrors",
  default: {} as { [key: string]: string },
});

export const serviceFormFieldErrors = selectorFamily({
  key: "serviceFormFieldErrors",
  get:
    (field: string) =>
    ({ get }) =>
      get(serviceFormErrors)[field],
});

// Function that validates that two services do not have the same name
export const ServiceNameValidator = (props: ValidatorProps) => {
  props.query.forEach((q) => {
    if (q.name === props.element.name) {
      props.errors((prev: any) => ({ ...prev, name: "Invalid value" }));
    }
  });
};
