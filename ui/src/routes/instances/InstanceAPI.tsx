import { Instance, InstanceService } from "../../recoil/atoms/instances";
import { DefaultService } from "../../recoil/atoms/services";

export const fetchServices = async (instance: Instance) => {
  await new Promise((resolve) => setTimeout(resolve, 800));
  let services;

  let mockService = {
    ...DefaultService,
    running: false,
    proxy: 0,
    name: "CoAP",
  };

  services = [mockService, mockService, mockService];

  return services;
};

export const patchService = async (id: number, service: InstanceService) => {
  await new Promise((resolve) => setTimeout(resolve, 800));
  let servs;

  return service;
};
