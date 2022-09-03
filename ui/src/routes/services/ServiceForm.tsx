import { InteractionOptions, NetworkOptions } from "../../constants/globals";

export const ServiceFormFields = [
  {
    name: "name",
    help: "Name of the service",
    fieldattrs: {
      type: "text",
      required: true,
      placeholder: "Type a name for the service",
    },
  },
  {
    name: "host",
    help: "Host in where the service can be reached. E.g.: localhost, 0.0.0.0, etc.",
    fieldattrs: {
      type: "text",
      required: true,
      placeholder: "localhost",
    },
  },
  {
    name: "port",
    help: "Port address of the service. A number between 1 and 65535",
    fieldattrs: {
      type: "number",
      required: true,
      min: 1,
      max: 65535,
      placeholder: 1111,
    },
  },
  {
    name: "network",
    readOnly: true,
    type: "select",
    help: "Network layer protocol",
    fieldattrs: {
      options: NetworkOptions,
    },
  },
  {
    name: "interaction",
    readOnly: true,
    type: "select",
    help: "Service interaction level. An emulated service should be low; whereas a fully flagged service will be high.",
    fieldattrs: {
      options: InteractionOptions,
    },
  },
];
