import { InteractionOptions, NetworkOptions } from "../../constants/globals";
import { Instance, InstanceProxyService } from "../../recoil/atoms/instances";
import { Service } from "../../recoil/atoms/services";

export const fetchProxy = async (instance: Instance) => {
  await new Promise((resolve) => setTimeout(resolve, 800));
  let services;

  let mockProxy = {
    id: "",
    port: 2,
    status: "stopped",
    service: {
      id: "",
      name: "CoAP",
      network: NetworkOptions[0],
      interaction: InteractionOptions[0],
      host: "localhost",
      port: 2022,
    },
  };

  services = [mockProxy, mockProxy, mockProxy];

  return services;
};

export const patchService = async (id: number, service: Service) => {
  await new Promise((resolve) => setTimeout(resolve, 800));
  let servs;

  return service;
};

export const patchProxy = async (id: number, proxy: InstanceProxyService) => {
  await new Promise((resolve) => setTimeout(resolve, 800));
  let servs;

  return proxy;
};

export const changeProxyStatus = async (
  instanceId: number,
  proxyId: string,
  status: string
) => {
  await new Promise((resolve) => setTimeout(resolve, 800));
  return status;
};
